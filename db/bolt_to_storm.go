package db

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"time"

	"github.com/Bit-Nation/panthalassa/crypto/aes"
	"github.com/Bit-Nation/panthalassa/keyManager"
	"github.com/Bit-Nation/panthalassa/queue"
	"github.com/Bit-Nation/x3dh"
	"github.com/asdine/storm"
	dr "github.com/tiabc/doubleratchet"
	bolt "go.etcd.io/bbolt"
)

type BoltToStormMigration struct {
	Km *keyManager.KeyManager
}

func (m *BoltToStormMigration) Migrate(db *storm.DB) error {
	
	// migrate queue jobs
	err := db.Bolt.Update(func(tx *bolt.Tx) error {
		queueStorage := queue.NewStorage(db.WithTransaction(tx))

		b := tx.Bucket([]byte("queue_storage"))
		if b == nil {
			return nil
		}

		type Job struct {
			ID   string                 `json:"id"`
			Type string                 `json:"type"`
			Data map[string]interface{} `json:"data"`
		}

		err := b.ForEach(func(oldJob, v []byte) error {

			// unmarshal old job
			var j Job
			if err := json.Unmarshal(v, &j); err != nil {
				return err
			}

			// persist job
			err := queueStorage.PersistJob(&queue.Job{
				Type: j.Type,
				Data: j.Data,
			})
			if err != nil {
				return err
			}

			// delete
			return b.Delete(oldJob)

		})
		if err != nil {
			return err
		}

		return tx.DeleteBucket([]byte("queue_storage"))
	})
	if err != nil {
		return err
	}

	// migrate double ratchet keys
	err = db.Bolt.Update(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte("double_ratchet_key_store"))
		if b == nil {
			return nil
		}

		drKeyStorage := NewBoltDRKeyStorage(db.WithTransaction(tx), m.Km)

		err := b.ForEach(func(plainDrKey, v []byte) error {

			if len(plainDrKey) != 32 {
				return errors.New("got dr key that isn't 32 bytes long")
			}

			drKey := dr.Key{}
			copy(drKey[:], plainDrKey)

			// key
			keys := b.Bucket(plainDrKey)
			if keys == nil {
				return nil
			}

			err := keys.ForEach(func(msgNumber, msgKey []byte) error {

				ct, err := aes.Unmarshal(msgKey)
				if err != nil {
					return err
				}

				plainSecret, err := m.Km.AESDecrypt(ct)
				if err != nil {
					return err
				}

				if len(plainSecret) != 32 {
					return errors.New("got invalid message key")
				}

				mk := dr.Key{}
				copy(mk[:], plainSecret)

				drKeyStorage.Put(
					drKey,
					uint(binary.BigEndian.Uint64(msgNumber)),
					mk,
				)

				return nil

			})
			if err != nil {
				return err
			}

			return b.Delete(plainDrKey)

		})
		if err != nil {
			return err
		}

		return tx.DeleteBucket([]byte("double_ratchet_key_store"))

	})
	if err != nil {
		return err
	}

	// migrate one time pre keys
	err = db.Bolt.Update(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte("one_time_pre_keys"))
		if b == nil {
			return nil
		}

		oneTimePreKeyStorage := NewBoltOneTimePreKeyStorage(db.WithTransaction(tx), m.Km)

		err := b.ForEach(func(pubKey, privKeyCT []byte) error {

			xPub := x3dh.PublicKey{}
			xPriv := x3dh.PrivateKey{}

			if len(pubKey) != 32 {
				return errors.New("got invalid public key during signed pre key migration")
			}

			ct, err := aes.Unmarshal(privKeyCT)
			if err != nil {
				return err
			}

			plainPrivKey, err := m.Km.AESDecrypt(ct)
			if err != nil {
				return err
			}

			if len(plainPrivKey) != 32 {
				return errors.New("got invalid private key during one time key migration")
			}

			copy(xPub[:], pubKey)
			copy(xPriv[:], plainPrivKey)

			return oneTimePreKeyStorage.Put([]x3dh.KeyPair{
				x3dh.KeyPair{
					PublicKey:  xPub,
					PrivateKey: xPriv,
				},
			})

			return b.Delete(pubKey)

		})
		if err != nil {
			return err
		}

		return tx.DeleteBucket([]byte("one_time_pre_keys"))
	})
	if err != nil {
		return err
	}

	// migrate signed pre keys
	err = db.Bolt.Update(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte("signed_pre_keys"))
		if b == nil {
			return nil
		}

		signedPreKeyStorage := NewBoltSignedPreKeyStorage(db.WithTransaction(tx), m.Km)

		type OldSignedPreKey struct {
			ValidTill  int64           `json:"valid_till"`
			PrivateKey x3dh.PrivateKey `json:"private_key"`
			PublicKey  x3dh.PublicKey  `json:"public_key"`
			Version    uint            `json:"version"`
		}

		err := b.ForEach(func(pubKey, privKeyCT []byte) error {
			if len(pubKey) != 32 {
				return errors.New("got invalid public key during signed pre key migration")
			}

			ct, err := aes.Unmarshal(privKeyCT)
			if err != nil {
				return err
			}

			rawSignedPrivKey, err := m.Km.AESDecrypt(ct)
			if err != nil {
				return err
			}

			signedPreKey := OldSignedPreKey{}
			if err := json.Unmarshal(rawSignedPrivKey, &signedPreKey); err != nil {
				return err
			}

			privKey := signedPreKey.PrivateKey

			if len(privKey) != 32 {
				return errors.New("got invalid private key during signed pre key migration")
			}

			xPub := x3dh.PublicKey{}
			copy(xPub[:], pubKey)

			return signedPreKeyStorage.Put(x3dh.KeyPair{
				PublicKey:  xPub,
				PrivateKey: privKey,
			})

			return b.Delete(pubKey)

		})
		if err != nil {
			return err
		}

		return tx.DeleteBucket([]byte("signed_pre_keys"))
	})
	if err != nil {
		return err
	}

	// migrate chats
	err = db.Bolt.Update(func(tx *bolt.Tx) error {

		chats := tx.Bucket([]byte("private_chat"))
		if chats == nil {
			return nil
		}

		chatDB := NewChatStorage(db.WithTransaction(tx), []func(event MessagePersistedEvent){}, m.Km)

		type OldDAppMessage struct {
			DAppPublicKey []byte                 `json:"dapp_public_key"`
			Type          string                 `json:"type"`
			Params        map[string]interface{} `json:"params"`
			ShouldSend    bool                   `json:"should_send"`
		}

		type OldMessage struct {
			ID         string       `json:"message_id"`
			Version    uint         `json:"version"`
			Status     Status       `json:"status"`
			Received   bool         `json:"received"`
			DApp       *DAppMessage `json:"dapp"`
			Message    []byte       `json:"message"`
			CreatedAt  int64        `json:"created_at"`
			Sender     []byte       `json:"sender"`
			DatabaseID int64        `json:"db_id"`
		}

		return chats.ForEach(func(partner, _ []byte) error {

			if err := chatDB.CreateChat(partner); err != nil {
				return err
			}

			chat, err := chatDB.GetChat(partner)
			if err != nil {
				return err
			}
			if chat == nil {
				return errors.New("got nil chat after creating it")
			}

			partnerChat := chats.Bucket(partner)
			if partnerChat == nil {
				return nil
			}

			return partnerChat.ForEach(func(_, rawEncryptedMessage []byte) error {

				ct, err := aes.Unmarshal(rawEncryptedMessage)
				if err != nil {
					return err
				}

				rawMessage, err := m.Km.AESDecrypt(ct)
				if err != nil {
					return err
				}

				m := OldMessage{}
				if err := json.Unmarshal(rawMessage, &m); err != nil {
					return err
				}

				newMessage := Message{
					ID:        m.ID,
					Version:   m.Version,
					Status:    Status(m.Status),
					Received:  m.Received,
					Message:   m.Message,
					CreatedAt: m.CreatedAt,
					Sender:    m.Sender,
				}

				if m.DApp != nil {
					newMessage.DApp = &DAppMessage{
						DAppPublicKey: m.DApp.DAppPublicKey,
						Type:          m.DApp.Type,
						Params:        m.DApp.Params,
						ShouldSend:    m.DApp.ShouldSend,
					}
				}

				return chat.PersistMessage(newMessage)

			})

		})

	})
	if err != nil {
		return err
	}

	// shared secret migration
	err = db.Bolt.Update(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte("shared_secrets"))
		if b == nil {
			return nil
		}

		sharedSecretStorage := NewBoltSharedSecretStorage(db.WithTransaction(tx), m.Km)

		type OldSharedSecret struct {
			X3dhSS                x3dh.SharedSecret `json:"-"`
			Accepted              bool              `json:"accepted"`
			CreatedAt             time.Time         `json:"created_at"`
			DestroyAt             *time.Time        `json:"destroy_at"`
			EphemeralKey          x3dh.PublicKey    `json:"ephemeral_key"`
			EphemeralKeySignature []byte            `json:"ephemeral_key_signature"`
			UsedSignedPreKey      x3dh.PublicKey    `json:"used_signed_pre_key"`
			UsedOneTimePreKey     *x3dh.PublicKey   `json:"used_one_time_pre_key"`
			// the base id chosen by the initiator of the chat
			BaseID []byte `json:"base_id"`
			// id used for indexing (calculated based on a few parameters)
			ID []byte `json:"id"`
			// the id based on the chat init params
			IDInitParams []byte `json:"id_init_params"`
		}

		type PersistedSharedSecret struct {
			OldSharedSecret
			X3dhSS aes.CipherText `json:"x3dh_shared_secret"`
		}

		err := b.ForEach(func(partner, _ []byte) error {

			if len(partner) != 32 {
				return errors.New("got invalid partner during shared secret migration")
			}

			partnerBucket := b.Bucket(partner)
			if partnerBucket == nil {
				return nil
			}

			return partnerBucket.ForEach(func(_, rawSharedSecret []byte) error {

				s := PersistedSharedSecret{}
				if err := json.Unmarshal(rawSharedSecret, &s); err != nil {
					return err
				}

				// decrypt shared secret
				plainSharedSecret, err := m.Km.AESDecrypt(s.X3dhSS)
				if err != nil {
					return err
				}

				// make sure shared secret is correct
				if len(plainSharedSecret) != 32 {
					return errors.New("got invalid shared secret")
				}

				ss := x3dh.SharedSecret{}
				copy(ss[:], plainSharedSecret)

				newSharedSecret := SharedSecret{
					x3dhSS:                ss,
					Accepted:              s.Accepted,
					CreatedAt:             s.CreatedAt,
					DestroyAt:             s.DestroyAt,
					Partner:               partner,
					ID:                    s.BaseID,
					UsedOneTimePreKey:     s.UsedOneTimePreKey,
					UsedSignedPreKey:      s.UsedSignedPreKey,
					EphemeralKey:          s.EphemeralKey,
					EphemeralKeySignature: s.EphemeralKeySignature,
				}

				return sharedSecretStorage.Put(newSharedSecret)

			})

		})
		if err != nil {
			return err
		}

		return tx.DeleteBucket([]byte("shared_secrets"))

	})
	if err != nil {
		return err
	}

	return nil

}

func (m *BoltToStormMigration) Version() uint32 {
	return 1
}

package db

import (
	"encoding/binary"
	"encoding/json"
	"errors"

	aes "github.com/Bit-Nation/panthalassa/crypto/aes"
	keyManager "github.com/Bit-Nation/panthalassa/keyManager"
	queue "github.com/Bit-Nation/panthalassa/queue"
	x3dh "github.com/Bit-Nation/x3dh"
	storm "github.com/asdine/storm"
	bolt "github.com/coreos/bbolt"
	dr "github.com/tiabc/doubleratchet"
)

type BoltToStormMigration struct {
	Km *keyManager.KeyManager
}

func (m *BoltToStormMigration) Migrate(db *storm.DB) error {

	// migrate queue jobs
	queueStorage := queue.NewStorage(db)
	err := db.Bolt.Update(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte("queue_storage"))
		if b == nil {
			return nil
		}

		type Job struct {
			ID   string                 `json:"id"`
			Type string                 `json:"type"`
			Data map[string]interface{} `json:"data"`
		}

		return b.ForEach(func(_, v []byte) error {

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
			return tx.DeleteBucket([]byte("queue_storage"))

		})

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

		drKeyStorage := NewBoltDRKeyStorage(db, m.Km)

		return b.ForEach(func(plainDrKey, v []byte) error {

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

			return tx.DeleteBucket(plainDrKey)

		})

	})
	if err != nil {
		return nil
	}

	// migrate one time pre keys
	err = db.Bolt.Update(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte("one_time_pre_keys"))
		if b == nil {
			return nil
		}

		signedPreKeyStorage := NewBoltSignedPreKeyStorage(db, m.Km)

		return b.ForEach(func(pubKey, privKeyCT []byte) error {

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
				return errors.New("got invalid private key during signed pre key migration")
			}

			copy(xPub[:], pubKey)
			copy(xPriv[:], plainPrivKey)

			return signedPreKeyStorage.Put(x3dh.KeyPair{
				PublicKey:  xPub,
				PrivateKey: xPriv,
			})

			return b.Delete(pubKey)

		})
	})
	if err != nil {
		return err
	}

	/**
	// shared secret migration
	err = m.db.Bolt.Update(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte("shared_secrets"))
		if b == nil {
			return nil
		}

		type SharedSecret struct {
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

		return b.ForEach(func(k, rawSharedSecret []byte) error {

		})

	})
	if err != nil {
		return err
	}
	*/

	// migrate chats
	err = db.Bolt.Update(func(tx *bolt.Tx) error {

		chats := tx.Bucket([]byte("private_chat"))
		if chats == nil {
			return nil
		}

		chatDB := NewChatStorage(db, []func(event MessagePersistedEvent){}, m.Km)

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

				m := Message{}
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

	return nil

}

func (m *BoltToStormMigration) Version() uint32 {
	return 1
}

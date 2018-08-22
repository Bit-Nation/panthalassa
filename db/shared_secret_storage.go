package db

import (
	"bytes"
	"encoding/json"
	"errors"
	"sort"
	"time"

	aes "github.com/Bit-Nation/panthalassa/crypto/aes"
	keyManager "github.com/Bit-Nation/panthalassa/keyManager"
	x3dh "github.com/Bit-Nation/x3dh"
	bolt "github.com/coreos/bbolt"
	ed25519 "golang.org/x/crypto/ed25519"
)

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

// the persistedSharedSecret is almost the same as SharedSecret except for
// that the X3dhSS value is an AES cipher text.
type persistedSharedSecret struct {
	SharedSecret
	X3dhSS aes.CipherText `json:"x3dh_shared_secret"`
}

var (
	sharedSecretBucketName = []byte("shared_secrets")
)

// use this to decrypt a persisted shared secret
var decryptPersistedSharedSecret = func(ss persistedSharedSecret, km *keyManager.KeyManager) (*SharedSecret, error) {

	// decrypt shared secret
	plainSharedSecret, err := km.AESDecrypt(ss.X3dhSS)
	if err != nil {
		return nil, err
	}

	// shared secret must have a length of 32 bytes
	if len(plainSharedSecret) != 32 {
		return nil, errors.New("shared secret must have a length of 32 bytes")
	}

	// copy shared secret over
	copy(ss.SharedSecret.X3dhSS[:], plainSharedSecret[:])

	return &ss.SharedSecret, nil

}

type SharedSecretStorage interface {
	HasAny(key ed25519.PublicKey) (bool, error)
	// must return an error if no shared secret found
	GetYoungest(key ed25519.PublicKey) (*SharedSecret, error)
	Put(chatPartner ed25519.PublicKey, ss SharedSecret) error
	// check if a secret for a chat initialization message exists
	SecretForChatInitMsg(partner ed25519.PublicKey, chatInitID []byte) (*SharedSecret, error)
	// accept will mark the given shared secret as accepted
	// and will set a destroy date for all other shared secrets
	Accept(partner ed25519.PublicKey, sharedSec *SharedSecret) error
	// get sender public key and shared secret id
	Get(key ed25519.PublicKey, sharedSecretID []byte) (*SharedSecret, error)
}

func NewBoltSharedSecretStorage(db *bolt.DB, km *keyManager.KeyManager) *BoltSharedSecretStorage {
	return &BoltSharedSecretStorage{
		db: db,
		km: km,
	}
}

type BoltSharedSecretStorage struct {
	db *bolt.DB
	km *keyManager.KeyManager
}

func (b *BoltSharedSecretStorage) HasAny(partner ed25519.PublicKey) (bool, error) {
	hasAny := false
	err := b.db.View(func(tx *bolt.Tx) error {

		// shared secrets bucket
		sharedSecretBucket := tx.Bucket(sharedSecretBucketName)
		if sharedSecretBucket == nil {
			return nil
		}

		// shared secrets with partner
		sharedSecretsPartner := sharedSecretBucket.Bucket(partner)
		if sharedSecretsPartner == nil {
			return nil
		}

		// in the case there are key values pairs,
		// we do have some shared secrets
		if sharedSecretsPartner.Stats().KeyN != 0 {
			hasAny = true
		}

		return nil

	})
	return hasAny, err
}

func (b *BoltSharedSecretStorage) GetYoungest(partner ed25519.PublicKey) (*SharedSecret, error) {
	shSec := new(SharedSecret)
	shSec = nil
	err := b.db.View(func(tx *bolt.Tx) error {

		// shared secrets bucket
		sharedSecretBucket := tx.Bucket(sharedSecretBucketName)
		if sharedSecretBucket == nil {
			return nil
		}

		// shared secrets with partner
		sharedSecretsPartner := sharedSecretBucket.Bucket(partner)
		if sharedSecretsPartner == nil {
			return nil
		}

		// fetch all shared secrets
		sharedSecrets := []persistedSharedSecret{}
		err := sharedSecretsPartner.ForEach(func(k, v []byte) error {
			tmpShSec := persistedSharedSecret{}
			if err := json.Unmarshal(v, &tmpShSec); err != nil {
				return err
			}
			sharedSecrets = append(sharedSecrets, tmpShSec)
			return nil
		})
		if err != nil {
			return err
		}
		logger.Debugf("fetched %d shared secrets from shared secret storage", len(sharedSecrets))

		// sort the shared secrets based on the createdAt
		sort.Slice(sharedSecrets, func(i, j int) bool {
			return sharedSecrets[i].CreatedAt.After(sharedSecrets[j].CreatedAt)
		})

		// if we couldn't find any shared secrets we can just return
		if len(sharedSecrets) <= 0 {
			return nil
		}

		// decrypt persisted shared secret
		shSec, err = decryptPersistedSharedSecret(sharedSecrets[0], b.km)
		return err
	})
	return shSec, err
}

func (b *BoltSharedSecretStorage) Put(chatPartner ed25519.PublicKey, ss SharedSecret) error {

	logger.Debugf("going to persist shared secret - for: %x, init params id: ", chatPartner, ss.IDInitParams)

	if len(ss.BaseID) != 32 {
		return errors.New("can't persisted shared secret with a base id len != 32")
	}

	return b.db.Update(func(tx *bolt.Tx) error {

		// shared secret bucket
		sharedSecretBucket, err := tx.CreateBucketIfNotExists(sharedSecretBucketName)
		if err != nil {
			return err
		}

		// shared secret partner bucket
		sharedSecPartnerBucket, err := sharedSecretBucket.CreateBucketIfNotExists(chatPartner)
		if err != nil {
			return err
		}

		// created persisted representation of shared secret
		persistedSharedSec := persistedSharedSecret{}
		persistedSharedSec.SharedSecret = ss
		persistedSharedSec.X3dhSS, err = b.km.AESEncrypt(ss.X3dhSS[:])

		// marshal persisted shared secret
		rawPersSharedSec, err := json.Marshal(persistedSharedSec)
		if err != nil {
			return err
		}

		return sharedSecPartnerBucket.Put(ss.BaseID, rawPersSharedSec)

	})
}

func (b *BoltSharedSecretStorage) SecretForChatInitMsg(chatPartner ed25519.PublicKey, chatInitID []byte) (*SharedSecret, error) {
	ss := new(SharedSecret)
	ss = nil

	err := b.db.View(func(tx *bolt.Tx) error {

		// shared secret bucket
		sharedSecretBucket := tx.Bucket(sharedSecretBucketName)
		if sharedSecretBucket == nil {
			return nil
		}

		// shared secret partner bucket
		sharedSecPartnerBucket := sharedSecretBucket.Bucket(chatPartner)
		if sharedSecPartnerBucket == nil {
			return nil
		}

		return sharedSecPartnerBucket.ForEach(func(k, rawPersSharedSecret []byte) error {

			// unmarshal raw persisted shared secret
			persSS := persistedSharedSecret{}
			if err := json.Unmarshal(rawPersSharedSecret, &persSS); err != nil {
				return err
			}

			// check if the secrets are equal
			if bytes.Equal(persSS.SharedSecret.IDInitParams, chatInitID) {
				fetchedSS, err := b.Get(chatPartner, persSS.BaseID)
				if err != nil {
					return err
				}
				ss = fetchedSS
			}

			return nil

		})

	})

	return ss, err

}

func (b *BoltSharedSecretStorage) Accept(partner ed25519.PublicKey, sharedSec *SharedSecret) error {
	sharedSec.Accepted = true
	return b.Put(partner, *sharedSec)
}

func (b *BoltSharedSecretStorage) Get(partner ed25519.PublicKey, sharedSecretID []byte) (*SharedSecret, error) {
	ss := new(SharedSecret)
	ss = nil

	err := b.db.View(func(tx *bolt.Tx) error {

		// shared secret bucket
		sharedSecretBucket := tx.Bucket(sharedSecretBucketName)
		if sharedSecretBucket == nil {
			return nil
		}

		// shared secret partner bucket
		sharedSecPartnerBucket := sharedSecretBucket.Bucket(partner)
		if sharedSecPartnerBucket == nil {
			return nil
		}

		// fetch raw persisted shared secret
		rawPersSharedSecret := sharedSecPartnerBucket.Get(sharedSecretID)
		if rawPersSharedSecret == nil {
			return nil
		}

		// unmarshal shared secret
		persSharedSecret := persistedSharedSecret{}
		if err := json.Unmarshal(rawPersSharedSecret, &persSharedSecret); err != nil {
			return nil
		}

		// decrypt shared secret
		shSec, err := decryptPersistedSharedSecret(persSharedSecret, b.km)
		if err != nil {
			return err
		}
		ss = shSec

		return nil
	})

	return ss, err
}

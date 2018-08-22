package db

import (
	"encoding/json"
	"fmt"
	"time"

	aes "github.com/Bit-Nation/panthalassa/crypto/aes"
	keyManager "github.com/Bit-Nation/panthalassa/keyManager"
	x3dh "github.com/Bit-Nation/x3dh"
	bolt "github.com/coreos/bbolt"
)

var (
	signedPreKeyBucketName = []byte("signed_pre_keys")
)

const (
	SignedPreKeyValidTimeFrame = time.Hour * 24 * 60
)

type SignedPreKeyStorage interface {
	// persist the signed pre key
	// @todo don't forget to give the option to register a listener in the put function
	// @todo publish signed pre key to backend
	Put(signedPreKey x3dh.KeyPair) error
	Get(publicKey x3dh.PublicKey) (*x3dh.PrivateKey, error)
	All() []*x3dh.KeyPair
}

type SignedPreKey struct {
	ValidTill  int64           `json:"valid_till"`
	PrivateKey x3dh.PrivateKey `json:"private_key"`
	PublicKey  x3dh.PublicKey  `json:"public_key"`
	Version    uint            `json:"version"`
}

type BoltSignedPreKeyStorage struct {
	db *bolt.DB
	km *keyManager.KeyManager
}

func NewBoltSignedPreKeyStorage(db *bolt.DB, km *keyManager.KeyManager) *BoltSignedPreKeyStorage {
	return &BoltSignedPreKeyStorage{
		db: db,
		km: km,
	}
}

func (s *BoltSignedPreKeyStorage) Put(signedPreKey x3dh.KeyPair) error {
	return s.db.Update(func(tx *bolt.Tx) error {

		// signed pre key bucket
		signedPreKeyBucket, err := tx.CreateBucketIfNotExists(signedPreKeyBucketName)
		if err != nil {
			return err
		}

		spk := SignedPreKey{
			ValidTill:  time.Now().Add(SignedPreKeyValidTimeFrame).Unix(),
			PrivateKey: signedPreKey.PrivateKey,
			PublicKey:  signedPreKey.PublicKey,
			Version:    1,
		}
		// raw signed pre key
		rawSignedPreKey, err := json.Marshal(spk)
		if err != nil {
			return err
		}

		// encrypt signed pre key
		ct, err := s.km.AESEncrypt(rawSignedPreKey)
		if err != nil {
			return err
		}

		// raw encrypted signed pre  key
		rawEncSignedPreKey, err := ct.Marshal()
		if err != nil {
			return err
		}

		// persist signed pre key
		return signedPreKeyBucket.Put(signedPreKey.PublicKey[:], rawEncSignedPreKey)

	})
}

func (s *BoltSignedPreKeyStorage) getSignedPreKey(publicKey x3dh.PublicKey) (*SignedPreKey, error) {
	signedPreKey := new(SignedPreKey)
	signedPreKey = nil

	err := s.db.View(func(tx *bolt.Tx) error {

		// signed pre key bucket
		signedPreKeyBucket := tx.Bucket(signedPreKeyBucketName)
		if signedPreKeyBucket == nil {
			return nil
		}

		// raw encrypted signed pre key
		rawEncryptedSignedPreKey := signedPreKeyBucket.Get(publicKey[:])
		if rawEncryptedSignedPreKey == nil {
			return nil
		}

		// unmarshal aes cipher text
		ct := aes.CipherText{}
		if err := json.Unmarshal(rawEncryptedSignedPreKey, &ct); err != nil {
			return err
		}

		// decrypt signed pre key
		rawSignedPreKey, err := s.km.AESDecrypt(ct)
		if err != nil {
			return err
		}

		skp := SignedPreKey{}
		if err := json.Unmarshal(rawSignedPreKey, &skp); err != nil {
			return err
		}
		signedPreKey = &skp

		return nil
	})

	return signedPreKey, err

}

func (s *BoltSignedPreKeyStorage) Get(publicKey x3dh.PublicKey) (*x3dh.PrivateKey, error) {
	privKey := new(x3dh.PrivateKey)
	privKey = nil

	err := s.db.View(func(tx *bolt.Tx) error {

		// fetch signed pre key
		signedPreKey, err := s.getSignedPreKey(publicKey)
		if err != nil {
			return err
		}

		privKey = &signedPreKey.PrivateKey

		return nil
	})

	return privKey, err
}

func (s *BoltSignedPreKeyStorage) All() []*x3dh.KeyPair {

	signedPreKeys := []*x3dh.KeyPair{}

	err := s.db.View(func(tx *bolt.Tx) error {

		// signed pre key bucket
		signedPreKeyBucket := tx.Bucket(signedPreKeyBucketName)
		if signedPreKeyBucket == nil {
			return nil
		}

		return signedPreKeyBucket.ForEach(func(pubKey, _ []byte) error {

			// key must be a x3dh public key
			if len(pubKey) != 32 {
				return fmt.Errorf("got invalid public key: %x - must have length of 32 bytes", pubKey)
			}

			pub := x3dh.PublicKey{}
			copy(pub[:], pubKey)

			signedPreKey, err := s.getSignedPreKey(pub)
			if err != nil {
				return err
			}

			if signedPreKey == nil {
				return fmt.Errorf("tried to fetch signed pre key of public key: %x", pubKey)
			}

			if signedPreKey.PrivateKey == [32]byte{} {
				return fmt.Errorf("got invalid private key (32x0) for public key: %x", pubKey)
			}

			// append signed pre key
			signedPreKeys = append(signedPreKeys, &x3dh.KeyPair{
				PrivateKey: signedPreKey.PrivateKey,
				PublicKey:  signedPreKey.PublicKey,
			})

			return nil

		})

	})

	// @todo we should return the error instead of just logging it
	if err != nil {
		logger.Error(err)
	}

	return signedPreKeys

}

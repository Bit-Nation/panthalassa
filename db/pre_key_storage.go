package db

import (
	"errors"

	aes "github.com/Bit-Nation/panthalassa/crypto/aes"
	km "github.com/Bit-Nation/panthalassa/keyManager"
	x3dh "github.com/Bit-Nation/x3dh"
	bolt "github.com/coreos/bbolt"
)

var preKeyStoreBucket = []byte("pre_key_store_bucket")

type PreKeyStorage struct {
	km *km.KeyManager
	db *bolt.DB
}

func (s *PreKeyStorage) Put(keyPair x3dh.KeyPair) error {

	return s.db.Update(func(tx *bolt.Tx) error {

		storage, err := tx.CreateBucketIfNotExists(preKeyStoreBucket)
		if err != nil {
			return err
		}

		encryptedPrivateKey, err := s.km.AESEncrypt(keyPair.PrivateKey[:])
		if err != nil {
			return err
		}

		rawEncryptedPrivateKey, err := encryptedPrivateKey.Marshal()
		if err != nil {
			return err
		}

		return storage.Put(keyPair.PublicKey[:], rawEncryptedPrivateKey)

	})

}

func (s *PreKeyStorage) Has(pubKey x3dh.PublicKey) (bool, error) {

	exist := false

	err := s.db.View(func(tx *bolt.Tx) error {

		// fetch pre key store bucket
		bucket := tx.Bucket(preKeyStoreBucket)
		if bucket == nil {
			return nil
		}

		// in the case the key does not exist nil will be returned
		if bucket.Get(pubKey[:]) != nil {
			exist = true
		}

		return nil
	})

	return exist, err

}

// get will return an error if the private key does not exist
// make sure to check first if the key exist with "has"
func (s *PreKeyStorage) Get(pubKey x3dh.PublicKey) (x3dh.PrivateKey, error) {

	var privKey x3dh.PrivateKey

	err := s.db.View(func(tx *bolt.Tx) error {

		// fetch pre key store bucket
		bucket := tx.Bucket(preKeyStoreBucket)
		if bucket == nil {
			return nil
		}

		encryptedPrivateKey := bucket.Get(pubKey[:])
		if encryptedPrivateKey == nil {
			return errors.New("couldn't find private key for given public key")
		}

		ct, err := aes.Unmarshal(encryptedPrivateKey)
		if err != nil {
			return err
		}

		rawPrivateKey, err := s.km.AESDecrypt(ct)
		if err != nil {
			return err
		}

		if len(rawPrivateKey) != 32 {
			return errors.New("private key must have a length of 32 bytes")
		}

		copy(privKey[:], rawPrivateKey)

		return nil
	})

	return privKey, err

}

func (s *PreKeyStorage) Delete(pubKey x3dh.PublicKey) error {

	return s.db.Update(func(tx *bolt.Tx) error {

		// fetch pre key store bucket
		bucket := tx.Bucket(preKeyStoreBucket)
		if bucket == nil {
			return nil
		}

		return bucket.Delete(pubKey[:])

	})

}

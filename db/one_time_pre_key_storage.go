package db

import (
	"errors"

	aes "github.com/Bit-Nation/panthalassa/crypto/aes"
	km "github.com/Bit-Nation/panthalassa/keyManager"
	x3dh "github.com/Bit-Nation/x3dh"
	bolt "github.com/coreos/bbolt"
)

var (
	oneTimePreKeyStorageBucketName = []byte("one_time_pre_keys")
)

type OneTimePreKeyStorage interface {
	Cut(pubKey x3dh.PublicKey) (*x3dh.PrivateKey, error)
	Count() (uint32, error)
	Put(keyPairs []x3dh.KeyPair) error
}

type BoltOneTimePreKeyStorage struct {
	db *bolt.DB
	km *km.KeyManager
}

func NewBoltOneTimePreKeyStorage(db *bolt.DB, km *km.KeyManager) *BoltChatMessageStorage {
	return &BoltChatMessageStorage{
		db: db,
		km: km,
	}
}

func (s *BoltChatMessageStorage) Cut(pubKey x3dh.PublicKey) (*x3dh.PrivateKey, error) {

	var privKey x3dh.PrivateKey

	err := s.db.Update(func(tx *bolt.Tx) error {

		// one time pre keys
		oneTimePreKeys, err := tx.CreateBucketIfNotExists(oneTimePreKeyStorageBucketName)
		if err != nil {
			return err
		}

		// fetch encrypted one time pre key of public key
		encryptedRawPriv := oneTimePreKeys.Get(pubKey[:])
		if encryptedRawPriv == nil {
			return errors.New("failed to fetch one time pre key private key for given public key")
		}

		// turn from byte slice to cipher text
		encryptedPriv, err := aes.Unmarshal(encryptedRawPriv)
		if err != nil {
			if err := oneTimePreKeys.Delete(pubKey[:]); err != nil {
				logger.Error(err)
			}
			return err
		}

		// decrypt the cipher text
		plainPriv, err := s.km.AESDecrypt(encryptedPriv)
		if err != nil {
			if err := oneTimePreKeys.Delete(pubKey[:]); err != nil {
				logger.Error(err)
			}
			return err
		}

		// make sure that we have a valid private key
		if len(plainPriv) != 32 {
			if err := oneTimePreKeys.Delete(pubKey[:]); err != nil {
				logger.Error(err)
			}
			return errors.New("got invalid x3dh private key for public key")
		}

		// copy over
		copy(privKey[:], plainPriv)

		// since this is the cut method we delete the private key once we fetched it
		if err := oneTimePreKeys.Delete(pubKey[:]); err != nil {
			logger.Error(err)
		}

		return nil

	})

	return &privKey, err

}

func (s *BoltChatMessageStorage) Count() (uint32, error) {

	var amount uint32

	err := s.db.View(func(tx *bolt.Tx) error {

		// one time pre key bucket
		oneTimePreKeys := tx.Bucket(oneTimePreKeyStorageBucketName)
		if oneTimePreKeys == nil {
			return nil
		}

		// fetch the amount of one time pre keys
		return oneTimePreKeys.ForEach(func(k, v []byte) error {
			amount++
			return nil
		})

	})

	return amount, err
}

func (s *BoltChatMessageStorage) Put(keyPairs []x3dh.KeyPair) error {

	return s.db.Update(func(tx *bolt.Tx) error {

		oneTimePreKeyBucket, err := tx.CreateBucketIfNotExists(oneTimePreKeyStorageBucketName)
		if err != nil {
			return err
		}

		// persist keys
		for _, keyPair := range keyPairs {

			// encrypt private key
			encryptedPriv, err := s.km.AESEncrypt(keyPair.PrivateKey[:])
			if err != nil {
				return err
			}

			// marshal encrypted private key
			rawEncryptedPriv, err := encryptedPriv.Marshal()
			if err != nil {
				return err
			}

			// persist private key
			err = oneTimePreKeyBucket.Put(keyPair.PublicKey[:], rawEncryptedPriv)
			if err != nil {
				return err
			}

		}

		return nil

	})

}

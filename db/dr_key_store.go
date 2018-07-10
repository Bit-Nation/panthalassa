package db

import (
	"encoding/binary"
	"errors"
	"fmt"

	aes "github.com/Bit-Nation/panthalassa/crypto/aes"
	km "github.com/Bit-Nation/panthalassa/keyManager"
	bolt "github.com/coreos/bbolt"
	dr "github.com/tiabc/doubleratchet"
)

type Store struct {
	db *bolt.DB
	km *km.KeyManager
}

var (
	doubleRatchetKeyStoreBucket = []byte("double_ratchet_key_store")
)

func uintToBytes(uint uint) []byte {
	num := make([]byte, 8)
	binary.LittleEndian.PutUint64(num, uint64(uint))
	return num
}

func bytesToUint(uint []byte) uint64 {
	return binary.LittleEndian.Uint64(uint)
}

func (s *Store) Get(k dr.Key, msgNum uint) (mk dr.Key, ok bool) {

	exist := false
	key := dr.Key{}

	err := s.db.View(func(tx *bolt.Tx) error {
		// create Fetch pre key store
		preKeyStore := tx.Bucket(doubleRatchetKeyStoreBucket)
		if preKeyStore == nil {
			return nil
		}
		// message store
		messageStore := preKeyStore.Bucket(k[:])
		if messageStore == nil {
			return nil
		}
		// encrypted pre key
		encryptedRawDRKey := messageStore.Get(uintToBytes(msgNum))
		if encryptedRawDRKey == nil {
			return nil
		}
		encryptedPreKey, err := aes.Unmarshal(encryptedRawDRKey)
		if err != nil {
			return err
		}
		// decrypted double
		plainText, err := s.km.AESDecrypt(encryptedPreKey)
		if err != nil {
			return err
		}
		if len(plainText) != 32 {
			return errors.New(fmt.Sprintf("double ratchet key is invalid (length: %d)", len(plainText)))
		}
		copy(key[:], plainText)
		exist = true
		return nil
	})

	if err != nil {
		//@todo not yet sure what todo
	}

	return key, exist

}

func (s *Store) Put(k dr.Key, msgNum uint, mk dr.Key) {

	err := s.db.Update(func(tx *bolt.Tx) error {
		preKeyStore, err := tx.CreateBucketIfNotExists(doubleRatchetKeyStoreBucket)
		if err != nil {
			return err
		}
		messageKeys, err := preKeyStore.CreateBucketIfNotExists(k[:])
		if err != nil {
			return err
		}
		ctStruct, err := s.km.AESEncrypt(mk[:])
		if err != nil {
			return err
		}
		ct, err := ctStruct.Marshal()
		if err != nil {
			return err
		}
		return messageKeys.Put(uintToBytes(msgNum), ct)
	})

	if err != nil {
		//@todo not sure what todo with it
	}

}

func (s *Store) DeleteMk(k dr.Key, msgNum uint) {

	err := s.db.Update(func(tx *bolt.Tx) error {
		drKeyStore := tx.Bucket(doubleRatchetKeyStoreBucket)
		if drKeyStore == nil {
			return nil
		}
		messageStore := drKeyStore.Bucket(k[:])
		if messageStore == nil {
			return nil
		}
		return messageStore.Delete(uintToBytes(msgNum))
	})

	if err != nil {
		//@todo not sure what todo with it
	}

}

func (s *Store) DeletePk(k dr.Key) {

	err := s.db.Update(func(tx *bolt.Tx) error {
		drKeyStore := tx.Bucket(doubleRatchetKeyStoreBucket)
		if drKeyStore == nil {
			return nil
		}
		return drKeyStore.DeleteBucket(k[:])
	})

	if err != nil {
		//@todo not sure what todo with it
	}

}

func (s *Store) Count(k dr.Key) uint {

	count := 0

	err := s.db.Update(func(tx *bolt.Tx) error {
		drKeyStore := tx.Bucket(doubleRatchetKeyStoreBucket)
		if drKeyStore == nil {
			return nil
		}
		messageStore := drKeyStore.Bucket(k[:])
		if messageStore == nil {
			return nil
		}
		return messageStore.ForEach(func(k, v []byte) error {
			count++
			return nil
		})
	})

	if err != nil {
		//@todo not sure what todo with it
	}

	return uint(count)

}

func (s *Store) All() map[dr.Key]map[uint]dr.Key {

	keys := map[dr.Key]map[uint]dr.Key{}

	err := s.db.Update(func(tx *bolt.Tx) error {
		// get double ratchet key store
		drKeyStore := tx.Bucket(doubleRatchetKeyStoreBucket)
		if drKeyStore == nil {
			return nil
		}
		return drKeyStore.ForEach(func(k, v []byte) error {
			// get message key store
			messageKeyStore := tx.Bucket(k)
			if messageKeyStore == nil {
				return errors.New("we MUST get a bucket here since all keys should have a bucket which is a mapping between a uint and a key")
			}
			if len(k) != 32 {
				return errors.New("invalid key size. must be exactly 32 bytes long")
			}
			var drKey dr.Key
			copy(drKey[:], k)
			messageKeys := map[uint]dr.Key{}
			err := messageKeyStore.ForEach(func(k, v []byte) error {
				key, exist := s.Get(drKey, uint(bytesToUint(k)))
				if !exist {
					return nil
				}
				messageKeys[uint(bytesToUint(k))] = key
				return nil
			})
			if err != nil {
				return err
			}
			keys[drKey] = messageKeys
			return nil
		})
	})

	if err != nil {
		//@todo not sure what todo with it
	}

	return keys

}

package db

import (
	"encoding/binary"
	"errors"
	"fmt"

	aes "github.com/Bit-Nation/panthalassa/crypto/aes"
	km "github.com/Bit-Nation/panthalassa/keyManager"
	bolt "github.com/coreos/bbolt"
	log "github.com/ipfs/go-log"
	dr "github.com/tiabc/doubleratchet"
)

type DRKeyStorage interface {
	Get(k dr.Key, msgNum uint) (mk dr.Key, ok bool)
	Put(k dr.Key, msgNum uint, mk dr.Key)
	DeleteMk(k dr.Key, msgNum uint)
	DeletePk(k dr.Key)
	Count(k dr.Key) uint
	All() map[dr.Key]map[uint]dr.Key
}

type BoltDRKeyStorage struct {
	db *bolt.DB
	km *km.KeyManager
}

func NewBoltDRKeyStorage(db *bolt.DB, km *km.KeyManager) *BoltDRKeyStorage {
	return &BoltDRKeyStorage{
		db: db,
		km: km,
	}
}

var logger = log.Logger("database")

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

func (s *BoltDRKeyStorage) Get(k dr.Key, msgNum uint) (mk dr.Key, ok bool) {

	exist := false
	key := dr.Key{}

	err := s.db.View(func(tx *bolt.Tx) error {
		// double ratchet key store
		preKeyStore := tx.Bucket(doubleRatchetKeyStoreBucket)
		if preKeyStore == nil {
			return nil
		}
		// message key store
		messageKeyStore := preKeyStore.Bucket(k[:])
		if messageKeyStore == nil {
			return nil
		}
		// encrypted message key
		encryptedRawDRKey := messageKeyStore.Get(uintToBytes(msgNum))
		if encryptedRawDRKey == nil {
			return nil
		}
		encryptedMessageKey, err := aes.Unmarshal(encryptedRawDRKey)
		if err != nil {
			return err
		}
		// decrypted message key
		plainText, err := s.km.AESDecrypt(encryptedMessageKey)
		if err != nil {
			return err
		}
		if len(plainText) != 32 {
			return errors.New(fmt.Sprintf("message key is invalid (length: %d)", len(plainText)))
		}
		copy(key[:], plainText)
		exist = true
		return nil
	})

	// @todo we need to change the dr package to use a better interface
	if err != nil {
		logger.Error(err)
	}

	return key, exist

}

func (s *BoltDRKeyStorage) Put(k dr.Key, msgNum uint, mk dr.Key) {

	err := s.db.Update(func(tx *bolt.Tx) error {
		// get double ratchet key store
		drKeyStore, err := tx.CreateBucketIfNotExists(doubleRatchetKeyStoreBucket)
		if err != nil {
			return err
		}
		// message keys
		messageKeys, err := drKeyStore.CreateBucketIfNotExists(k[:])
		if err != nil {
			return err
		}
		// encrypt the message key
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

	// @todo we need to change the dr package to use a better interface
	if err != nil {
		logger.Error(err)
	}

}

func (s *BoltDRKeyStorage) DeleteMk(k dr.Key, msgNum uint) {

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

	// @todo we need to change the dr package to use a better interface
	if err != nil {
		logger.Error(err)
	}

}

func (s *BoltDRKeyStorage) DeletePk(k dr.Key) {

	err := s.db.Update(func(tx *bolt.Tx) error {
		drKeyStore := tx.Bucket(doubleRatchetKeyStoreBucket)
		if drKeyStore == nil {
			return nil
		}
		return drKeyStore.DeleteBucket(k[:])
	})

	// @todo we need to change the dr package to use a better interface
	if err != nil {
		logger.Error(err)
	}

}

func (s *BoltDRKeyStorage) Count(k dr.Key) uint {

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
		logger.Error(err)
	}

	return uint(count)

}

func (s *BoltDRKeyStorage) All() map[dr.Key]map[uint]dr.Key {

	keys := map[dr.Key]map[uint]dr.Key{}

	err := s.db.View(func(tx *bolt.Tx) error {
		// get double ratchet key store
		drKeyStore := tx.Bucket(doubleRatchetKeyStoreBucket)
		if drKeyStore == nil {
			return nil
		}
		return drKeyStore.ForEach(func(k, v []byte) error {
			// get message key store
			messageKeyStore := drKeyStore.Bucket(k)
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
		logger.Error(err)
	}

	return keys

}

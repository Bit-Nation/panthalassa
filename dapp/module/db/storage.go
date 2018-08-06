package db

import (
	"errors"

	bolt "github.com/coreos/bbolt"
	ed25519 "golang.org/x/crypto/ed25519"
)

var dAppDBBucketName = []byte("dapp_db_bucket")

type Storage interface {
	Put(key, value []byte) error
	Get(key []byte) ([]byte, error)
	Has(key []byte) (bool, error)
	Delete(key []byte) error
}

type BoltStorage struct {
	db             *bolt.DB
	dAppSigningKey ed25519.PublicKey
}

func NewBoltStorage(db *bolt.DB, signingKey ed25519.PublicKey) (*BoltStorage, error) {
	return &BoltStorage{
		db:             db,
		dAppSigningKey: signingKey,
	}, nil
}

func fetchDAppDB(tx *bolt.Tx, dAppSigningKey ed25519.PublicKey) (*bolt.Bucket, error) {

	// make sure dAppSigningKey
	if len(dAppSigningKey) != 32 {
		return nil, errors.New("dApp signing key is too short")
	}

	// dApp bucket
	dAppDB, err := tx.CreateBucketIfNotExists(dAppDBBucketName)
	if err != nil {
		return nil, err
	}

	return dAppDB.CreateBucketIfNotExists(dAppSigningKey)

}

func fetchDAppDBView(tx *bolt.Tx, dAppSigningKey ed25519.PublicKey) (*bolt.Bucket, error) {

	// make sure dAppSigningKey
	if len(dAppSigningKey) != 32 {
		return nil, errors.New("dApp signing key is too short")
	}

	// dApp bucket
	dAppDB := tx.Bucket(dAppDBBucketName)
	if dAppDB == nil {
		return nil, nil
	}

	return dAppDB.Bucket(dAppSigningKey), nil

}

func (s *BoltStorage) Put(key, value []byte) error {
	return s.db.Update(func(tx *bolt.Tx) error {

		// dApp bucket
		dAppDB, err := fetchDAppDB(tx, s.dAppSigningKey)
		if err != nil {
			return err
		}

		return dAppDB.Put(key, value)

	})
}

func (s *BoltStorage) Get(key []byte) ([]byte, error) {
	var value []byte
	err := s.db.View(func(tx *bolt.Tx) error {

		// dApp bucket
		dAppDB, err := fetchDAppDBView(tx, s.dAppSigningKey)
		if err != nil {
			return err
		}
		if dAppDB == nil {
			return nil
		}

		value = dAppDB.Get(key)
		return nil

	})
	return value, err
}

func (s *BoltStorage) Has(key []byte) (bool, error) {
	exist := false
	err := s.db.View(func(tx *bolt.Tx) error {

		// dApp bucket
		dAppDB, err := fetchDAppDBView(tx, s.dAppSigningKey)
		if err != nil {
			return err
		}
		if dAppDB == nil {
			return nil
		}

		if dAppDB.Get(key) != nil {
			exist = true
		}
		return nil

	})
	return exist, err
}

func (s *BoltStorage) Delete(key []byte) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		// dApp bucket
		dAppDB, err := fetchDAppDB(tx, s.dAppSigningKey)
		if err != nil {
			return err
		}
		return dAppDB.Delete(key)
	})
}

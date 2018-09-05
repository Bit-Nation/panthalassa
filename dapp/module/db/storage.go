package db

import (
	storm "github.com/asdine/storm"
	sq "github.com/asdine/storm/q"
	ed25519 "golang.org/x/crypto/ed25519"
)

type Record struct {
	ID             int               `storm:"id,increment"`
	DAppSigningKey ed25519.PublicKey `storm:"index"`
	Key            []byte            `storm:"index"`
	Value          []byte
}

type Storage interface {
	Put(key, value []byte) error
	Get(key []byte) ([]byte, error)
	Has(key []byte) (bool, error)
	Delete(key []byte) error
}

type BoltStorage struct {
	db             *storm.DB
	dAppSigningKey ed25519.PublicKey
}

func NewBoltStorage(db *storm.DB, signingKey ed25519.PublicKey) (*BoltStorage, error) {
	return &BoltStorage{
		db:             db,
		dAppSigningKey: signingKey,
	}, nil
}

func (s *BoltStorage) Put(key, value []byte) error {
	return s.db.Save(&Record{
		DAppSigningKey: s.dAppSigningKey,
		Key:            key,
		Value:          value,
	})
}

func (s *BoltStorage) Get(key []byte) ([]byte, error) {
	q := s.db.Select(sq.And(
		sq.Eq("Key", key),
		sq.Eq("DAppSigningKey", s.dAppSigningKey),
	))

	var r Record
	return r.Value, q.First(&r)
}

func (s *BoltStorage) Has(key []byte) (bool, error) {

	q := s.db.Select(sq.And(
		sq.Eq("Key", key),
		sq.Eq("DAppSigningKey", s.dAppSigningKey),
	))

	amount, err := q.Count(&Record{})
	if err != nil {
		return false, err
	}

	if amount > 0 {
		return true, nil
	}

	return false, nil

}

func (s *BoltStorage) Delete(key []byte) error {
	q := s.db.Select(sq.And(
		sq.Eq("Key", key),
		sq.Eq("DAppSigningKey", s.dAppSigningKey),
	))
	return q.Delete(&Record{})
}

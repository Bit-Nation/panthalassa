package db

import (
	preKey "github.com/Bit-Nation/panthalassa/chat/prekey"
	storm "github.com/asdine/storm"
	sq "github.com/asdine/storm/q"
	ed25519 "golang.org/x/crypto/ed25519"
)

// user storage store meta data about users
type UserStorage interface {
	// don't forget to verify the signature when implementing this
	GetSignedPreKey(idKey ed25519.PublicKey) (*preKey.PreKey, error)
	PutSignedPreKey(idKey ed25519.PublicKey, key preKey.PreKey) error
}

type User struct {
	IdKey        ed25519.PublicKey `storm:"index,id,unique"`
	SignedPreKey preKey.PreKey
	Version      uint
}

type BoltUserStorage struct {
	db *storm.DB
}

func NewBoltUserStorage(db *storm.DB) *BoltUserStorage {
	return &BoltUserStorage{db: db}
}

func (s *BoltUserStorage) GetSignedPreKey(idKey ed25519.PublicKey) (*preKey.PreKey, error) {

	q := s.db.Select(sq.Eq("IdKey", idKey))
	amount, err := q.Count(&User{})
	if err != nil {
		return nil, err
	}
	if amount == 0 {
		return nil, nil
	}

	var u User
	if err := q.First(&u); err != nil {
		return nil, err
	}

	return &u.SignedPreKey, nil

}

func (s *BoltUserStorage) PutSignedPreKey(idKey ed25519.PublicKey, key preKey.PreKey) error {

	return s.db.Save(&User{
		IdKey:        idKey,
		SignedPreKey: key,
		Version:      1,
	})

}

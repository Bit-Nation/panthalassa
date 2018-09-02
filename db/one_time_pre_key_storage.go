package db

import (
	"errors"

	aes "github.com/Bit-Nation/panthalassa/crypto/aes"
	km "github.com/Bit-Nation/panthalassa/keyManager"
	x3dh "github.com/Bit-Nation/x3dh"
	storm "github.com/asdine/storm"
	sq "github.com/asdine/storm/q"
)

type oneTimePreKey struct {
	ID      int            `storm:"id,increment"`
	PubKey  x3dh.PublicKey `storm:"unique,index"`
	PrivKey aes.CipherText
}

type OneTimePreKeyStorage interface {
	Cut(pubKey []byte) (*x3dh.PrivateKey, error)
	Count() (uint32, error)
	Put(keyPairs []x3dh.KeyPair) error
}

type BoltOneTimePreKeyStorage struct {
	db storm.Node
	km *km.KeyManager
}

func NewBoltOneTimePreKeyStorage(db storm.Node, km *km.KeyManager) *BoltOneTimePreKeyStorage {
	return &BoltOneTimePreKeyStorage{
		db: db,
		km: km,
	}
}

func (s *BoltOneTimePreKeyStorage) Cut(pubKey []byte) (*x3dh.PrivateKey, error) {

	if len(pubKey) != 32 {
		return nil, errors.New("got invalid public key with length != 32")
	}
	x3dhPub := x3dh.PublicKey{}
	copy(x3dhPub[:], pubKey)

	privKey := new(x3dh.PrivateKey)

	// check if a record exist
	q := s.db.Select(sq.Eq("PubKey", x3dhPub))
	amount, err := q.Count(&oneTimePreKey{})
	if err != nil {
		return nil, err
	}
	if amount <= 0 {
		return nil, nil
	}

	// find one time pre key
	var oneTimePreKey oneTimePreKey
	if err := q.First(&oneTimePreKey); err != nil {
		return nil, err
	}

	// decrypt private key and copy over
	plainPrivKey, err := s.km.AESDecrypt(oneTimePreKey.PrivKey)
	if err != nil {
		return nil, err
	}
	if len(plainPrivKey) != 32 {
		return nil, errors.New("got invalid private key with len != 32")
	}
	copy(privKey[:], plainPrivKey)

	return privKey, s.db.DeleteStruct(&oneTimePreKey)

}

func (s *BoltOneTimePreKeyStorage) Count() (uint32, error) {
	amount, err := s.db.Count(&oneTimePreKey{})
	return uint32(amount), err
}

func (s *BoltOneTimePreKeyStorage) Put(keyPairs []x3dh.KeyPair) error {

	for _, keyPair := range keyPairs {

		// encrypt the private key
		ct, err := s.km.AESEncrypt(keyPair.PrivateKey[:])
		if err != nil {
			return err
		}

		// persist one time pre key
		err = s.db.Save(&oneTimePreKey{
			PubKey:  keyPair.PublicKey,
			PrivKey: ct,
		})
		if err != nil {
			return err
		}

	}

	return nil

}

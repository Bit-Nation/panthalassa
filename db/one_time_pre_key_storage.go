package db

import (
	"errors"

	aes "github.com/Bit-Nation/panthalassa/crypto/aes"
	km "github.com/Bit-Nation/panthalassa/keyManager"
	x3dh "github.com/Bit-Nation/x3dh"
	storm "github.com/asdine/storm"
)

type oneTimePreKey struct {
	PubKey  x3dh.PublicKey `storm:"unique,index"`
	PrivKey aes.CipherText
}

type OneTimePreKeyStorage interface {
	Cut(pubKey []byte) (*x3dh.PrivateKey, error)
	Count() (uint32, error)
	Put(keyPairs []x3dh.KeyPair) error
}

type BoltOneTimePreKeyStorage struct {
	db *storm.DB
	km *km.KeyManager
}

func NewBoltOneTimePreKeyStorage(db *storm.DB, km *km.KeyManager) *BoltOneTimePreKeyStorage {
	return &BoltOneTimePreKeyStorage{
		db: db,
		km: km,
	}
}

func (s *BoltOneTimePreKeyStorage) Cut(pubKey []byte) (*x3dh.PrivateKey, error) {

	privKey := new(x3dh.PrivateKey)

	// find one time pre key
	otpk := new(oneTimePreKey)
	if err := s.db.Find("PubKey", pubKey, otpk); err != nil {
		return nil, err
	}
	if otpk == nil {
		return nil, nil
	}

	// decrypt private key and copy over
	plainPrivKey, err := s.km.AESDecrypt(otpk.PrivKey)
	if err != nil {
		return nil, err
	}
	if len(plainPrivKey) != 32 {
		return nil, errors.New("got invalid private key with len != 32")
	}
	copy(privKey[:], plainPrivKey)

	return privKey, s.db.DeleteStruct(otpk)

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

package db

import (
	"time"
	"errors"

	aes "github.com/Bit-Nation/panthalassa/crypto/aes"
	keyManager "github.com/Bit-Nation/panthalassa/keyManager"
	x3dh "github.com/Bit-Nation/x3dh"
	storm "github.com/asdine/storm"
	sq "github.com/asdine/storm/q"
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
	All() ([]*x3dh.KeyPair, error)
}

type SignedPreKey struct {
	ValidTill  int64           `storm:"index"`
	EncryptedPrivateKey aes.CipherText
	privateKey x3dh.PrivateKey
	PublicKey  x3dh.PublicKey  `storm:"index,id"`
	Version    uint
}

func (s *SignedPreKey) PrivateKey() x3dh.PrivateKey {
	return s.privateKey
}

type BoltSignedPreKeyStorage struct {
	db *storm.DB
	km *keyManager.KeyManager
}

func NewBoltSignedPreKeyStorage(db *storm.DB, km *keyManager.KeyManager) *BoltSignedPreKeyStorage {
	return &BoltSignedPreKeyStorage{
		db: db,
		km: km,
	}
}

func (s *BoltSignedPreKeyStorage) Put(signedPreKey x3dh.KeyPair) error {
	
	privKeyCT, err := s.km.AESEncrypt(signedPreKey.PrivateKey[:])
	if err != nil {
		return err
	}
	
	return s.db.Save(&SignedPreKey{
		ValidTill:  time.Now().Add(SignedPreKeyValidTimeFrame).Unix(),
		EncryptedPrivateKey: privKeyCT,
		PublicKey:  signedPreKey.PublicKey,
		Version:    1,
	})
}

func (s *BoltSignedPreKeyStorage) Get(publicKey x3dh.PublicKey) (*x3dh.PrivateKey, error) {
	
	// check if signed pre key exist
	q := s.db.Select(sq.Eq("PublicKey", publicKey))
	amount, err := q.Count(&SignedPreKey{})
	if err != nil {
		return nil, err
	}
	if amount == 0 {
		return nil, nil
	}
	
	var spk SignedPreKey
	if err := q.First(&spk); err != nil {
		return nil, err
	}
	plainPrivKey, err := s.km.AESDecrypt(spk.EncryptedPrivateKey)
	if len(plainPrivKey) != 32 {
		return nil, errors.New("received invalid private key (len != 32)")
	}
	if err != nil {
		return nil, err
	}
	privKey := &x3dh.PrivateKey{}
	copy(privKey[:], plainPrivKey)
	
	return privKey, nil
}

func (s *BoltSignedPreKeyStorage) All() ([]*x3dh.KeyPair, error) {

	signedPreKeys := []*x3dh.KeyPair{}

	var persistedSignedPreKeys []SignedPreKey
	if err := s.db.All(&persistedSignedPreKeys); err != nil {
		return nil, err
	}
	
	for _, signedPreKey := range persistedSignedPreKeys {
		priv, err := s.Get(signedPreKey.PublicKey)
		if err != nil {
			return nil, err
		}
		signedPreKeys = append(signedPreKeys, &x3dh.KeyPair{
			PublicKey: signedPreKey.PublicKey,
			PrivateKey: *priv,
		})
	}
	
	return signedPreKeys, nil

}

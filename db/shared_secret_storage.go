package db

import (
	"errors"
	"time"

	aes "github.com/Bit-Nation/panthalassa/crypto/aes"
	keyManager "github.com/Bit-Nation/panthalassa/keyManager"
	x3dh "github.com/Bit-Nation/x3dh"
	storm "github.com/asdine/storm"
	sq "github.com/asdine/storm/q"
	ed25519 "golang.org/x/crypto/ed25519"
)

type SharedSecret struct {
	// this prop is filled in during saving the shared secret
	X3dhSS                aes.CipherText
	x3dhSS                x3dh.SharedSecret
	Accepted              bool              `storm:"index"`
	CreatedAt             time.Time         `storm:"index"`
	DestroyAt             *time.Time        `storm:"index"`
	Partner               ed25519.PublicKey `storm:"index"`
	ID                    []byte            `storm:"index"`
	DBID                  int               `storm:"id,increment"`
	UsedOneTimePreKey     *x3dh.PublicKey
	UsedSignedPreKey      x3dh.PublicKey
	EphemeralKey          x3dh.PublicKey
	EphemeralKeySignature []byte
}

// @todo I don't like this solution. However, we still need to figure out
// @todo how to ignore fields from storm
func (ss *SharedSecret) GetX3dhSecret() x3dh.SharedSecret {
	return ss.x3dhSS
}
func (ss *SharedSecret) SetX3dhSecret(secret x3dh.SharedSecret) {
	ss.x3dhSS = secret
}

type SharedSecretStorage interface {
	HasAny(key ed25519.PublicKey) (bool, error)
	// must return an error if no shared secret found
	GetYoungest(key ed25519.PublicKey) (*SharedSecret, error)
	Put(ss SharedSecret) error
	// accept will mark the given shared secret as accepted
	// and will set a destroy date for all other shared secrets
	Accept(sharedSec SharedSecret) error
	// get sender public key and shared secret id
	Get(key ed25519.PublicKey, sharedSecretID []byte) (*SharedSecret, error)
}

func NewBoltSharedSecretStorage(db storm.Node, km *keyManager.KeyManager) *BoltSharedSecretStorage {
	return &BoltSharedSecretStorage{
		db: db,
		km: km,
	}
}

type BoltSharedSecretStorage struct {
	db storm.Node
	km *keyManager.KeyManager
}

func (b *BoltSharedSecretStorage) HasAny(partner ed25519.PublicKey) (bool, error) {

	q := b.db.Select(sq.Eq("Partner", partner))

	// count shared secrets
	amount, err := q.Count(&SharedSecret{})
	if err != nil {
		return false, err
	}

	if amount > 0 {
		return true, nil
	}
	return false, nil

}

func (b *BoltSharedSecretStorage) GetYoungest(partner ed25519.PublicKey) (*SharedSecret, error) {

	var shSec SharedSecret

	// fetch shared secret
	q := b.db.Select(sq.Eq("Partner", partner))
	amount, err := q.OrderBy("CreatedAt").Count(&shSec)
	if err != nil {
		return nil, err
	}

	if amount == 0 {
		return nil, nil
	}

	// fetch first shared secret
	if err := q.OrderBy("CreatedAt").Reverse().First(&shSec); err != nil {
		return nil, err
	}

	// decrypt shared secret
	plainSharedSec, err := b.km.AESDecrypt(shSec.X3dhSS)
	if err != nil {
		return nil, err
	}
	if len(plainSharedSec) != 32 {
		return nil, errors.New("got invalid plain shared secret with len != 32")
	}
	copy(shSec.x3dhSS[:], plainSharedSec)

	return &shSec, nil

}

func (b *BoltSharedSecretStorage) Put(ss SharedSecret) error {

	if len(ss.ID) != 32 {
		return errors.New("can't persisted shared secret with id len != 32")
	}

	if len(ss.Partner) != 32 {
		return errors.New("chat partner must have a length of 32")
	}

	if ss.x3dhSS == [32]byte{} {
		return errors.New("can't persist empty shared secret")
	}

	var err error
	// encrypt shared secret
	ss.X3dhSS, err = b.km.AESEncrypt(ss.x3dhSS[:])
	if err != nil {
		return err
	}

	return b.db.Save(&ss)

}

func (b *BoltSharedSecretStorage) Accept(sharedSec SharedSecret) error {
	q := b.db.Select(sq.And(
		sq.Eq("Partner", sharedSec.Partner),
		sq.Eq("ID", sharedSec.ID),
	))

	var ss SharedSecret
	if err := q.First(&ss); err != nil {
		return err
	}

	ss.Accepted = true

	return b.db.Update(&ss)
}

func (b *BoltSharedSecretStorage) Get(partner ed25519.PublicKey, id []byte) (*SharedSecret, error) {

	var ss SharedSecret

	// check if shared secret exist
	q := b.db.Select(sq.And(sq.Eq("Partner", partner), sq.Eq("ID", id)))
	amount, err := q.Count(&SharedSecret{})
	if err != nil {
		return nil, err
	}
	if amount == 0 {
		return nil, nil
	}

	if err := q.First(&ss); err != nil {
		return nil, err
	}

	// decrypt plain shared secret
	plainSharedSec, err := b.km.AESDecrypt(ss.X3dhSS)
	if err != nil {
		return nil, err
	}
	if len(plainSharedSec) != 32 {
		return nil, errors.New("invalid plain shared secret with len != 32")
	}
	copy(ss.x3dhSS[:], plainSharedSec)

	return &ss, err
}

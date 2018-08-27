package db

import (
	"errors"
	"time"

	"fmt"
	aes "github.com/Bit-Nation/panthalassa/crypto/aes"
	keyManager "github.com/Bit-Nation/panthalassa/keyManager"
	x3dh "github.com/Bit-Nation/x3dh"
	storm "github.com/asdine/storm"
	sq "github.com/asdine/storm/q"
	ed25519 "golang.org/x/crypto/ed25519"
)

type SharedSecret struct {
	// @todo figure out how to ignore fields in storm
	X3dhSS    x3dh.SharedSecret
	x3dhSS    aes.CipherText
	Accepted  bool              `storm:"index"`
	CreatedAt time.Time         `storm:"index"`
	DestroyAt *time.Time        `storm:"index"`
	Partner   ed25519.PublicKey `storm:"index"`
	ID        []byte            `storm:"index"`
}

// the persistedSharedSecret is almost the same as SharedSecret except for
// that the X3dhSS value is an AES cipher text.
type persistedSharedSecret struct {
	SharedSecret
	X3dhSS aes.CipherText `json:"x3dh_shared_secret"`
}

type SharedSecretStorage interface {
	HasAny(key ed25519.PublicKey) (bool, error)
	// must return an error if no shared secret found
	GetYoungest(key ed25519.PublicKey) (*SharedSecret, error)
	Put(ss SharedSecret) error
	// accept will mark the given shared secret as accepted
	// and will set a destroy date for all other shared secrets
	Accept(partner ed25519.PublicKey, sharedSec *SharedSecret) error
	// get sender public key and shared secret id
	Get(key ed25519.PublicKey, sharedSecretID []byte) (*SharedSecret, error)
}

func NewBoltSharedSecretStorage(db *storm.DB, km *keyManager.KeyManager) *BoltSharedSecretStorage {
	return &BoltSharedSecretStorage{
		db: db,
		km: km,
	}
}

type BoltSharedSecretStorage struct {
	db *storm.DB
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
	shSec := new(SharedSecret)
	shSec = nil

	// fetch shared secret
	q := b.db.Select(sq.Eq("Partner", partner))
	err := q.OrderBy("CreatedAt").First(shSec)
	if err != nil {
		return nil, err
	}
	if shSec == nil {
		return nil, nil
	}

	// decrypt shared secret
	plainSharedSec, err := b.km.AESDecrypt(shSec.x3dhSS)
	if err != nil {
		return nil, err
	}
	if len(plainSharedSec) != 32 {
		return nil, errors.New("got invalid plain shared secret with len != 32")
	}
	copy(shSec.X3dhSS[:], plainSharedSec)

	return shSec, nil

}

func (b *BoltSharedSecretStorage) Put(ss SharedSecret) error {

	if len(ss.ID) != 32 {
		return errors.New("can't persisted shared secret with id len != 32")
	}

	if len(ss.Partner) != 32 {
		return errors.New("chat partner must have a length of 32")
	}

	if ss.X3dhSS == [32]byte{} {
		return errors.New("can't persist empty shared secret")
	}

	var err error

	// encrypt shared secret
	ss.x3dhSS, err = b.km.AESEncrypt(ss.X3dhSS[:])
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

	ss := new(SharedSecret)
	if err := q.First(&ss); err != nil {
		return err
	}
	if ss == nil {
		return fmt.Errorf("couldn't find shared secret for partner %x with id %x", sharedSec.Partner, sharedSec.ID)
	}

	ss.Accepted = true

	return b.db.Update(ss)
}

func (b *BoltSharedSecretStorage) Get(partner ed25519.PublicKey, id []byte) (*SharedSecret, error) {
	ss := new(SharedSecret)
	ss = nil

	// fetch shared secret
	q := b.db.Select(sq.And(sq.Eq("Partner", partner), sq.Eq("ID", id)))
	if err := q.First(&ss); err != nil {
		return nil, err
	}
	if ss == nil {
		return nil, nil
	}

	// decrypt plain shared secret
	plainSharedSec, err := b.km.AESDecrypt(ss.x3dhSS)
	if err != nil {
		return nil, err
	}
	if len(plainSharedSec) != 32 {
		return nil, errors.New("invalid plain shared secret with len != 32")
	}

	return ss, err
}

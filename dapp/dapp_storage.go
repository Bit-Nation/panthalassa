package dapp

import (
	"encoding/hex"
	"errors"
	"fmt"

	uiapi "github.com/Bit-Nation/panthalassa/uiapi"
	storm "github.com/asdine/storm"
	sq "github.com/asdine/storm/q"
	ed25519 "golang.org/x/crypto/ed25519"
)

type Storage interface {
	SaveDApp(dApp Data) error
	All() ([]*Data, error)
	Get(signingKey ed25519.PublicKey) (*Data, error)
}

type BoltDAppStorage struct {
	db    *storm.DB
	uiApi *uiapi.Api
}

func NewDAppStorage(db *storm.DB, api *uiapi.Api) *BoltDAppStorage {
	return &BoltDAppStorage{
		db:    db,
		uiApi: api,
	}
}

func (s *BoltDAppStorage) SaveDApp(dApp Data) error {

	if dApp.Version < 1 {
		return errors.New("version must be at least 1")
	}

	// validate DApp signature
	valid, err := dApp.VerifySignature()
	if err != nil {
		return err
	}
	if !valid {
		return fmt.Errorf("invalid signature for DApp: %x", dApp.UsedSigningKey)
	}

	// save DApp
	if err := s.db.Save(&dApp); err != nil {
		return err
	}

	// update ui
	s.uiApi.Send("DAPP:PERSISTED", map[string]interface{}{
		"dapp_signing_key": hex.EncodeToString(dApp.UsedSigningKey),
	})

	return nil
}

func (s *BoltDAppStorage) All() ([]*Data, error) {
	var dApps []*Data
	return dApps, s.db.All(&dApps)
}

func (s *BoltDAppStorage) Get(signingKey ed25519.PublicKey) (*Data, error) {

	q := s.db.Select(sq.Eq("UsedSigningKey", signingKey))

	amount, err := q.Count(&Data{})
	if err != nil {
		return nil, err
	}
	if amount == 0 {
		return nil, nil
	}

	var dApp Data
	return &dApp, q.First(&dApp)

}

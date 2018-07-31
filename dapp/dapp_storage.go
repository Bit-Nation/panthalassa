package dapp

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"

	uiapi "github.com/Bit-Nation/panthalassa/uiapi"
	bolt "github.com/coreos/bbolt"
	"golang.org/x/crypto/ed25519"
)

var (
	dAppStoreBucketName = []byte("dapps")
)

type Storage interface {
	SaveDApp(dApp Data) error
	All() ([]*Data, error)
	Get(signingKey ed25519.PublicKey) (*Data, error)
}

type BoltDAppStorage struct {
	db    *bolt.DB
	uiApi *uiapi.Api
}

func NewDAppStorage(db *bolt.DB, api *uiapi.Api) *BoltDAppStorage {
	return &BoltDAppStorage{
		db:    db,
		uiApi: api,
	}
}

func (s *BoltDAppStorage) SaveDApp(dApp Data) error {
	return s.db.Update(func(tx *bolt.Tx) error {

		if dApp.Version < 1 {
			return errors.New("version must be at least 1")
		}

		tx.OnCommit(func() {
			s.uiApi.Send("DAPP:PERSISTED", map[string]interface{}{
				"dapp_signing_key": hex.EncodeToString(dApp.UsedSigningKey),
			})
		})

		valid, err := dApp.VerifySignature()
		if err != nil {
			return err
		}
		if !valid {
			return fmt.Errorf("invalid signature for DApp: %x", dApp.UsedSigningKey)
		}

		// fetch dApp storage bucket
		dAppStorageBucket, err := tx.CreateBucketIfNotExists(dAppStoreBucketName)
		if err != nil {
			return err
		}

		// marshal dApp
		rawDApp, err := json.Marshal(dApp)
		if err != nil {
			return err
		}

		// persist dApp
		return dAppStorageBucket.Put(dApp.UsedSigningKey, rawDApp)

	})
}

func (s *BoltDAppStorage) All() ([]*Data, error) {

	var dApps []*Data

	err := s.db.View(func(tx *bolt.Tx) error {

		// fetch dApp's bucket
		dAppStorage := tx.Bucket(dAppStoreBucketName)
		if dAppStorage == nil {
			return nil
		}

		return dAppStorage.ForEach(func(_, rawDApp []byte) error {

			// unmarshal build
			d := Data{}
			if err := json.Unmarshal(rawDApp, &d); err != nil {
				return err
			}

			// add to list of Dapps
			dApps = append(dApps, &d)

			return nil

		})

	})

	return dApps, err
}

func (s *BoltDAppStorage) Get(signingKey ed25519.PublicKey) (*Data, error) {

	var dApp *Data
	dApp = nil

	err := s.db.View(func(tx *bolt.Tx) error {
		var dApp *Data
		dApp = nil
		err := s.db.View(func(tx *bolt.Tx) error {
			d := Data{}
			if dAppStorage := tx.Bucket(dAppStoreBucketName); dAppStorage != nil {
				if rawDAppData := dAppStorage.Get(signingKey); rawDAppData != nil {
					if err := json.Unmarshal(rawDAppData, &d); err != nil {
						return err
					}
				}
			}
			dApp = &d
			return nil
		})
		return err
	})

	return dApp, err

}

package registry

import (
	"encoding/hex"
	"encoding/json"
	"testing"
	"time"

	dapp "github.com/Bit-Nation/panthalassa/dapp"
	keyManager "github.com/Bit-Nation/panthalassa/keyManager"
	keyStore "github.com/Bit-Nation/panthalassa/keyStore"
	mnemonic "github.com/Bit-Nation/panthalassa/mnemonic"
	require "github.com/stretchr/testify/require"
	ed25519 "golang.org/x/crypto/ed25519"
)

const testDApp = `{"name":{"en-us":"DApp Name"},"engine":"0.1.0","image":"aW1hZ2U=","used_signing_key":"ff1fd817be47bfe6d3e055dcbe62447069b86c698132e782bbba2e70124b5448","code":"var i = 1","version":"1","signature":"b79caecf98224430777b54d8bd2f2a209548d64ce62dda8d8f0cbd1d9df8ac1750c7b0a86563b978890bc19719dbfc700b78c66f3ed5c5df40e5f60190dbd307"}`

type memDAppStorage struct {
	saveDApp func(dApp dapp.Data) error
	all      func() ([]*dapp.Data, error)
	get      func(signingKey ed25519.PublicKey) (*dapp.Data, error)
}

func (s *memDAppStorage) SaveDApp(dApp dapp.Data) error {
	return s.saveDApp(dApp)
}

func (s *memDAppStorage) All() ([]*dapp.Data, error) {
	return s.all()
}

func (s *memDAppStorage) Get(signingKey ed25519.PublicKey) (*dapp.Data, error) {
	return s.get(signingKey)
}

func TestRegistry_StartDApp(t *testing.T) {

	// signing key
	signingKey, err := hex.DecodeString("ff1fd817be47bfe6d3e055dcbe62447069b86c698132e782bbba2e70124b5448")
	require.Nil(t, err)

	// key manager
	mne, err := mnemonic.New()
	require.Nil(t, err)
	ks, err := keyStore.NewFromMnemonic(mne)
	require.Nil(t, err)
	km := keyManager.CreateFromKeyStore(ks)

	dAppStorage := memDAppStorage{
		get: func(signingKey ed25519.PublicKey) (*dapp.Data, error) {
			rawDApp := dapp.RawData{}
			require.Nil(t, json.Unmarshal([]byte(testDApp), &rawDApp))
			dAppData, err := dapp.ParseJsonToData(rawDApp)
			return &dAppData, err
		},
		saveDApp: func(dApp dapp.Data) error {
			return nil
		},
	}

	reg, err := NewDAppRegistry(nil, Config{}, nil, km, &dAppStorage, nil, nil)
	require.Nil(t, err)
	require.Nil(t, reg.StartDApp(signingKey, time.Second*2))

}

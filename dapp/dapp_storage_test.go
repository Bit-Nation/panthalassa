package dapp

import (
	"crypto/rand"
	"os"
	"path/filepath"
	"testing"
	"time"

	"fmt"
	uiApi "github.com/Bit-Nation/panthalassa/uiapi"
	bolt "github.com/coreos/bbolt"
	"github.com/segmentio/objconv/json"
	require "github.com/stretchr/testify/require"
	ed25519 "golang.org/x/crypto/ed25519"
)

type testUpstream struct {
}

func (u testUpstream) Send(string string) {

}

func TestBoltStorage_SaveDApp(t *testing.T) {

	db := createDB()

	dAppStorage := BoltDAppStorage{
		db:    db,
		uiApi: uiApi.New(&testUpstream{}),
	}

	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	require.Nil(t, err)

	dAppJson := JsonBuild{
		Name: map[string]string{
			"en-us": "send and request money",
			"de":    "sende und fordere geld an",
		},
		UsedSigningKey: pub,
		Code:           []byte(`var wallet = "0x930aa9a843266bdb02847168d571e7913907dd84"`),
		Image:          []byte("base64..."),
		Engine:         SV{1, 2, 3},
		Version:        1,
	}

	// hash the DApp
	dAppHash, err := dAppJson.Hash()
	require.Nil(t, err)

	// add signature to Dapp
	dAppJson.Signature = ed25519.Sign(priv, dAppHash)

	// persist dApp
	require.Nil(t, dAppStorage.SaveDApp(dAppJson))

	err = db.View(func(tx *bolt.Tx) error {

		// dApp storage
		dAppStorage := tx.Bucket(dAppStoreBucketName)
		require.NotNil(t, dAppStorage)

		// raw DApp
		rawDApp := dAppStorage.Get(pub)
		require.NotNil(t, rawDApp)

		// make sure that the dApps are the same
		// since we persisted the whole Dapp
		dApp := JsonBuild{}
		require.Nil(t, json.Unmarshal(rawDApp, &dApp))
		require.Equal(t, dAppJson, dApp)

		return nil

	})
	require.Nil(t, err)

}

func TestBoltStorage_SaveDAppInvalidSignature(t *testing.T) {

	db := createDB()

	dAppStorage := BoltDAppStorage{
		db:    db,
		uiApi: uiApi.New(&testUpstream{}),
	}

	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	require.Nil(t, err)

	dAppJson := JsonBuild{
		Name: map[string]string{
			"en-us": "send and request money",
			"de":    "sende und fordere geld an",
		},
		UsedSigningKey: pub,
		Code:           []byte(`var wallet = "0x930aa9a843266bdb02847168d571e7913907dd84"`),
		Image:          []byte("base64..."),
		Engine:         SV{1, 2, 3},
		Version:        1,
	}

	// hash the DApp
	dAppHash, err := dAppJson.Hash()
	require.Nil(t, err)

	// add signature to Dapp
	dAppJson.Signature = ed25519.Sign(priv, dAppHash)
	// fake the signature here to make sure the verification will fail
	dAppJson.Signature[3] = 0xf3

	// persist dApp
	err = dAppStorage.SaveDApp(dAppJson)
	require.EqualError(t, err, fmt.Sprintf("invalid signature for DApp: %x", pub))

}

func TestBoltStorage_All(t *testing.T) {

	db := createDB()

	dAppStorage := BoltDAppStorage{
		db:    db,
		uiApi: uiApi.New(&testUpstream{}),
	}

	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	require.Nil(t, err)

	dAppJson := JsonBuild{
		Name: map[string]string{
			"en-us": "send and request money",
			"de":    "sende und fordere geld an",
		},
		UsedSigningKey: pub,
		Code:           []byte(`var wallet = "0x930aa9a843266bdb02847168d571e7913907dd84"`),
		Image:          []byte("base64..."),
		Engine:         SV{1, 2, 3},
		Version:        1,
	}

	// hash the DApp
	dAppHash, err := dAppJson.Hash()
	require.Nil(t, err)

	// add signature to DApp
	dAppJson.Signature = ed25519.Sign(priv, dAppHash)

	// persist DApp
	require.Nil(t, dAppStorage.SaveDApp(dAppJson))

	allDapps, err := dAppStorage.All()
	require.Nil(t, err)

	// make sure the first DApp is the DApp we persisted
	require.Equal(t, dAppJson, *allDapps[0])

}

func createDB() *bolt.DB {
	dbPath, err := filepath.Abs(os.TempDir() + "/" + time.Now().String())
	if err != nil {
		panic(err)
	}
	db, err := bolt.Open(dbPath, 0600, &bolt.Options{Timeout: time.Second})
	if err != nil {
		panic(err)
	}
	return db
}

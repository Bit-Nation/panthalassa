package db

import (
	km "github.com/Bit-Nation/panthalassa/keyManager"
	ks "github.com/Bit-Nation/panthalassa/keyStore"
	mnemonic "github.com/Bit-Nation/panthalassa/mnemonic"
	bolt "github.com/coreos/bbolt"
	"os"
	"time"
	"path/filepath"
)

func createKeyManager() *km.KeyManager {

	mne, err := mnemonic.New()
	if err != nil {
		panic(err)
	}

	keyStore, err := ks.NewFromMnemonic(mne)
	if err != nil {
		panic(err)
	}

	return km.CreateFromKeyStore(keyStore)

}

func createDB() *bolt.DB {
	dbPath, err := filepath.Abs(os.TempDir()+"/"+time.Now().String())
	if err != nil {
		panic(err)
	}
	db, err := bolt.Open(dbPath, 0600, &bolt.Options{Timeout: time.Second})
	if err != nil {
		panic(err)
	}
	return db
}

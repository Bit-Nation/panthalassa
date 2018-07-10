package db

import (
	km "github.com/Bit-Nation/panthalassa/keyManager"
	ks "github.com/Bit-Nation/panthalassa/keyStore"
	mnemonic "github.com/Bit-Nation/panthalassa/mnemonic"
	bolt "github.com/coreos/bbolt"
	"os"
	"time"
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
	db, err := bolt.Open(os.TempDir()+time.Now().String(), 0600, &bolt.Options{Timeout: time.Second})
	if err != nil {
		panic(err)
	}
	return db
}

package db

import (
	"crypto/rand"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	bolt "github.com/coreos/bbolt"
	require "github.com/stretchr/testify/require"
	ed25519 "golang.org/x/crypto/ed25519"
)

func createDB() *bolt.DB {
	dbPath, err := filepath.Abs(os.TempDir() + "/" + fmt.Sprint(time.Now().UnixNano()))
	if err != nil {
		panic(err)
	}
	db, err := bolt.Open(dbPath, 0600, &bolt.Options{Timeout: time.Second})
	if err != nil {
		panic(err)
	}
	return db
}

func TestBoltStorage_CRUD(t *testing.T) {

	db := createDB()
	pub, _, err := ed25519.GenerateKey(rand.Reader)
	require.Nil(t, err)

	s, err := NewBoltStorage(db, pub)
	require.Nil(t, err)

	// put value into storage
	require.Nil(t, s.Put([]byte("key"), []byte("value")))

	// make sure the right bucket got created
	err = db.View(func(tx *bolt.Tx) error {

		// fetch dApp bucket
		dAppBucket := tx.Bucket(dAppDBBucketName)
		require.NotNil(t, dAppDBBucketName)

		// dApp DB
		// bucket must not be null
		dAppDB := dAppBucket.Bucket(pub)
		require.NotNil(t, dAppDB)

		return nil
	})
	require.Nil(t, err)

	// key should exist in storage
	exist, err := s.Has([]byte("key"))
	require.Nil(t, err)
	require.True(t, exist)

	// value of key should exist
	value, err := s.Get([]byte("key"))
	require.Nil(t, err)
	require.Equal(t, []byte("value"), value)

	// delete key
	require.Nil(t, s.Delete([]byte("key")))

	// should not exist
	exist, err = s.Has([]byte("key"))
	require.Nil(t, err)
	require.False(t, exist)

}

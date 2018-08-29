package db

import (
	"crypto/rand"
	"os"
	"path/filepath"
	"testing"
	"time"

	storm "github.com/asdine/storm"
	require "github.com/stretchr/testify/require"
	ed25519 "golang.org/x/crypto/ed25519"
)

func createStorm() *storm.DB {
	dbPath, err := filepath.Abs(os.TempDir() + "/" + time.Now().String())
	if err != nil {
		panic(err)
	}
	db, err := storm.Open(dbPath)
	if err != nil {
		panic(err)
	}
	return db
}

func TestBoltStorage_CRUD(t *testing.T) {

	db := createStorm()
	pub, _, err := ed25519.GenerateKey(rand.Reader)
	require.Nil(t, err)

	s, err := NewBoltStorage(db, pub)
	require.Nil(t, err)

	// put value into storage
	require.Nil(t, s.Put([]byte("key"), []byte("value")))

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

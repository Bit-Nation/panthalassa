package db

import (
	"crypto/rand"
	"testing"
	"time"

	x3dh "github.com/Bit-Nation/x3dh"
	require "github.com/stretchr/testify/require"
	ed25519 "golang.org/x/crypto/ed25519"
)

func TestBoltSharedSecretStorage_Put(t *testing.T) {

	// setup
	db := createStorm()
	km := createKeyManager()
	storage := NewBoltSharedSecretStorage(db, km)
	pub, _, err := ed25519.GenerateKey(rand.Reader)
	require.Nil(t, err)

	ID := make([]byte, 32)
	_, err = rand.Read(ID)
	require.Nil(t, err)

	// persist shared secret with invalid id
	require.EqualError(t, storage.Put(SharedSecret{}), "can't persisted shared secret with id len != 32")

	// persist shared secret with invalid chat partner
	require.EqualError(t, storage.Put(SharedSecret{
		ID: ID,
	}), "chat partner must have a length of 32")

	// persist shared secret with invalid x3dh shared secret
	require.EqualError(t, storage.Put(SharedSecret{
		ID:      ID,
		Partner: pub,
	}), "can't persist empty shared secret")

	err = storage.Put(SharedSecret{
		ID:      ID,
		Partner: pub,
		x3dhSS:  x3dh.SharedSecret{1},
	})
	require.Nil(t, err)

	sharedSec, err := storage.Get(pub, ID)
	require.Nil(t, err)
	require.NotNil(t, sharedSec)

	require.Equal(t, x3dh.SharedSecret{1}, sharedSec.x3dhSS)

}

func TestBoltSharedSecretStorage_Accept(t *testing.T) {

	// setup
	db := createStorm()
	km := createKeyManager()
	storage := NewBoltSharedSecretStorage(db, km)
	pub, _, err := ed25519.GenerateKey(rand.Reader)
	require.Nil(t, err)

	ID := make([]byte, 32)
	_, err = rand.Read(ID)
	require.Nil(t, err)

	// persist shared secret
	require.Nil(t, storage.Put(SharedSecret{
		x3dhSS:  [32]byte{1, 2},
		ID:      ID,
		Partner: pub,
	}))

	// accept shared secret that doesn't exist
	err = storage.Accept(SharedSecret{
		Partner: make([]byte, 32),
		ID:      make([]byte, 32),
	})
	require.EqualError(t, err, "not found")

	// accept shared secret
	require.Nil(t, storage.Accept(SharedSecret{
		x3dhSS:  [32]byte{1, 2},
		ID:      ID,
		Partner: pub,
	}))

	// fetch shared secret
	ss, err := storage.Get(pub, ID)
	require.Nil(t, err)
	require.NotNil(t, ss)
	require.True(t, ss.Accepted)

}

func TestBoltSharedSecretStorage_Get(t *testing.T) {

	// setup
	db := createStorm()
	km := createKeyManager()
	storage := NewBoltSharedSecretStorage(db, km)
	pub, _, err := ed25519.GenerateKey(rand.Reader)
	require.Nil(t, err)

	ID := make([]byte, 32)
	_, err = rand.Read(ID)
	require.Nil(t, err)

	// persist shared secret
	require.Nil(t, storage.Put(SharedSecret{
		x3dhSS:  [32]byte{1, 2},
		ID:      ID,
		Partner: pub,
	}))

	// fetch shared secret
	ss, err := storage.Get(pub, ID)
	require.Nil(t, err)
	require.NotNil(t, ss)
	require.Equal(t, [32]byte{1, 2}, ss.x3dhSS)

	// shared secret should be nil if shared secret doesn't exist
	ss, err = storage.Get(pub, make([]byte, 32))
	require.Nil(t, err)
	require.Nil(t, ss)

}

func TestBoltSharedSecretStorage_GetYoungest(t *testing.T) {

	// setup
	db := createStorm()
	km := createKeyManager()
	storage := NewBoltSharedSecretStorage(db, km)
	pub, _, err := ed25519.GenerateKey(rand.Reader)
	require.Nil(t, err)

	ID := make([]byte, 32)
	_, err = rand.Read(ID)
	require.Nil(t, err)

	// persist first shared secret
	require.Nil(t, storage.Put(SharedSecret{
		x3dhSS:    [32]byte{1},
		CreatedAt: time.Now().Truncate(time.Minute),
		ID:        ID,
		Partner:   pub,
	}))

	ID = make([]byte, 32)
	_, err = rand.Read(ID)
	require.Nil(t, err)

	// persist second shared secret
	require.Nil(t, storage.Put(SharedSecret{
		x3dhSS:    [32]byte{1, 2},
		ID:        ID,
		CreatedAt: time.Now(),
		Partner:   pub,
	}))

	// fetch youngest shared secret
	ss, err := storage.GetYoungest(pub)
	require.Nil(t, err)
	require.NotNil(t, ss)

	require.Equal(t, [32]byte{1, 2}, ss.x3dhSS)

}

func TestBoltSharedSecretStorage_HasAny(t *testing.T) {

	// setup
	db := createStorm()
	km := createKeyManager()
	storage := NewBoltSharedSecretStorage(db, km)
	pub, _, err := ed25519.GenerateKey(rand.Reader)
	require.Nil(t, err)

	// should be false since we didn't persist a shared secret
	has, err := storage.HasAny(pub)
	require.Nil(t, err)
	require.False(t, has)

	// persist first shared secret
	require.Nil(t, storage.Put(SharedSecret{
		x3dhSS:    [32]byte{1, 2},
		ID:        make([]byte, 32),
		CreatedAt: time.Now().Truncate(time.Minute),
		Partner:   pub,
	}))

	// must be true since we persisted a shared secret
	has, err = storage.HasAny(pub)
	require.Nil(t, err)
	require.True(t, has)

}

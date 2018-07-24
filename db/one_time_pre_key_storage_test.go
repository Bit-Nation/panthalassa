package db

import (
	"crypto/rand"
	"testing"

	x3dh "github.com/Bit-Nation/x3dh"
	require "github.com/stretchr/testify/require"
)

func TestNewBoltOneTimePreKeyStorage(t *testing.T) {

	storage := NewBoltOneTimePreKeyStorage(createDB(), createKeyManager())
	c := x3dh.NewCurve25519(rand.Reader)

	// test key pair
	keyPair, err := c.GenerateKeyPair()
	require.Nil(t, err)

	// persist key pair
	require.Nil(t, storage.Put([]x3dh.KeyPair{keyPair}))

	// count key pairs (should be 1 since we persisted one)
	keys, err := storage.Count()
	require.Nil(t, err)
	require.Equal(t, uint32(1), keys)

	// make sure fetched pre key is the one we passed in
	privKey, err := storage.Cut(keyPair.PublicKey)
	require.Nil(t, err)
	require.Equal(t, keyPair.PrivateKey, *privKey)

	// must fail since the previous cut deleted the key from the storage
	_, err = storage.Cut(keyPair.PublicKey)
	require.EqualError(t, err, "failed to fetch one time pre key private key for given public key")
}

package db

import (
	"crypto/rand"
	"testing"

	x3dh "github.com/Bit-Nation/x3dh"
	require "github.com/stretchr/testify/require"
)

func TestBoltOneTimePreKeyStorage_Put(t *testing.T) {

	db := createStorm()

	storage := NewBoltOneTimePreKeyStorage(db, createKeyManager())
	c := x3dh.NewCurve25519(rand.Reader)

	// test key pair
	keyPair, err := c.GenerateKeyPair()
	require.Nil(t, err)

	// persist private key
	require.Nil(t, storage.Put([]x3dh.KeyPair{keyPair}))

	// fetch private key
	fetchedPriv, err := storage.Cut(keyPair.PublicKey[:])
	require.Nil(t, err)
	require.NotNil(t, fetchedPriv)

	// make sure private keys are equal
	require.Equal(t, keyPair.PrivateKey, *fetchedPriv)

}

func TestBoltOneTimePreKeyStorage_Count(t *testing.T) {

	db := createStorm()

	storage := NewBoltOneTimePreKeyStorage(db, createKeyManager())
	c := x3dh.NewCurve25519(rand.Reader)

	// test key pair
	keyPair, err := c.GenerateKeyPair()
	require.Nil(t, err)

	// persist private key
	require.Nil(t, storage.Put([]x3dh.KeyPair{keyPair}))

	// count all one time pre keys
	amount, err := storage.Count()
	require.Nil(t, err)
	require.Equal(t, uint32(1), amount)

}

func TestBoltOneTimePreKeyStorage_Cut(t *testing.T) {

	db := createStorm()

	storage := NewBoltOneTimePreKeyStorage(db, createKeyManager())
	c := x3dh.NewCurve25519(rand.Reader)

	// test key pair
	keyPair, err := c.GenerateKeyPair()
	require.Nil(t, err)

	// persist private key
	require.Nil(t, storage.Put([]x3dh.KeyPair{keyPair}))

	// cut private key
	fetchedPriv, err := storage.Cut(keyPair.PublicKey[:])
	require.Nil(t, err)

	// make sure private keys are equal
	require.Equal(t, keyPair.PrivateKey, *fetchedPriv)

	// cut again result should be nil
	priv, err := storage.Cut(keyPair.PublicKey[:])
	require.Nil(t, err)
	require.Nil(t, priv)

}

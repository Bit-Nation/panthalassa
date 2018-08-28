package db

import (
	"crypto/rand"
	"encoding/hex"
	"testing"

	x3dh "github.com/Bit-Nation/x3dh"
	require "github.com/stretchr/testify/require"
)

func TestBoltSignedPreKeyStorage_Put(t *testing.T) {

	// setup
	db := createStorm()
	km := createKeyManager()
	curve := x3dh.NewCurve25519(rand.Reader)
	signedPreKeyStorage := NewBoltSignedPreKeyStorage(db, km)

	// test key pair
	keyPair, err := curve.GenerateKeyPair()
	require.Nil(t, err)

	require.Nil(t, signedPreKeyStorage.Put(keyPair))

	privKey, err := signedPreKeyStorage.Get(keyPair.PublicKey)
	require.Nil(t, err)
	require.NotNil(t, privKey)

	require.Equal(t, &keyPair.PrivateKey, privKey)

}

func TestBoltSignedPreKeyStorage_Get(t *testing.T) {

	// setup
	db := createStorm()
	km := createKeyManager()
	curve := x3dh.NewCurve25519(rand.Reader)
	signedPreKeyStorage := NewBoltSignedPreKeyStorage(db, km)

	// test key pair
	keyPair, err := curve.GenerateKeyPair()
	require.Nil(t, err)

	require.Nil(t, signedPreKeyStorage.Put(keyPair))

	// get and make sure key is as expected
	privKey, err := signedPreKeyStorage.Get(keyPair.PublicKey)
	require.Nil(t, err)
	require.NotNil(t, privKey)
	require.Equal(t, &keyPair.PrivateKey, privKey)

	// try to fetch private key that doesn't exist
	privKey, err = signedPreKeyStorage.Get(x3dh.PublicKey{1, 2, 3})
	require.Nil(t, err)
	require.Nil(t, privKey)

}

func TestBoltSignedPreKeyStorage_All(t *testing.T) {

	// setup
	db := createStorm()
	km := createKeyManager()
	curve := x3dh.NewCurve25519(rand.Reader)
	signedPreKeyStorage := NewBoltSignedPreKeyStorage(db, km)

	// test key pair
	pairOne, err := curve.GenerateKeyPair()
	require.Nil(t, err)
	pairTwo, err := curve.GenerateKeyPair()
	require.Nil(t, err)

	// persist key pairs
	require.Nil(t, signedPreKeyStorage.Put(pairOne))
	require.Nil(t, signedPreKeyStorage.Put(pairTwo))

	// fetch private key
	pairOnePrivKey, err := signedPreKeyStorage.Get(pairOne.PublicKey)
	require.Nil(t, err)
	require.Equal(t, pairOne.PrivateKey, *pairOnePrivKey)

	// fetch all private keys
	keyPairs, err := signedPreKeyStorage.All()
	require.Nil(t, err)

	require.True(t, hex.EncodeToString(pairOne.PrivateKey[:]) == hex.EncodeToString(keyPairs[0].PrivateKey[:]) || hex.EncodeToString(pairOne.PrivateKey[:]) == hex.EncodeToString(keyPairs[1].PrivateKey[:]))
	require.True(t, hex.EncodeToString(pairTwo.PrivateKey[:]) == hex.EncodeToString(keyPairs[0].PrivateKey[:]) || hex.EncodeToString(pairTwo.PrivateKey[:]) == hex.EncodeToString(keyPairs[1].PrivateKey[:]))

}

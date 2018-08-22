package db

import (
	"crypto/rand"
	"encoding/hex"
	"testing"

	x3dh "github.com/Bit-Nation/x3dh"
	require "github.com/stretchr/testify/require"
)

func TestBoltSignedPreKeyStorage(t *testing.T) {

	// setup
	db := createDB()
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
	keyPairs := signedPreKeyStorage.All()

	require.Equal(t, hex.EncodeToString(pairOne.PrivateKey[:]), hex.EncodeToString(keyPairs[0].PrivateKey[:]))
	require.Equal(t, hex.EncodeToString(pairTwo.PrivateKey[:]), hex.EncodeToString(keyPairs[1].PrivateKey[:]))

}

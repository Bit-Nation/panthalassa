package db

import (
	"crypto/rand"
	"encoding/hex"
	"testing"

	preKey "github.com/Bit-Nation/panthalassa/chat/prekey"
	x3dh "github.com/Bit-Nation/x3dh"
	require "github.com/stretchr/testify/require"
)

func TestBoltUserStorage_PuSignedPreKey(t *testing.T) {

	b := createStorm()
	km := createKeyManager()
	userStorage := BoltUserStorage{
		db: b,
	}

	curve := x3dh.NewCurve25519(rand.Reader)

	// pre key
	keyPair, err := curve.GenerateKeyPair()
	require.Nil(t, err)
	signedPreKey := preKey.PreKey{}
	signedPreKey.PrivateKey = keyPair.PrivateKey
	signedPreKey.PublicKey = keyPair.PublicKey
	require.Nil(t, signedPreKey.Sign(*km))

	// id pub key
	idPubKeyStr, err := km.IdentityPublicKey()
	require.Nil(t, err)
	idPubKey, err := hex.DecodeString(idPubKeyStr)
	require.Nil(t, err)

	// persist signed pre key
	require.Nil(t, userStorage.PutSignedPreKey(idPubKey, signedPreKey))

	// fetch signed pre key
	pk, err := userStorage.GetSignedPreKey(idPubKey)
	require.Nil(t, err)

	require.Equal(t, signedPreKey.PublicKey, pk.PublicKey)

}

func TestBoltUserStorage_GetSignedPreKey(t *testing.T) {

	db := createStorm()
	km := createKeyManager()
	userStorage := BoltUserStorage{
		db: db,
	}

	curve := x3dh.NewCurve25519(rand.Reader)

	// pre key
	keyPair, err := curve.GenerateKeyPair()
	require.Nil(t, err)
	signedPreKey := preKey.PreKey{}
	signedPreKey.PrivateKey = keyPair.PrivateKey
	signedPreKey.PublicKey = keyPair.PublicKey
	require.Nil(t, signedPreKey.Sign(*km))

	// id pub key
	idPubKeyStr, err := km.IdentityPublicKey()
	require.Nil(t, err)
	idPubKey, err := hex.DecodeString(idPubKeyStr)
	require.Nil(t, err)

	// persist signed pre key
	require.Nil(t, userStorage.PutSignedPreKey(idPubKey, signedPreKey))

	// fetch signed pre key
	fetchedSignedPreKey, err := userStorage.GetSignedPreKey(idPubKey)
	require.Nil(t, err)
	require.NotNil(t, signedPreKey)

	// verify signature
	valid, err := fetchedSignedPreKey.VerifySignature(idPubKey)
	require.Nil(t, err)
	require.True(t, valid)

}

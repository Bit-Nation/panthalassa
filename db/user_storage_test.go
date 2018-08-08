package db

import (
	"crypto/rand"
	"encoding/hex"
	"testing"

	preKey "github.com/Bit-Nation/panthalassa/chat/prekey"
	bpb "github.com/Bit-Nation/protobuffers"
	x3dh "github.com/Bit-Nation/x3dh"
	bolt "github.com/coreos/bbolt"
	proto "github.com/gogo/protobuf/proto"
	require "github.com/stretchr/testify/require"
)

func TestBoltUserStorage_PutSignedPreKey(t *testing.T) {

	b := createDB()
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

	err = b.View(func(tx *bolt.Tx) error {

		// user storage bucket
		userStorageBucket := tx.Bucket(userStorageBucketName)
		require.NotNil(t, userStorageBucketName)

		// user bucket based on pub key
		userBucket := userStorageBucket.Bucket(idPubKey)
		require.NotNil(t, userBucket)

		// raw signed pre key
		rawSignedPreKey := userBucket.Get(signedPreKeyName)
		require.NotNil(t, rawSignedPreKey)

		// unmarshal signed  pre key
		fetchedSignedPreKey := bpb.PreKey{}
		require.Nil(t, proto.Unmarshal(rawSignedPreKey, &fetchedSignedPreKey))

		// do assertions on the pre key
		require.Equal(t, signedPreKey.PublicKey[:], fetchedSignedPreKey.Key)

		return nil

	})
	require.Nil(t, err)

}

func TestBoltUserStorage_GetSignedPreKey(t *testing.T) {

	b := createDB()
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
	fetchedSignedPreKey, err := userStorage.GetSignedPreKey(idPubKey)
	require.Nil(t, err)
	require.NotNil(t, signedPreKey)

	// verify signature
	valid, err := fetchedSignedPreKey.VerifySignature(idPubKey)
	require.Nil(t, err)
	require.True(t, valid)

}

package db

import (
	"crypto/rand"
	"encoding/hex"
	"github.com/Bit-Nation/panthalassa/crypto/aes"
	"github.com/Bit-Nation/x3dh"
	"github.com/coreos/bbolt"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestPreKeyStorage_Put(t *testing.T) {

	db := createDB()
	km := createKeyManager()

	storage := PreKeyStorage{
		km: km,
		db: db,
	}

	// create key pair
	c := x3dh.NewCurve25519(rand.Reader)
	keyPair, err := c.GenerateKeyPair()
	require.Nil(t, err)

	require.Nil(t, storage.Put(keyPair))

	err = db.View(func(tx *bolt.Tx) error {

		// the pre key store bucket
		preKeyStore := tx.Bucket(preKeyStoreBucket)

		// encrypted private key
		encryptedPrivateKey := preKeyStore.Get(keyPair.PublicKey[:])

		// unmarshal private key
		cipherText, err := aes.Unmarshal(encryptedPrivateKey)
		require.Nil(t, err)

		// decrypt private key
		plainPrivateKey, err := km.AESDecrypt(cipherText)
		require.Nil(t, err)

		require.Equal(t, hex.EncodeToString(keyPair.PrivateKey[:]), hex.EncodeToString(plainPrivateKey))

		return nil
	})
	require.Nil(t, err)

}

// in the case there is no bucket for pre keys
// we expect false for exist and nil for error
func TestPreKeyStorage_HasWithoutBucket(t *testing.T) {

	db := createDB()
	km := createKeyManager()

	storage := PreKeyStorage{
		km: km,
		db: db,
	}

	// create key pair
	c := x3dh.NewCurve25519(rand.Reader)
	keyPair, err := c.GenerateKeyPair()
	require.Nil(t, err)

	exist, err := storage.Has(keyPair.PublicKey)
	require.Nil(t, err)
	require.False(t, exist)

}

func TestPreKeyStorage_HasNoPrivateKeyForPubKey(t *testing.T) {

	db := createDB()
	km := createKeyManager()

	// create bucket
	err := db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucket(preKeyStoreBucket)
		return err
	})
	require.Nil(t, err)

	storage := PreKeyStorage{
		km: km,
		db: db,
	}

	// create key pair
	c := x3dh.NewCurve25519(rand.Reader)
	keyPair, err := c.GenerateKeyPair()
	require.Nil(t, err)

	exist, err := storage.Has(keyPair.PublicKey)
	require.Nil(t, err)
	require.False(t, exist)

}

func TestPreKeyStorage_Has(t *testing.T) {

	db := createDB()
	km := createKeyManager()

	storage := PreKeyStorage{
		km: km,
		db: db,
	}

	// create key pair
	c := x3dh.NewCurve25519(rand.Reader)
	keyPair, err := c.GenerateKeyPair()
	require.Nil(t, err)

	// persist private key
	err = db.Update(func(tx *bolt.Tx) error {
		preKeyStorage, err := tx.CreateBucket(preKeyStoreBucket)
		require.Nil(t, err)

		return preKeyStorage.Put(keyPair.PublicKey[:], []byte("dummy_value"))
	})
	require.Nil(t, err)

	exist, err := storage.Has(keyPair.PublicKey)
	require.Nil(t, err)
	require.True(t, exist)

}

func TestPreKeyStorage_Get(t *testing.T) {
	db := createDB()
	km := createKeyManager()

	storage := PreKeyStorage{
		km: km,
		db: db,
	}

	// create key pair
	c := x3dh.NewCurve25519(rand.Reader)
	keyPair, err := c.GenerateKeyPair()
	require.Nil(t, err)

	// persist private key
	err = db.Update(func(tx *bolt.Tx) error {
		preKeyStorage, err := tx.CreateBucket(preKeyStoreBucket)
		require.Nil(t, err)

		encryptedPrivateKey, err := km.AESEncrypt(keyPair.PrivateKey[:])
		require.Nil(t, err)

		rawEncryptedPrivateKey, err := encryptedPrivateKey.Marshal()
		require.Nil(t, err)

		return preKeyStorage.Put(keyPair.PublicKey[:], rawEncryptedPrivateKey)
	})
	require.Nil(t, err)

	fetchedPrivateKey, err := storage.Get(keyPair.PublicKey)
	require.Nil(t, err)
	require.Equal(t, keyPair.PrivateKey, fetchedPrivateKey)

}

// get should return an error if there is no private key present
func TestPreKeyStorage_GetErrorIfNoKey(t *testing.T) {
	db := createDB()
	km := createKeyManager()

	storage := PreKeyStorage{
		km: km,
		db: db,
	}

	// create key pair
	c := x3dh.NewCurve25519(rand.Reader)
	keyPair, err := c.GenerateKeyPair()
	require.Nil(t, err)

	// create bucket
	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucket(preKeyStoreBucket)
		return err
	})
	require.Nil(t, err)

	fetchedPrivateKey, err := storage.Get(keyPair.PublicKey)
	require.EqualError(t, err, "couldn't find private key for given public key")
	require.Equal(t, x3dh.PrivateKey{}, fetchedPrivateKey)

}

func TestPreKeyStorage_Delete(t *testing.T) {

	db := createDB()
	km := createKeyManager()

	storage := PreKeyStorage{
		km: km,
		db: db,
	}

	// create key pair
	c := x3dh.NewCurve25519(rand.Reader)
	keyPair, err := c.GenerateKeyPair()
	require.Nil(t, err)

	// persist dummy value under public key
	err = db.Update(func(tx *bolt.Tx) error {
		preKeyStorage, err := tx.CreateBucket(preKeyStoreBucket)
		require.Nil(t, err)
		return preKeyStorage.Put(keyPair.PublicKey[:], []byte("I am a dummy value"))
	})
	require.Nil(t, err)

	// delete key
	require.Nil(t, storage.Delete(keyPair.PublicKey))

	// make sure key was really deleted
	err = db.Update(func(tx *bolt.Tx) error {
		preKeyStorage := tx.Bucket(preKeyStoreBucket)
		require.NotNil(t, preKeyStorage)
		require.Nil(t, err)
		require.Nil(t, preKeyStorage.Get(keyPair.PublicKey[:]))
		return nil
	})
	require.Nil(t, err)

}

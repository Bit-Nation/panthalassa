package db

import (
	"encoding/hex"
	"testing"

	aes "github.com/Bit-Nation/panthalassa/crypto/aes"
	bolt "github.com/coreos/bbolt"
	require "github.com/stretchr/testify/require"
	dr "github.com/tiabc/doubleratchet"
)

// hint, this tests use "keyPair.PublicKey()" as the "key"
// this is not used in real cases.

func TestStore_Put(t *testing.T) {

	km := createKeyManager()
	db := createDB()

	s := Store{
		db: db,
		km: km,
	}

	crypto := dr.DefaultCrypto{}
	keyPair, err := crypto.GenerateDH()
	require.Nil(t, err)

	s.Put(keyPair.PublicKey(), 3, keyPair.PrivateKey())

	err = db.View(func(tx *bolt.Tx) error {

		drKeyStore := tx.Bucket(doubleRatchetKeyStoreBucket)
		require.NotNil(t, drKeyStore)

		k := keyPair.PublicKey()
		messageStore := drKeyStore.Bucket(k[:])
		require.NotNil(t, messageStore)

		rawCipherText := messageStore.Get(uintToBytes(3))
		require.NotNil(t, rawCipherText)

		aesCipherText, err := aes.Unmarshal(rawCipherText)
		require.Nil(t, err)

		privKey, err := km.AESDecrypt(aesCipherText)
		require.Nil(t, err)

		pk := keyPair.PrivateKey()
		require.Equal(t, hex.EncodeToString(pk[:]), hex.EncodeToString(privKey))

		return nil
	})
	require.Nil(t, err)
}

func TestStore_Get(t *testing.T) {

	km := createKeyManager()
	db := createDB()

	s := Store{
		db: db,
		km: km,
	}

	crypto := dr.DefaultCrypto{}
	keyPair, err := crypto.GenerateDH()
	require.Nil(t, err)

	// persist private key
	s.db.Update(func(tx *bolt.Tx) error {
		// double ratchet key store bucket
		drKeyStore, err := tx.CreateBucketIfNotExists(doubleRatchetKeyStoreBucket)
		require.Nil(t, err)
		pubKey := keyPair.PublicKey()
		privKey := keyPair.PrivateKey()

		// message key store
		messageKeyStore, err := drKeyStore.CreateBucketIfNotExists(pubKey[:])
		require.Nil(t, err)

		// encrypt private key
		encryptedPrivKey, err := km.AESEncrypt(privKey[:])
		require.Nil(t, err)

		// raw encrypted private key
		rawEncryptedPrivKey, err := encryptedPrivKey.Marshal()
		require.Nil(t, err)

		messageKeyStore.Put(uintToBytes(4), rawEncryptedPrivKey)

		return nil
	})

	// should be false if key does not exist
	messageKey, exist := s.Get([32]byte{}, 0)
	require.False(t, exist)
	require.Equal(t, dr.Key{}, messageKey)

	// should be false if message number does not exist
	messageKey, exist = s.Get(keyPair.PublicKey(), 999)
	require.False(t, exist)
	require.Equal(t, dr.Key{}, messageKey)

	// should be true if key exist
	messageKey, exist = s.Get(keyPair.PublicKey(), 4)
	require.True(t, exist)
	require.Equal(t, keyPair.PrivateKey(), messageKey)

}

func TestStore_DeleteMk(t *testing.T) {

	km := createKeyManager()
	db := createDB()

	s := Store{
		db: db,
		km: km,
	}

	crypto := dr.DefaultCrypto{}
	keyPair, err := crypto.GenerateDH()
	pubKey := keyPair.PublicKey()
	require.Nil(t, err)

	// create message key
	s.db.Update(func(tx *bolt.Tx) error {
		// double ratchet key store bucket
		drKeyStore, err := tx.CreateBucketIfNotExists(doubleRatchetKeyStoreBucket)
		require.Nil(t, err)

		// message key store
		messageKeyStore, err := drKeyStore.CreateBucketIfNotExists(pubKey[:])
		require.Nil(t, err)

		messageKeyStore.Put(uintToBytes(2), []byte("hi"))

		return nil
	})

	// make sure message key exist
	s.db.View(func(tx *bolt.Tx) error {
		// double ratchet key store bucket
		drKeyStore := tx.Bucket(doubleRatchetKeyStoreBucket)
		require.NotNil(t, drKeyStore)

		// message key store
		messageKeyStore := drKeyStore.Bucket(pubKey[:])
		require.NotNil(t, messageKeyStore)

		require.Equal(t, []byte("hi"), messageKeyStore.Get(uintToBytes(2)))

		return nil
	})

	// delete message key
	s.DeleteMk(pubKey, 2)

	// make sure message key does not exist
	s.db.View(func(tx *bolt.Tx) error {
		// double ratchet key store bucket
		drKeyStore := tx.Bucket(doubleRatchetKeyStoreBucket)
		require.NotNil(t, drKeyStore)

		// message key store
		messageKeyStore := drKeyStore.Bucket(pubKey[:])
		require.NotNil(t, messageKeyStore)

		require.Nil(t, messageKeyStore.Get(uintToBytes(2)))

		return nil
	})

}

func TestStore_DeletePk(t *testing.T) {

	km := createKeyManager()
	db := createDB()

	s := Store{
		db: db,
		km: km,
	}

	crypto := dr.DefaultCrypto{}
	keyPair, err := crypto.GenerateDH()
	pubKey := keyPair.PublicKey()
	require.Nil(t, err)

	// create message keys bucket
	s.db.Update(func(tx *bolt.Tx) error {
		// double ratchet key store bucket
		drKeyStore, err := tx.CreateBucketIfNotExists(doubleRatchetKeyStoreBucket)
		require.Nil(t, err)

		// message key store
		_, err = drKeyStore.CreateBucketIfNotExists(pubKey[:])
		require.Nil(t, err)

		return nil
	})

	// make sure message key bucket exist
	s.db.View(func(tx *bolt.Tx) error {
		// double ratchet key store bucket
		drKeyStore := tx.Bucket(doubleRatchetKeyStoreBucket)
		require.NotNil(t, drKeyStore)

		// message key store
		messageKeyStore := drKeyStore.Bucket(pubKey[:])
		require.NotNil(t, messageKeyStore)

		return nil
	})

	// delete message key bucket
	s.DeletePk(pubKey)

	// make sure message key bucket does not exist
	s.db.View(func(tx *bolt.Tx) error {
		// double ratchet key store bucket
		drKeyStore := tx.Bucket(doubleRatchetKeyStoreBucket)
		require.NotNil(t, drKeyStore)

		// message key store
		messageKeyStore := drKeyStore.Bucket(pubKey[:])
		require.Nil(t, messageKeyStore)

		return nil
	})

}

func TestStore_Count(t *testing.T) {

	km := createKeyManager()
	db := createDB()

	s := Store{
		db: db,
		km: km,
	}

	crypto := dr.DefaultCrypto{}
	keyPair, err := crypto.GenerateDH()
	pubKey := keyPair.PublicKey()
	require.Nil(t, err)

	// persist a few test message keys
	s.db.Update(func(tx *bolt.Tx) error {
		// double ratchet key store bucket
		drKeyStore, err := tx.CreateBucketIfNotExists(doubleRatchetKeyStoreBucket)
		require.Nil(t, err)

		// message key store
		messageStore, err := drKeyStore.CreateBucketIfNotExists(pubKey[:])
		require.Nil(t, err)

		// persist a few message message
		messageStore.Put(uintToBytes(1), []byte(""))
		messageStore.Put(uintToBytes(3), []byte(""))
		messageStore.Put(uintToBytes(4), []byte(""))
		messageStore.Put(uintToBytes(199), []byte(""))

		return nil
	})

	// delete message key bucket
	require.Equal(t, uint(4), s.Count(pubKey))

}

func TestStore_AllSuccess(t *testing.T) {

	km := createKeyManager()
	db := createDB()

	s := Store{
		db: db,
		km: km,
	}

	crypto := dr.DefaultCrypto{}
	keyPairOne, err := crypto.GenerateDH()
	keyPairTwo, err := crypto.GenerateDH()
	require.Nil(t, err)

	// persist a few message keys
	prepare := func(keyPair dr.DHPair, msgKeyNumber uint) {
		db.Update(func(tx *bolt.Tx) error {

			drKeyStore, err := tx.CreateBucketIfNotExists(doubleRatchetKeyStoreBucket)
			require.Nil(t, err)

			pubKey := keyPair.PublicKey()

			// message key store
			messageStore, err := drKeyStore.CreateBucketIfNotExists(pubKey[:])
			require.Nil(t, err)

			// encrypt message key
			privKey := keyPair.PrivateKey()
			privKeyCipherText, err := km.AESEncrypt(privKey[:])
			require.Nil(t, err)

			// raw message key cipher text
			rawMsgKeyCt, err := privKeyCipherText.Marshal()
			require.Nil(t, err)

			require.Nil(t, messageStore.Put(uintToBytes(msgKeyNumber), rawMsgKeyCt))

			return nil

		})
	}

	prepare(keyPairOne, 2)
	prepare(keyPairTwo, 4)

	allKeys := s.All()

	require.Equal(t, keyPairOne.PrivateKey(), allKeys[keyPairOne.PublicKey()][2])
	require.Equal(t, keyPairTwo.PrivateKey(), allKeys[keyPairTwo.PublicKey()][4])

}

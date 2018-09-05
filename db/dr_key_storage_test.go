package db

import (
	"encoding/hex"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	storm "github.com/asdine/storm"
	sq "github.com/asdine/storm/q"
	require "github.com/stretchr/testify/require"
	dr "github.com/tiabc/doubleratchet"
)

func createStorm() *storm.DB {
	dbPath, err := filepath.Abs(os.TempDir() + "/" + strconv.Itoa(int(time.Now().UnixNano())))
	if err != nil {
		panic(err)
	}
	db, err := storm.Open(dbPath)
	if err != nil {
		panic(err)
	}
	return db
}

// hint, this tests use "keyPair.PublicKey()" as the "key"
// this is not used in real cases.

func TestStore_Put(t *testing.T) {

	km := createKeyManager()
	db := createStorm()

	s := NewBoltDRKeyStorage(db, km)

	crypto := dr.DefaultCrypto{}
	keyPair, err := crypto.GenerateDH()
	priv := keyPair.PrivateKey()
	require.Nil(t, err)

	s.Put(keyPair.PublicKey(), 3, keyPair.PrivateKey())

	q := db.Select(sq.And(
		sq.Eq("Key", keyPair.PublicKey()),
		sq.Eq("MsgNum", 3),
	))
	var drk DRKey
	require.Nil(t, q.First(&drk))

	// a few redundant check's - but mooooore is better
	require.Equal(t, keyPair.PublicKey(), drk.Key)
	require.Equal(t, uint(3), drk.MsgNum)

	// check the private key
	plainKey, err := km.AESDecrypt(drk.MsgKey)
	require.Nil(t, err)
	require.Equal(t, hex.EncodeToString(priv[:]), hex.EncodeToString(plainKey))

}

func TestStore_Get(t *testing.T) {

	km := createKeyManager()
	db := createStorm()

	s := NewBoltDRKeyStorage(db, km)

	crypto := dr.DefaultCrypto{}
	keyPair, err := crypto.GenerateDH()
	priv := keyPair.PrivateKey()
	require.Nil(t, err)

	encryptedPrivKey, err := km.AESEncrypt(priv[:])
	require.Nil(t, err)

	// persist test dr key
	err = db.Save(&DRKey{
		Key:    keyPair.PublicKey(),
		MsgNum: 4,
		MsgKey: encryptedPrivKey,
	})
	require.Nil(t, err)

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
	db := createStorm()

	s := NewBoltDRKeyStorage(db, km)

	crypto := dr.DefaultCrypto{}
	keyPair, err := crypto.GenerateDH()
	pubKey := keyPair.PublicKey()
	priv := keyPair.PrivateKey()
	require.Nil(t, err)

	ct, err := km.AESEncrypt(priv[:])
	require.Nil(t, err)

	// persist double ratchet key
	err = db.Save(&DRKey{
		Key:    pubKey,
		MsgNum: 2,
		MsgKey: ct,
	})
	require.Nil(t, err)

	// delete message key
	s.DeleteMk(pubKey, 2)

	// make sure message key does not exist
	q := db.Select()
	amount, err := q.Count(&DRKey{})
	require.Nil(t, err)
	require.Equal(t, 0, amount)

}

func TestStore_DeletePk(t *testing.T) {

	km := createKeyManager()
	db := createStorm()

	s := NewBoltDRKeyStorage(db, km)

	crypto := dr.DefaultCrypto{}
	keyPair, err := crypto.GenerateDH()
	pubKey := keyPair.PublicKey()
	require.Nil(t, err)

	// persist dr key
	err = db.Save(&DRKey{
		Key:    pubKey,
		MsgNum: 3,
	})
	require.Nil(t, err)

	// delete message key bucket
	s.DeletePk(pubKey)

	q := db.Select()
	amount, err := q.Count(&DRKey{})
	require.Nil(t, err)
	require.Equal(t, 0, amount)

}

func TestStore_Count(t *testing.T) {

	km := createKeyManager()
	db := createStorm()

	s := NewBoltDRKeyStorage(db, km)

	// persist a few test message keys
	saveTestKey := func(k dr.Key, msgNum uint) {
		err := db.Save(&DRKey{
			Key:    k,
			MsgNum: msgNum,
		})
		require.Nil(t, err)
	}

	// persist the first set of message keys
	saveTestKey(dr.Key{1, 2, 3}, 3)
	saveTestKey(dr.Key{1, 2, 3}, 4)
	saveTestKey(dr.Key{1, 2, 3}, 5)
	saveTestKey(dr.Key{1, 2, 3}, 6)

	// persist another set of messages with a different Key
	// to make sure count works correct
	saveTestKey(dr.Key{1}, 3)
	saveTestKey(dr.Key{1}, 4)

	// delete message key bucket
	require.Equal(t, uint(4), s.Count(dr.Key{1, 2, 3}))

}

func TestStore_AllSuccess(t *testing.T) {

	km := createKeyManager()
	db := createStorm()

	s := NewBoltDRKeyStorage(db, km)

	crypto := dr.DefaultCrypto{}
	keyPairOne, err := crypto.GenerateDH()
	require.Nil(t, err)
	keyPairTwo, err := crypto.GenerateDH()
	require.Nil(t, err)

	// persist the first key pair
	privKeyOne := keyPairOne.PrivateKey()
	ct, err := km.AESEncrypt(privKeyOne[:])
	require.Nil(t, err)
	err = db.Save(&DRKey{
		Key:    keyPairOne.PublicKey(),
		MsgNum: 2,
		MsgKey: ct,
	})
	require.Nil(t, err)

	// persist second key pair
	privKeyTwo := keyPairTwo.PrivateKey()
	ct, err = km.AESEncrypt(privKeyTwo[:])
	require.Nil(t, err)
	err = db.Save(&DRKey{
		Key:    keyPairTwo.PublicKey(),
		MsgNum: 4,
		MsgKey: ct,
	})
	require.Nil(t, err)

	allKeys := s.All()

	require.Equal(t, keyPairOne.PrivateKey(), allKeys[keyPairOne.PublicKey()][2])
	require.Equal(t, keyPairTwo.PrivateKey(), allKeys[keyPairTwo.PublicKey()][4])

}

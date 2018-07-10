package db

import (
	"testing"

	require "github.com/stretchr/testify/require"
	dr "github.com/tiabc/doubleratchet"
)

func TestStore_PutGet(t *testing.T) {

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

	key, exist := s.Get(keyPair.PublicKey(), 3)
	require.True(t, exist)

	require.Equal(t, keyPair.PrivateKey(), key)

}

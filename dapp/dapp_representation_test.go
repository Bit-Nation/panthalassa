package dapp

import (
	"bytes"
	"crypto/rand"
	"golang.org/x/crypto/ed25519"
	"testing"

	mh "github.com/multiformats/go-multihash"
	require "github.com/stretchr/testify/require"
)

func TestDAppRepresentationHash(t *testing.T) {

	pub, _, err := ed25519.GenerateKey(rand.Reader)
	require.Nil(t, err)

	rep := JsonRepresentation{
		Name:               "Send / Receive Money",
		Code:               []byte(`var wallet = "0x930aa9a843266bdb02847168d571e7913907dd84"`),
		SignaturePublicKey: pub,
	}

	// calculate hash manually
	// name + code + signature public key
	buff := bytes.NewBuffer([]byte(rep.Name))

	_, err = buff.Write([]byte(rep.Code))
	require.Nil(t, err)

	_, err = buff.Write([]byte(rep.SignaturePublicKey))
	require.Nil(t, err)

	expectedHash, err := mh.Sum(buff.Bytes(), mh.SHA3_256, -1)
	require.Nil(t, err)

	// calculate hash
	calculateHash, err := rep.hash()
	require.Nil(t, err)

	// check if hashes match
	require.Equal(t, string(expectedHash), string(calculateHash))

}

func TestDAppVerifySignature(t *testing.T) {

	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	require.Nil(t, err)

	rep := JsonRepresentation{
		Name:               "Send / Receive Money",
		Code:               []byte(`var wallet = "0x930aa9a843266bdb02847168d571e7913907dd84"`),
		SignaturePublicKey: pub,
	}

	// validate signature
	// should be invalid since it doesn't exist
	valid, err := rep.VerifySignature()
	require.Nil(t, err)
	require.False(t, valid)

	// hash the representation
	calculatedHash, err := rep.hash()
	require.Nil(t, err)

	// sign representation
	rep.Signature = ed25519.Sign(priv, calculatedHash)

	// validate signature
	valid, err = rep.VerifySignature()
	require.Nil(t, err)
	require.True(t, valid)

}

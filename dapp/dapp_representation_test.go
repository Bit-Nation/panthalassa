package dapp

import (
	"crypto/rand"
	"crypto/sha256"
	"golang.org/x/crypto/ed25519"
	"testing"

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
	h := sha256.New()
	_, err = h.Write([]byte(rep.Name))
	require.Nil(t, err)
	_, err = h.Write(rep.Code)
	require.Nil(t, err)
	_, err = h.Write([]byte(rep.SignaturePublicKey))
	require.Nil(t, err)
	expectedHash := h.Sum(nil)

	// calculate hash
	calculateHash, err := rep.hash()
	require.Nil(t, err)

	// check if hashes match
	require.Equal(t, expectedHash, calculateHash)

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

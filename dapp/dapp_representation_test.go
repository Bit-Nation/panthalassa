package dapp

import (
	"crypto/rand"
	"encoding/hex"
	"testing"

	mh "github.com/multiformats/go-multihash"
	require "github.com/stretchr/testify/require"
	ed25519 "golang.org/x/crypto/ed25519"
)

func TestDAppRepresentationHash(t *testing.T) {

	// we need a test hash for the JsonToJsonBuild method
	testHash, err := mh.Sum(nil, mh.SHA2_256, -1)
	require.Nil(t, err)

	dAppJson := RawData{
		Name: map[string]string{
			"en-us": "send and request money",
			"de":    "sende und fordere geld an",
		},
		UsedSigningKey: "110c3ff292fb8ebf0084a9fc1e8c06418ab1c2cbd1058d87e78aa0fcdcbf5791",
		Code:           `var wallet = "0x930aa9a843266bdb02847168d571e7913907dd84"`,
		Image:          "aGk=",
		Engine:         "1.2.3",
		Signature:      testHash.String(),
		Version:        1,
	}

	// calculate hash
	jsonBuild, err := ParseJsonToData(dAppJson)
	require.Nil(t, err)
	calculatedHash, err := jsonBuild.Hash()
	require.Nil(t, err)

	// check if hashes match
	require.Equal(t, "122045b810a58b64b3d35e46a30c8b1d80dccb8f04142acb4f6fe86e01792316a3e4", hex.EncodeToString(calculatedHash))

}

func TestEngineVersionToSV(t *testing.T) {

	// should pass
	sv, err := engineVersionToSV("3.2.1")
	require.Nil(t, err)
	require.Equal(t, SV{
		Major: 3,
		Minor: 2,
		Patch: 1,
	}, sv)

	// should fail since we provided an invalid version
	_, err = engineVersionToSV("3.2")
	require.EqualError(t, err, "invalid version - a version must consist of 3 numbers separated by dots")

}

func TestDAppVerifySignature(t *testing.T) {

	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	require.Nil(t, err)

	// create fake signature
	fakeSignature, err := mh.Sum([]byte("invalid signature"), mh.SHA3_256, -1)
	require.Nil(t, err)

	dAppJson := RawData{
		Name: map[string]string{
			"en-us": "send and request money",
			"de":    "sende und fordere geld an",
		},
		UsedSigningKey: hex.EncodeToString(pub),
		Code:           `var wallet = "0x930aa9a843266bdb02847168d571e7913907dd84"`,
		Image:          "aGk=",
		Engine:         "1.2.3",
		Signature:      fakeSignature.String(),
	}

	// validate signature
	// should be invalid since it doesn't exist
	jsonDApp, err := ParseJsonToData(dAppJson)
	require.Nil(t, err)
	valid, err := jsonDApp.VerifySignature()
	require.Nil(t, err)
	require.False(t, valid)

	// add the signature to the json DApp
	jsonDAppHash, err := jsonDApp.Hash()
	require.Nil(t, err)
	jsonDApp.Signature = ed25519.Sign(priv, jsonDAppHash)

	// signature invalid
	valid, err = jsonDApp.VerifySignature()
	require.Nil(t, err)
	require.True(t, valid)

}

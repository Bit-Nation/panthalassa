package dapp

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"testing"

	mh "github.com/multiformats/go-multihash"
	require "github.com/stretchr/testify/require"
	ed25519 "golang.org/x/crypto/ed25519"
)

func TestDAppRepresentationHash(t *testing.T) {

	pub, _, err := ed25519.GenerateKey(rand.Reader)
	require.Nil(t, err)

	// we need a test hash for the JsonToJsonBuild method
	testHash, err := mh.Sum(nil, mh.SHA3_256, -1)
	require.Nil(t, err)

	dAppJson := jsonBuild{
		Name: map[string]string{
			"en-us": "send and request money",
			"de":    "sende und fordere geld an",
		},
		UsedSigningKey: hex.EncodeToString(pub),
		Code:           `var wallet = "0x930aa9a843266bdb02847168d571e7913907dd84"`,
		Image:          "base64...",
		Engine:         "1.2.3",
		Signature:      testHash.String(),
	}

	// calculate hash manually
	// name + code + signature public key
	buff := bytes.NewBuffer(nil)

	// write de first since we sort for the language codes
	_, err = buff.WriteString("de")
	require.Nil(t, err)
	_, err = buff.WriteString("sende und fordere geld an")
	require.Nil(t, err)

	// write en us
	_, err = buff.WriteString("en-us")
	require.Nil(t, err)
	_, err = buff.WriteString("send and request money")
	require.Nil(t, err)

	// write used signing key
	_, err = buff.Write(pub)
	require.Nil(t, err)

	// write code to buffer
	_, err = buff.WriteString(`var wallet = "0x930aa9a843266bdb02847168d571e7913907dd84"`)
	require.Nil(t, err)

	// write image
	_, err = buff.WriteString("base64...")
	require.Nil(t, err)

	// write version
	_, err = buff.WriteString("1.2.3")
	require.Nil(t, err)

	expectedHash, err := mh.Sum(buff.Bytes(), mh.SHA3_256, -1)
	require.Nil(t, err)

	// calculate hash
	jsonBuild, err := JsonToJsonBuild(dAppJson)
	require.Nil(t, err)
	calculatedHash, err := jsonBuild.Hash()
	require.Nil(t, err)

	// check if hashes match
	require.Equal(t, hex.EncodeToString(expectedHash), hex.EncodeToString(calculatedHash))

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

	dAppJson := jsonBuild{
		Name: map[string]string{
			"en-us": "send and request money",
			"de":    "sende und fordere geld an",
		},
		UsedSigningKey: hex.EncodeToString(pub),
		Code:           `var wallet = "0x930aa9a843266bdb02847168d571e7913907dd84"`,
		Image:          "base64...",
		Engine:         "1.2.3",
		Signature:      fakeSignature.String(),
	}

	// validate signature
	// should be invalid since it doesn't exist
	jsonDApp, err := JsonToJsonBuild(dAppJson)
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

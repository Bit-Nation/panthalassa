package dapp

import (
	"crypto/rand"
	"testing"
	"time"

	dAppMod "github.com/Bit-Nation/panthalassa/dapp/module"
	log "github.com/op/go-logging"
	require "github.com/stretchr/testify/require"
	ed25519 "golang.org/x/crypto/ed25519"
)

func TestStartDAppSuccess(t *testing.T) {

	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	require.Nil(t, err)

	app := JsonBuild{
		Name: map[string]string{
			"en-us": "send and request money",
		},
		UsedSigningKey: pub,
		Code:           []byte("var i = 1"),
		Image:          []byte("base64..."),
		Engine: SV{
			Major: 0,
			Minor: 1,
			Patch: 0,
		},
	}

	appHash, err := app.Hash()
	require.Nil(t, err)

	app.Signature = ed25519.Sign(priv, appHash)

	// mock the verify signature function since
	// we don't care about correctness

	closer := make(chan *JsonBuild)

	_, err = New(log.MustGetLogger(""), &app, []dAppMod.Module{}, closer, time.Second)
	require.Nil(t, err)

}

func TestStartDAppHalting(t *testing.T) {

	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	require.Nil(t, err)

	app := JsonBuild{
		Name: map[string]string{
			"en-us": "send and request money",
		},
		UsedSigningKey: pub,
		Code:           []byte("while(true){}"),
		Image:          []byte("base64..."),
		Engine: SV{
			Major: 0,
			Minor: 1,
			Patch: 0,
		},
	}

	appHash, err := app.Hash()
	require.Nil(t, err)

	app.Signature = ed25519.Sign(priv, appHash)

	// mock the verify signature function since
	// we don't care about correctness

	closer := make(chan *JsonBuild, 1)

	dApp, err := New(log.MustGetLogger(""), &app, []dAppMod.Module{}, closer, time.Second)
	require.Nil(t, dApp)
	require.EqualError(t, err, "timeout - failed to start DApp")

}

func TestStartInvalidSignature(t *testing.T) {

	pub, _, err := ed25519.GenerateKey(rand.Reader)
	require.Nil(t, err)

	_, invalidPriv, err := ed25519.GenerateKey(rand.Reader)
	require.Nil(t, err)

	app := JsonBuild{
		Name: map[string]string{
			"en-us": "send and request money",
		},
		UsedSigningKey: pub,
		Code:           []byte(""),
		Image:          []byte("base64..."),
		Engine: SV{
			Major: 0,
			Minor: 1,
			Patch: 0,
		},
	}

	appHash, err := app.Hash()
	require.Nil(t, err)

	// sign with invalid private key
	app.Signature = ed25519.Sign(invalidPriv, appHash)

	// mock the verify signature function since
	// we don't care about correctness

	closer := make(chan *JsonBuild, 1)

	dApp, err := New(log.MustGetLogger(""), &app, []dAppMod.Module{}, closer, time.Second)
	require.Nil(t, dApp)
	require.EqualError(t, err, "failed to verify signature for DApp")

}

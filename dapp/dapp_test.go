package dapp

import (
	"testing"
	"time"

	"crypto/rand"
	dappMod "github.com/Bit-Nation/panthalassa/dapp/module"
	log "github.com/op/go-logging"
	require "github.com/stretchr/testify/require"
	ed25519 "golang.org/x/crypto/ed25519"
)

const validSampleDApp = ``

type tEstasdf struct {
}

func TestStartDAppSuccess(t *testing.T) {

	app := JsonRepresentation{
		Name: "My DApp",
		Code: "var i = 0;",
	}

	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	require.Nil(t, err)

	app.SignaturePublicKey = pub

	appHash, err := app.Hash()
	require.Nil(t, err)

	app.Signature = ed25519.Sign(priv, appHash)

	// mock the verify signature function since
	// we don't care about correctness

	closer := make(chan *JsonRepresentation)

	_, err = New(log.MustGetLogger(""), &app, []dappMod.Module{}, closer, time.Second)
	require.Nil(t, err)

}

func TestStartDAppHalting(t *testing.T) {

	app := JsonRepresentation{
		Name: "My DApp",
		Code: "while (true) {}",
	}

	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	require.Nil(t, err)

	app.SignaturePublicKey = pub

	appHash, err := app.Hash()
	require.Nil(t, err)

	app.Signature = ed25519.Sign(priv, appHash)

	// mock the verify signature function since
	// we don't care about correctness

	closer := make(chan *JsonRepresentation, 1)

	dapp, err := New(log.MustGetLogger(""), &app, []dappMod.Module{}, closer, time.Second)
	require.Nil(t, dapp)
	require.EqualError(t, err, "timeout - failed to start DApp")

}

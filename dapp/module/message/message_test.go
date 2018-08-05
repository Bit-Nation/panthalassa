package message

import (
	"crypto/rand"
	"encoding/hex"
	"testing"
	"time"

	db "github.com/Bit-Nation/panthalassa/db"
	otto "github.com/robertkrimen/otto"
	require "github.com/stretchr/testify/require"
	ed25519 "golang.org/x/crypto/ed25519"
)

func TestPersistMessageSuccessfully(t *testing.T) {

	vm := otto.New()

	// dapp pub key
	dAppPubKey, _, err := ed25519.GenerateKey(rand.Reader)
	require.Nil(t, err)

	// chat pub key
	chatPubKey, _, err := ed25519.GenerateKey(rand.Reader)
	require.Nil(t, err)

	// storage mock
	calledPersist := false
	msgStorage := testMessageStorage{
		persistDAppMessage: func(partner ed25519.PublicKey, msg db.DAppMessage) error {

			require.Equal(t, "SEND_MONEY", msg.Type)
			require.Equal(t, "value", msg.Params["key"])
			require.True(t, msg.ShouldSend)

			calledPersist = true
			return nil
		},
	}

	msgModule := New(&msgStorage, dAppPubKey, nil)
	require.Nil(t, msgModule.Register(vm))

	msg := map[string]interface{}{
		"shouldSend": true,
		"params": map[string]interface{}{
			"key": "value",
		},
		"type": "SEND_MONEY",
	}

	closer := make(chan struct{}, 1)

	vm.Call("sendMessage", vm, hex.EncodeToString(chatPubKey), msg, func(call otto.FunctionCall) otto.Value {

		err := call.Argument(0)
		if err.IsDefined() {
			require.Fail(t, err.String())
		}

		// make sure the message got persisted
		require.True(t, calledPersist)

		closer <- struct{}{}

		return otto.Value{}
	})

	select {
	case <-closer:
	case <-time.After(time.Second):
		require.Fail(t, "timed out")
	}

}

func TestPersistInvalidFunctionCall(t *testing.T) {

	vm := otto.New()

	// dapp pub key
	dAppPubKey, _, err := ed25519.GenerateKey(rand.Reader)
	require.Nil(t, err)

	// chat pub key
	chatPubKey, _, err := ed25519.GenerateKey(rand.Reader)
	require.Nil(t, err)

	msgModule := New(nil, dAppPubKey, nil)
	require.Nil(t, msgModule.Register(vm))

	closer := make(chan struct{}, 1)

	vm.Call("sendMessage", vm, hex.EncodeToString(chatPubKey), nil, func(call otto.FunctionCall) otto.Value {

		err := call.Argument(0)
		require.Equal(t, "ValidationError: expected parameter 1 to be of type object", err.String())

		closer <- struct{}{}

		return otto.Value{}
	})

	select {
	case <-closer:
	case <-time.After(time.Second):
		require.Fail(t, "timed out")
	}

}

// test whats happening if the chat is too short
func TestPersistChatTooShort(t *testing.T) {

	vm := otto.New()

	msgModule := New(nil, nil, nil)
	require.Nil(t, msgModule.Register(vm))

	closer := make(chan struct{}, 1)

	vm.Call("sendMessage", vm, hex.EncodeToString(make([]byte, 16)), map[string]interface{}{"shouldSend": true}, func(call otto.FunctionCall) otto.Value {

		err := call.Argument(0)
		require.Equal(t, "chat must be 32 bytes long", err.String())

		closer <- struct{}{}

		return otto.Value{}
	})

	select {
	case <-closer:
	case <-time.After(time.Second):
		require.Fail(t, "timed out")
	}

}

// test what happens if we try to send without
func TestPersistWithoutShouldSend(t *testing.T) {

	vm := otto.New()

	msgModule := New(nil, nil, nil)
	require.Nil(t, msgModule.Register(vm))

	closer := make(chan struct{}, 1)

	vm.Call("sendMessage", vm, hex.EncodeToString(make([]byte, 16)), map[string]interface{}{}, func(call otto.FunctionCall) otto.Value {

		err := call.Argument(0)
		require.Equal(t, "ValidationError: Missing value for required key: shouldSend", err.String())

		closer <- struct{}{}

		return otto.Value{}
	})

	select {
	case <-closer:
	case <-time.After(time.Second):
		require.Fail(t, "timed out")
	}

}

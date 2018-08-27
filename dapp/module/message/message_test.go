package message

import (
	"crypto/rand"
	"encoding/hex"
	"testing"
	"time"

	db "github.com/Bit-Nation/panthalassa/db"
	require "github.com/stretchr/testify/require"
	ed25519 "golang.org/x/crypto/ed25519"
	duktape "gopkg.in/olebedev/go-duktape.v3"
)

func TestPersistMessageSuccessfully(t *testing.T) {

	vm := duktape.New()

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

	msg := `({
		"shouldSend": true,
		"params": {
			"key": "value",
		},
		"type": "SEND_MONEY",
	})`

	closer := make(chan struct{}, 1)

	_, err = vm.PushGlobalGoFunction("callbackSendMessage", func(context *duktape.Context) int {
		if !context.IsUndefined(0) {
			require.Fail(t, context.ToString(0))
		}

		// make sure the message got persisted
		require.True(t, calledPersist)

		closer <- struct{}{}

		return 0
	})
	require.Nil(t, err)
	vm.PevalString(`sendMessage(` + `"` + hex.EncodeToString(chatPubKey) + `"` + `,` + msg + `,` + `callbackSendMessage)`)
	select {
	case <-closer:
	case <-time.After(time.Second):
		require.Fail(t, "timed out")
	}

}

func TestPersistInvalidFunctionCall(t *testing.T) {

	vm := duktape.New()

	// dapp pub key
	dAppPubKey, _, err := ed25519.GenerateKey(rand.Reader)
	require.Nil(t, err)

	// chat pub key
	chatPubKey, _, err := ed25519.GenerateKey(rand.Reader)
	require.Nil(t, err)

	msgModule := New(nil, dAppPubKey, nil)
	require.Nil(t, msgModule.Register(vm))

	closer := make(chan struct{}, 1)

	_, err = vm.PushGlobalGoFunction("callbackSendMessage", func(context *duktape.Context) int {

		err := context.ToString(0)
		require.Equal(t, "ValidationError: expected parameter 1 to be of type object", err)
		closer <- struct{}{}
		return 0
	})
	require.Nil(t, err)
	vm.PevalString(`sendMessage(` + `"` + hex.EncodeToString(chatPubKey) + `"` + `,` + `"nil"` + `,` + `callbackSendMessage)`)
	select {
	case <-closer:
	case <-time.After(time.Second):
		require.Fail(t, "timed out")
	}

}

// test whats happening if the chat is too short
func TestPersistChatTooShort(t *testing.T) {

	vm := duktape.New()

	msgModule := New(nil, nil, nil)
	require.Nil(t, msgModule.Register(vm))

	closer := make(chan struct{}, 1)

	_, err := vm.PushGlobalGoFunction("callbackSendMessage", func(context *duktape.Context) int {
		err := context.ToString(0)
		require.Equal(t, "chat must be 32 bytes long", err)

		closer <- struct{}{}

		return 0

	})
	require.Nil(t, err)
	vm.PevalString(`sendMessage(` + `"` + hex.EncodeToString(make([]byte, 16)) + `"` + `,` + `({"shouldSend": true})` + `,` + `callbackSendMessage)`)
	select {
	case <-closer:
	case <-time.After(time.Second):
		require.Fail(t, "timed out")
	}

}

// test what happens if we try to send without
func TestPersistWithoutShouldSend(t *testing.T) {

	vm := duktape.New()

	msgModule := New(nil, nil, nil)
	require.Nil(t, msgModule.Register(vm))

	closer := make(chan struct{}, 1)

	_, err := vm.PushGlobalGoFunction("callbackSendMessage", func(context *duktape.Context) int {
		err := context.ToString(0)
		require.Equal(t, "ValidationError: Missing value for required key: shouldSend", err)
		closer <- struct{}{}
		return 0
	})
	require.Nil(t, err)
	vm.PevalString(`sendMessage(` + `"` + hex.EncodeToString(make([]byte, 16)) + `"` + `,` + `({})` + `,` + `callbackSendMessage)`)
	select {
	case <-closer:
	case <-time.After(time.Second):
		require.Fail(t, "timed out")
	}

}

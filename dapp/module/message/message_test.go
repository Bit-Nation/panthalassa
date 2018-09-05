package message

import (
	"crypto/rand"
	"encoding/hex"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	db "github.com/Bit-Nation/panthalassa/db"
	km "github.com/Bit-Nation/panthalassa/keyManager"
	ks "github.com/Bit-Nation/panthalassa/keyStore"
	mnemonic "github.com/Bit-Nation/panthalassa/mnemonic"
	storm "github.com/asdine/storm"
	require "github.com/stretchr/testify/require"
	ed25519 "golang.org/x/crypto/ed25519"
	duktape "gopkg.in/olebedev/go-duktape.v3"
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

func createKeyManager() *km.KeyManager {

	mne, err := mnemonic.New()
	if err != nil {
		panic(err)
	}

	keyStore, err := ks.NewFromMnemonic(mne)
	if err != nil {
		panic(err)
	}

	return km.CreateFromKeyStore(keyStore)

}

func TestPersistMessageSuccessfully(t *testing.T) {

	vm := duktape.New()

	// dapp pub key
	dAppPubKey, _, err := ed25519.GenerateKey(rand.Reader)
	require.Nil(t, err)

	// chat pub key
	chatPubKey, _, err := ed25519.GenerateKey(rand.Reader)
	require.Nil(t, err)

	// closer
	closer := make(chan struct{}, 1)

	// storage mock
	chatStorage := db.NewChatStorage(createStorm(), []func(e db.MessagePersistedEvent){}, createKeyManager())
	chatStorage.AddListener(func(e db.MessagePersistedEvent) {
		chat, _ := chatStorage.GetChat(chatPubKey)
		if err != nil {
			panic(err)
		}

		messages, err := chat.Messages(0, 1)
		if err != nil {
			panic(err)
		}

		m := messages[0]

		if m.DApp.Type != "SEND_MONEY" {
			panic("invalid type")
		}

		if !m.DApp.ShouldSend {
			panic("should send is supposed to be true")
		}

		if m.DApp.Params["key"] != "value" {
			panic("invalid value for key")
		}

		if hex.EncodeToString(m.DApp.DAppPublicKey) != hex.EncodeToString(dAppPubKey) {
			panic("invalid dapp pub key")
		}

		closer <- struct{}{}
	})
	// create chat
	require.Nil(t, chatStorage.CreateChat(chatPubKey))

	msgModule := New(chatStorage, dAppPubKey, nil)
	require.Nil(t, msgModule.Register(vm))

	msg := `({
		"shouldSend": true,
		"params": {
			"key": "value",
		},
		"type": "SEND_MONEY",
	})`

	_, err = vm.PushGlobalGoFunction("callbackSendMessage", func(context *duktape.Context) int {
		if !context.IsUndefined(0) {
			require.Fail(t, context.ToString(0))
		}

		return 0
	})
	require.Nil(t, err)
	vm.PevalString(`sendMessage(` + `"` + hex.EncodeToString(chatPubKey) + `"` + `,` + msg + `,` + `callbackSendMessage)`)
	select {
	case <-closer:
	case <-time.After(time.Second * 5):
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

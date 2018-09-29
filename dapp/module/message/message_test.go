package message

import (
	"crypto/rand"
	"encoding/hex"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	ed25519 "golang.org/x/crypto/ed25519"

	db "github.com/Bit-Nation/panthalassa/db"
	km "github.com/Bit-Nation/panthalassa/keyManager"
	ks "github.com/Bit-Nation/panthalassa/keyStore"
	mnemonic "github.com/Bit-Nation/panthalassa/mnemonic"
	storm "github.com/asdine/storm"
	otto "github.com/robertkrimen/otto"
	require "github.com/stretchr/testify/require"
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

	vm := otto.New()

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
		chat, _ := chatStorage.GetChatByPartner(chatPubKey)
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
	_, err = chatStorage.CreateChat(chatPubKey)
	require.Nil(t, err)

	msgModule := New(chatStorage, dAppPubKey, nil)
	require.Nil(t, msgModule.Register(vm))

	msg := map[string]interface{}{
		"shouldSend": true,
		"params": map[string]interface{}{
			"key": "value",
		},
		"type": "SEND_MONEY",
	}

	vm.Call("sendMessage", vm, hex.EncodeToString(chatPubKey), msg, func(call otto.FunctionCall) otto.Value {

		err := call.Argument(0)
		if err.IsDefined() {
			require.Fail(t, err.String())
		}

		closer <- struct{}{}

		return otto.Value{}
	})

	select {
	case <-closer:
	case <-time.After(time.Second):
		require.Fail(t, "timed out")
	}

}

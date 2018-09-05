package chat

import (
	"crypto/rand"
	"testing"

	db "github.com/Bit-Nation/panthalassa/db"
	"github.com/stretchr/testify/require"
	ed25519 "golang.org/x/crypto/ed25519"
)

func TestChat_SavePrivateMessage(t *testing.T) {

	pub, _, err := ed25519.GenerateKey(rand.Reader)
	require.Nil(t, err)

	km := createKeyManager()

	c := Chat{
		chatStorage: db.NewChatStorage(createStorm(), []func(e db.MessagePersistedEvent){}, km),
		km:          km,
	}

	require.Nil(t, c.SavePrivateMessage(pub, []byte("hi")))

}

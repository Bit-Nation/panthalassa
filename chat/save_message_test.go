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

	chatStorage := db.NewChatStorage(createStorm(), []func(e db.MessagePersistedEvent){}, km)

	c := Chat{
		chatStorage: chatStorage,
		km:          km,
	}

	require.Nil(t, chatStorage.CreateChat(pub))

	chat, err := chatStorage.GetChatByPartner(pub)
	require.Nil(t, err)

	require.Nil(t, c.SaveMessage(chat.ID, []byte("hi")))

}

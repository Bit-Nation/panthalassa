package chat

import (
	"time"

	db "github.com/Bit-Nation/panthalassa/db"
	"github.com/kataras/iris/core/errors"
	ed25519 "golang.org/x/crypto/ed25519"
)

var nowAsUnix = func() int64 {
	return time.Now().UnixNano()
}

// persist private message
func (c *Chat) SavePrivateMessage(to ed25519.PublicKey, rawMessage []byte) error {
	msg := db.Message{
		Message:   rawMessage,
		CreatedAt: nowAsUnix(),
	}
	// fetch chat
	chat, err := c.chatStorage.GetChat(to)
	if err != nil {
		return err
	}
	if chat == nil {
		// create chat if not exist
		if err := c.chatStorage.CreateChat(to); err != nil {
			return err
		}
	}
	// fetch chat again
	chat, err = c.chatStorage.GetChat(to)
	if err != nil {
		return err
	}
	if chat == nil {
		return errors.New("got invalid chat")
	}

	return chat.PersistMessage(msg)
}

package chat

import (
	"errors"
	"time"
	"encoding/hex"
	
	db "github.com/Bit-Nation/panthalassa/db"
	uuid "github.com/satori/go.uuid"
	ed25519 "golang.org/x/crypto/ed25519"
)

var nowAsUnix = func() int64 {
	return time.Now().UnixNano()
}

// persist private message
func (c *Chat) SavePrivateMessage(to ed25519.PublicKey, rawMessage []byte) error {
	id, err := uuid.NewV4()
	if err != nil {
		return err
	}

	senderStr, err := c.km.IdentityPublicKey()
	if err != nil {
		return err
	}
	sender, err := hex.DecodeString(senderStr)
	if err != nil {
		return err
	}

	msg := db.Message{
		ID:        id.String(),
		Message:   rawMessage,
		CreatedAt: nowAsUnix(),
		Version:   1,
		Status:    db.StatusPersisted,
		Sender:    sender,
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

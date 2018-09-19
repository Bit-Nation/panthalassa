package chat

import (
	"encoding/hex"
	"fmt"
	"time"

	db "github.com/Bit-Nation/panthalassa/db"
	uuid "github.com/satori/go.uuid"
)

var nowAsUnix = func() int64 {
	return time.Now().UnixNano()
}

func (c *Chat) SaveMessage(chatID int, rawMessage []byte) error {
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
	chat, err := c.chatStorage.GetChat(chatID)
	if err != nil {
		return err
	}
	if chat == nil {
		return fmt.Errorf("got invalid chat for id: %d", chatID)
	}

	return chat.PersistMessage(msg)
}

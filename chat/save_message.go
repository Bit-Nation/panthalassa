package chat

import (
	"time"

	ed25519 "golang.org/x/crypto/ed25519"
	db "github.com/Bit-Nation/panthalassa/db"
)

var nowAsUnix = func() int64 {
	return time.Now().UnixNano()
}

// persist private message
func (c *Chat) SavePrivateMessage(to ed25519.PublicKey, rawMessage []byte) error {
	msg := db.Message{
		Message: rawMessage,
		CreatedAt: nowAsUnix(),
	}
	return c.messageDB.PersistMessageToSend(to, msg)
}

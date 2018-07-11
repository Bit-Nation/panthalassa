package chat

import (
	"time"

	bpb "github.com/Bit-Nation/protobuffers"
	uuid "github.com/satori/go.uuid"
	ed25519 "golang.org/x/crypto/ed25519"
)

var nowAsUnix = func() int64 {
	return time.Now().Unix()
}

// persist private message
func (c *Chat) SavePrivateMessage(to ed25519.PublicKey, msg bpb.PlainChatMessage) error {
	msg.CreatedAt = nowAsUnix()
	id, err := uuid.NewV4()
	if err != nil {
		return err
	}
	msg.MessageID = id.String()
	return c.messageDB.PersistMessage(to, msg)
}

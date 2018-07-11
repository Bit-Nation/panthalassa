package chat

import (
	bpb "github.com/Bit-Nation/protobuffers"
	ed25519 "golang.org/x/crypto/ed25519"
)

type MessageDatabase interface {
	PersistMessage(to ed25519.PublicKey, msg bpb.PlainChatMessage) error
}

type Chat struct {
	messageDB MessageDatabase
}

func NewChat(msgDB MessageDatabase) *Chat {
	return &Chat{
		messageDB: msgDB,
	}
}

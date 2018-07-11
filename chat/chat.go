package chat

import (
	db "github.com/Bit-Nation/panthalassa/db"
	keyManager "github.com/Bit-Nation/panthalassa/keyManager"
	bpb "github.com/Bit-Nation/protobuffers"
	x3dh "github.com/Bit-Nation/x3dh"
	dr "github.com/tiabc/doubleratchet"
	ed25519 "golang.org/x/crypto/ed25519"
)

// represents the source were we get things
// like the pre key bundles from
type Backend interface {
	FetchPreKeyBundle(userIDPubKey ed25519.PublicKey) (x3dh.PreKeyBundle, error)
	SubmitMessage(msg bpb.ChatMessage) error
}

type Chat struct {
	messageDB        db.ChatMessageStorage
	backend          Backend
	sharedSecStorage db.SharedSecretStorage
	x3dh             x3dh.X3dh
	km               *keyManager.KeyManager
	drKeyStorage     dr.KeysStorage
}

func NewChat(msgDB db.ChatMessageStorage) *Chat {
	return &Chat{
		messageDB: msgDB,
	}
}

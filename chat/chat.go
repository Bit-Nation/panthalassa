package chat

import (
	"time"

	prekey "github.com/Bit-Nation/panthalassa/chat/prekey"
	db "github.com/Bit-Nation/panthalassa/db"
	keyManager "github.com/Bit-Nation/panthalassa/keyManager"
	bpb "github.com/Bit-Nation/protobuffers"
	x3dh "github.com/Bit-Nation/x3dh"
	dr "github.com/tiabc/doubleratchet"
	ed25519 "golang.org/x/crypto/ed25519"
)

const (
	SignedPreKeyValidTimeFrame = time.Hour * 24 * 60
)

// represents the source were we get things
// like the pre key bundles from
type Backend interface {
	FetchPreKeyBundle(userIDPubKey ed25519.PublicKey) (x3dh.PreKeyBundle, error)
	SubmitMessage(msg bpb.ChatMessage) error
	FetchSignedPreKey(userIdPubKey ed25519.PublicKey) (prekey.PreKey, error)
}

type Chat struct {
	messageDB           db.ChatMessageStorage
	backend             Backend
	sharedSecStorage    db.SharedSecretStorage
	x3dh                *x3dh.X3dh
	km                  *keyManager.KeyManager
	drKeyStorage        dr.KeysStorage
	signedPreKeyStorage db.SignedPreKeyStorage
	userStorage         db.UserStorage
}

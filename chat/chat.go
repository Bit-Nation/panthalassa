package chat

import (
	"time"

	preKey "github.com/Bit-Nation/panthalassa/chat/prekey"
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

type Backend interface {
	FetchPreKeyBundle(userIDPubKey ed25519.PublicKey) (x3dh.PreKeyBundle, error)
	SubmitMessage(msg bpb.ChatMessage) error
	FetchSignedPreKey(userIdPubKey ed25519.PublicKey) (preKey.PreKey, error)
}

type Chat struct {
	messageDB            db.ChatMessageStorage
	backend              Backend
	sharedSecStorage     db.SharedSecretStorage
	x3dh                 *x3dh.X3dh
	km                   *keyManager.KeyManager
	drKeyStorage         dr.KeysStorage
	signedPreKeyStorage  db.SignedPreKeyStorage
	oneTimePreKeyStorage db.OneTimePreKeyStorage
	userStorage          db.UserStorage
}

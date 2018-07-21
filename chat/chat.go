package chat

import (
	"time"

	backend "github.com/Bit-Nation/panthalassa/backend"
	preKey "github.com/Bit-Nation/panthalassa/chat/prekey"
	db "github.com/Bit-Nation/panthalassa/db"
	keyManager "github.com/Bit-Nation/panthalassa/keyManager"
	bpb "github.com/Bit-Nation/protobuffers"
	x3dh "github.com/Bit-Nation/x3dh"
	log "github.com/ipfs/go-log"
	dr "github.com/tiabc/doubleratchet"
	ed25519 "golang.org/x/crypto/ed25519"
)

const (
	SignedPreKeyValidTimeFrame = time.Hour * 24 * 60
)

var logger = log.Logger("chat")

type Backend interface {
	FetchPreKeyBundle(userIDPubKey ed25519.PublicKey) (x3dh.PreKeyBundle, error)
	SubmitMessages(messages []*bpb.ChatMessage) error
	FetchSignedPreKey(userIdPubKey ed25519.PublicKey) (preKey.PreKey, error)
	AddRequestHandler(handler backend.RequestHandler)
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

func NewChat(b Backend) (*Chat, error) {

	c := &Chat{}

	// register messages handler
	msgHandler := c.messagesHandler
	b.AddRequestHandler(&msgHandler)

	return c, nil
}

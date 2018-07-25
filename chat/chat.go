package chat

import (
	"crypto/rand"
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

func (c *Chat) AllChats() ([]ed25519.PublicKey, error) {
	return c.messageDB.AllChats()
}

func (c *Chat) Messages(partner ed25519.PublicKey, start int64, amount uint) (map[int64]db.Message, error) {
	return c.messageDB.Messages(partner, start, amount)
}

type Config struct {
	MessageDB            db.ChatMessageStorage
	Backend              Backend
	SharedSecretDB       db.SharedSecretStorage
	KM                   *keyManager.KeyManager
	DRKeyStorage         dr.KeysStorage
	SignedPreKeyStorage  db.SignedPreKeyStorage
	OneTimePreKeyStorage db.OneTimePreKeyStorage
	UserStorage          db.UserStorage
}

// @todo we need to construct the x3dh instance correct
func NewChat(conf Config) (*Chat, error) {

	x3dh.NewCurve25519(rand.Reader)

	c := &Chat{
		messageDB:            conf.MessageDB,
		backend:              conf.Backend,
		sharedSecStorage:     conf.SharedSecretDB,
		x3dh:                 nil,
		km:                   conf.KM,
		drKeyStorage:         conf.DRKeyStorage,
		signedPreKeyStorage:  conf.SignedPreKeyStorage,
		oneTimePreKeyStorage: conf.OneTimePreKeyStorage,
		userStorage:          conf.UserStorage,
	}

	// register messages handler
	c.backend.AddRequestHandler(c.messagesHandler)

	// register now one time pre key handler
	c.backend.AddRequestHandler(c.oneTimePreKeysHandler)

	return c, nil
}

package chat

import (
	"crypto/rand"
	"crypto/sha256"
	"time"

	backend "github.com/Bit-Nation/panthalassa/backend"
	preKey "github.com/Bit-Nation/panthalassa/chat/prekey"
	db "github.com/Bit-Nation/panthalassa/db"
	keyManager "github.com/Bit-Nation/panthalassa/keyManager"
	queue "github.com/Bit-Nation/panthalassa/queue"
	uiapi "github.com/Bit-Nation/panthalassa/uiapi"
	bpb "github.com/Bit-Nation/protobuffers"
	x3dh "github.com/Bit-Nation/x3dh"
	log "github.com/ipfs/go-log"
	dr "github.com/tiabc/doubleratchet"
	ed25519 "golang.org/x/crypto/ed25519"
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
	uiApi                *uiapi.Api
	queue                *queue.Queue
}

func (c *Chat) AllChats() ([]ed25519.PublicKey, error) {
	return c.messageDB.AllChats()
}

func (c *Chat) Messages(partner ed25519.PublicKey, start int64, amount uint) ([]db.Message, error) {
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
	UiApi                *uiapi.Api
	Queue                *queue.Queue
}

func NewChat(conf Config) (*Chat, error) {

	// my chat id key pair
	myChatIDKeyPair, err := conf.KM.ChatIdKeyPair()
	if err != nil {
		return nil, err
	}

	// curve 25519
	c25519 := x3dh.NewCurve25519(rand.Reader)
	myX3dh := x3dh.New(&c25519, sha256.New(), "pangea-chat", myChatIDKeyPair)

	c := &Chat{
		messageDB:            conf.MessageDB,
		backend:              conf.Backend,
		sharedSecStorage:     conf.SharedSecretDB,
		x3dh:                 &myX3dh,
		km:                   conf.KM,
		drKeyStorage:         conf.DRKeyStorage,
		signedPreKeyStorage:  conf.SignedPreKeyStorage,
		oneTimePreKeyStorage: conf.OneTimePreKeyStorage,
		userStorage:          conf.UserStorage,
		uiApi:                conf.UiApi,
		queue:                conf.Queue,
	}

	err = c.queue.RegisterProcessor(&SubmitMessagesProcessor{
		chat:  c,
		msgDB: c.messageDB,
		queue: c.queue,
	})
	if err != nil {
		return nil, err
	}

	// add message handler that will inform the ui about updates
	c.messageDB.AddListener(c.handlePersistedMessage)

	// register messages handler
	c.backend.AddRequestHandler(c.messagesHandler)

	// register now one time pre key handler
	c.backend.AddRequestHandler(c.oneTimePreKeysHandler)

	return c, nil
}

package db

import (
	"errors"
	"fmt"

	aes "github.com/Bit-Nation/panthalassa/crypto/aes"
	km "github.com/Bit-Nation/panthalassa/keyManager"
	storm "github.com/asdine/storm"
	sq "github.com/asdine/storm/q"
	ed25519 "golang.org/x/crypto/ed25519"
)

// message status
type Status uint

const (
	StatusSent           Status = 100
	StatusFailedToSend   Status = 200
	StatusDelivered      Status = 300
	StatusFailedToHandle Status = 400
	StatusPersisted      Status = 500
	DAppMessageVersion   uint   = 1
)

var statuses = map[Status]bool{
	StatusSent:           true,
	StatusFailedToSend:   true,
	StatusDelivered:      true,
	StatusFailedToHandle: true,
	StatusPersisted:      true,
}

// validate a given message
var ValidMessage = func(m Message) error {

	// validate id
	if m.ID == "" {
		return errors.New("invalid message id (empty string)")
	}

	// validate version
	if m.Version == 0 {
		return errors.New("invalid version - got 0")
	}

	// validate version
	if _, exist := statuses[m.Status]; !exist {
		return fmt.Errorf("invalid status: %d (is not registered)", m.Status)
	}

	// validate "type" of message
	if m.DApp == nil && len(m.Message) == 0 {
		return errors.New("got invalid message - dapp and message are both nil")
	}

	// validate DApp
	if m.DApp != nil {

		// validate DApp public key
		if len(m.DApp.DAppPublicKey) != 32 {
			return fmt.Errorf("invalid dapp public key of length %d", len(m.DApp.DAppPublicKey))
		}

	}

	// validate created at
	// must be greater then the max unix time stamp
	// in seconds since we need the micro second timestamp
	if m.CreatedAt <= 2147483647 {
		return errors.New("invalid created at - must be bigger than 2147483647")
	}

	// validate sender
	if len(m.Sender) != 32 && m.DApp == nil {
		return fmt.Errorf("invalid sender of length %d", len(m.Sender))
	}

	return nil

}

type DAppMessage struct {
	DAppPublicKey []byte
	Type          string
	Params        map[string]interface{}
	ShouldSend    bool
}

type Message struct {
	DBID             int `storm:"id,increment"`
	ID               string
	Version          uint
	Status           Status
	Received         bool
	DApp             *DAppMessage
	Message          []byte `json:"-"`
	PersistedMessage aes.CipherText
	CreatedAt        int64
	Sender           []byte `storm:"index"`
	// the UniqueID is a unix nano timestamp
	// it's only unique in the relation with a chat id
	UniqueMsgID int64 `storm:"index"`
	ChatID      int   `storm:"index"`
}

type Chat struct {
	ID int `storm:"id,increment"`
	// partner will only be filled if this is a private chat
	Partner             ed25519.PublicKey `storm:"index,unique"`
	UnreadMessages      bool
	db                  *storm.DB
	km                  *km.KeyManager
	postPersistListener []func(event MessagePersistedEvent)
}

func (c *Chat) GetMessage(msgID int64) (*Message, error) {
	q := c.db.Select(sq.And(
		sq.Eq("ChatID", c.ID),
		sq.Eq("UniqueMsgID", msgID),
	))

	// count messages
	amount, err := q.Count(&Message{})
	if err != nil {
		return nil, err
	}
	if amount == 0 {
		return nil, nil
	}

	// fetch message
	var m Message
	if err := q.First(&m); err != nil {
		return &Message{}, nil
	}
	return &m, nil

}

// Persist Message
func (c *Chat) PersistMessage(msg Message) error {
	msg.ChatID = c.ID
	ct, err := c.km.AESEncrypt(msg.Message)
	if err != nil {
		return err
	}
	msg.PersistedMessage = ct

	// find message id
	tried := 0
	lastUniqueMsgID := msg.CreatedAt
	for {
		q := c.db.Select(sq.And(
			sq.Eq("ChatID", c.ID),
			sq.Eq("UniqueMsgID", lastUniqueMsgID),
		))
		amount, err := q.Count(&Message{})
		if err != nil {
			return err
		}
		// break if there are no messages with that id
		if amount == 0 {
			break
		}
		// it's a attack (pretty sure) when we can't find an id after 10K attempts
		if tried == 10000 {
			return errors.New("can't persist message - couldn't create unique msg id")
		}
		tried++
	}
	msg.UniqueMsgID = lastUniqueMsgID

	// validate message
	if err := ValidMessage(msg); err != nil {
		return err
	}

	// persist message
	if err := c.db.Save(&msg); err != nil {
		return err
	}

	// emit all registered event listeners
	for _, l := range c.postPersistListener {
		go func() {
			l(MessagePersistedEvent{
				Chat:    *c,
				Message: msg,
			})
		}()
	}

	return nil
}

func (c *Chat) Messages(start int64, amount uint) ([]Message, error) {

	// default query should only select from the last
	q := c.db.Select(sq.And(sq.Eq("ChatID", c.ID)))

	// if start is not 0 we should jump to the start position by
	// excluding everything that is < start
	if start != 0 {
		q = c.db.Select(sq.And(sq.Eq("ChatID", c.ID), sq.Gte("UniqueMsgID", start)))
	}

	var messages []Message
	if err := q.OrderBy("UniqueMsgID").Reverse().Limit(int(amount)).Find(&messages); err != nil {
		return []Message{}, err
	}

	// decrypt messages
	for _, m := range messages {
		plainMessage, err := c.km.AESDecrypt(m.PersistedMessage)
		if err != nil {
			return []Message{}, err
		}
		m.Message = plainMessage
	}

	return messages, nil

}

type ChatStorage interface {
	GetChat(partner ed25519.PublicKey) (*Chat, error)
	CreateChat(partner ed25519.PublicKey) error
	AddListener(func(e MessagePersistedEvent))
	AllChats() ([]Chat, error)
	UnreadMessages(c Chat) error
}

type MessagePersistedEvent struct {
	Chat    Chat
	Message Message
}

type BoltChatStorage struct {
	db                  *storm.DB
	postPersistListener []func(event MessagePersistedEvent)
	km                  *km.KeyManager
}

func (s *BoltChatStorage) GetChat(partner ed25519.PublicKey) (*Chat, error) {

	// check if partner chat exist
	q := s.db.Select(sq.Eq("Partner", partner))
	amount, err := q.Count(&Chat{})
	if err != nil {
		return nil, err
	}
	if amount == 0 {
		return nil, nil
	}

	// fetch partner chat
	var c Chat
	if err := q.First(&c); err != nil {
		return nil, err
	}
	c.km = s.km
	c.db = s.db
	c.postPersistListener = s.postPersistListener
	return &c, nil
}

func (s *BoltChatStorage) CreateChat(partner ed25519.PublicKey) error {
	return s.db.Save(&Chat{
		Partner: partner,
	})
}

func (s *BoltChatStorage) UnreadMessages(c Chat) error {
	var fc Chat
	if err := s.db.One("ID", c.ID, &fc); err != nil {
		return err
	}
	fc.UnreadMessages = true
	return s.db.Save(&fc)
}

func (s *BoltChatStorage) AddListener(fn func(e MessagePersistedEvent)) {
	s.postPersistListener = append(s.postPersistListener, fn)
}

func (s *BoltChatStorage) AllChats() ([]Chat, error) {
	chats := new([]Chat)
	return *chats, s.db.All(chats)
}

func NewChatStorage(db *storm.DB, listeners []func(event MessagePersistedEvent), km *km.KeyManager) *BoltChatStorage {
	return &BoltChatStorage{
		db:                  db,
		postPersistListener: listeners,
		km:                  km,
	}
}

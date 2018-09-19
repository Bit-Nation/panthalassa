package db

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"

	aes "github.com/Bit-Nation/panthalassa/crypto/aes"
	km "github.com/Bit-Nation/panthalassa/keyManager"
	storm "github.com/asdine/storm"
	sq "github.com/asdine/storm/q"
	uuid "github.com/satori/go.uuid"
	ed25519 "golang.org/x/crypto/ed25519"
	"time"
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

	// validate version
	if m.Version == 0 {
		return errors.New("invalid version - got 0")
	}

	// validate version
	if _, exist := statuses[m.Status]; !exist {
		return fmt.Errorf("invalid status: %d (is not registered)", m.Status)
	}

	// validate "type" of message
	if m.DApp == nil && m.AddUserToChat == nil && len(m.Message) == 0 {
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

type AddUserToChat struct {
	Users  []ed25519.PublicKey
	ChatID []byte
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
	UniqueMsgID   int64 `storm:"index"`
	ChatID        int   `storm:"index"`
	AddUserToChat *AddUserToChat
	GroupChatID   []byte
}

type Chat struct {
	ID int `storm:"id,increment"`
	// partner will only be filled if this is a private chat
	Partner             ed25519.PublicKey `storm:"index,unique"`
	Partners            []ed25519.PublicKey
	UnreadMessages      bool
	GroupChatRemoteID   []byte
	db                  storm.Node
	km                  *km.KeyManager
	postPersistListener []func(event MessagePersistedEvent)
}

// @todo maybe more checks
func (c *Chat) IsGroupChat() bool {
	return len(c.GroupChatRemoteID) != 0
}

var nowAsUnix = func() int64 {
	return time.Now().UnixNano()
}

// save a raw message
func (c *Chat) SaveMessage(rawMessage []byte) error {

	id, err := uuid.NewV4()
	if err != nil {
		return err
	}

	senderStr, err := c.km.IdentityPublicKey()
	if err != nil {
		return err
	}
	sender, err := hex.DecodeString(senderStr)
	if err != nil {
		return err
	}

	msg := Message{
		ID:        id.String(),
		Message:   rawMessage,
		CreatedAt: nowAsUnix(),
		Version:   1,
		Status:    StatusPersisted,
		Sender:    sender,
	}

	// fetch chat
	return c.PersistMessage(msg)
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

	// decrypt original message
	m.Message, err = c.km.AESDecrypt(m.PersistedMessage)
	if err != nil {
		return nil, err
	}

	return &m, nil

}

// Persist Message struct
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
		lastUniqueMsgID++
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
		go func(l func(e MessagePersistedEvent)) {
			l(MessagePersistedEvent{
				Chat:    *c,
				Message: msg,
			})
		}(l)
	}

	return nil
}

func (c *Chat) Messages(start int64, amount uint) ([]*Message, error) {

	// default query should only select from the last
	q := c.db.Select(sq.And(sq.Eq("ChatID", c.ID)))

	// if start is not 0 we should jump to the start position by
	// excluding everything that is < start
	if start != 0 {
		q = c.db.Select(sq.And(sq.Eq("ChatID", c.ID), sq.Gte("UniqueMsgID", start)))
	}

	var messages []*Message
	if err := q.OrderBy("UniqueMsgID").Reverse().Limit(int(amount)).Find(&messages); err != nil {
		return nil, err
	}

	// decrypt messages
	for _, m := range messages {
		plainMessage, err := c.km.AESDecrypt(m.PersistedMessage)
		if err != nil {
			return nil, err
		}
		m.Message = plainMessage
	}

	return messages, nil

}

type ChatStorage interface {
	GetChatByPartner(pubKey ed25519.PublicKey) (*Chat, error)
	GetChat(chatID int) (*Chat, error)
	GetGroupChatByRemoteID(id []byte) (*Chat, error)
	// returned int is the chat ID
	CreateChat(partner ed25519.PublicKey) error
	CreateGroupChat(partners []ed25519.PublicKey) (int, error)
	CreateGroupChatFromMsg(createMessage *AddUserToChat) error
	AddListener(func(e MessagePersistedEvent))
	AllChats() ([]Chat, error)
	// set state of chat to unread messages
	UnreadMessages(c Chat) error
	// set state of chat all messages read
	ReadMessages(partner ed25519.PublicKey) error
}

type MessagePersistedEvent struct {
	Chat    Chat
	Message Message
}

type BoltChatStorage struct {
	db                  storm.Node
	postPersistListener []func(event MessagePersistedEvent)
	km                  *km.KeyManager
}

func (s *BoltChatStorage) CreateGroupChatFromMsg(createMessage *AddUserToChat) error {

	c := Chat{
		Partners:          createMessage.Users,
		GroupChatRemoteID: createMessage.ChatID,
	}

	return s.db.Save(&c)

}

func (s *BoltChatStorage) GetGroupChatByRemoteID(id []byte) (*Chat, error) {

	// make sure chats exist
	amount, err := s.db.Count(&Chat{})
	if err != nil {
		return nil, err
	}
	if amount == 0 {
		return nil, nil
	}

	c := &Chat{}
	return c, s.db.One("GroupChatRemoteID", id, c)

}

func (s *BoltChatStorage) GetChat(chatID int) (*Chat, error) {

	amountChats, err := s.db.Count(&Chat{})
	if err != nil {
		return nil, err
	}

	if amountChats == 0 {
		return nil, nil
	}

	chat := &Chat{}
	chat.km = s.km
	chat.db = s.db
	chat.postPersistListener = s.postPersistListener

	return chat, s.db.One("ID", amountChats, chat)

}

func (s *BoltChatStorage) GetChatByPartner(partner ed25519.PublicKey) (*Chat, error) {

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

func (s *BoltChatStorage) ReadMessages(partner ed25519.PublicKey) error {
	var fc Chat
	if err := s.db.One("Partner", partner, &fc); err != nil {
		return err
	}
	fc.UnreadMessages = false
	return s.db.Save(&fc)
}

func (s *BoltChatStorage) AddListener(fn func(e MessagePersistedEvent)) {
	s.postPersistListener = append(s.postPersistListener, fn)
}

func (s *BoltChatStorage) AllChats() ([]Chat, error) {
	chats := new([]Chat)
	return *chats, s.db.All(chats)
}

func (s *BoltChatStorage) CreateGroupChat(partners []ed25519.PublicKey) (int, error) {

	// remote id
	ri := make([]byte, 200)
	if _, err := rand.Read(ri); err != nil {
		return 0, err
	}

	// group chat shared secret id
	ssID := make([]byte, 200)
	if _, err := rand.Read(ssID); err != nil {
		return 0, err
	}

	c := &Chat{
		Partners:          partners,
		GroupChatRemoteID: ri,
	}

	if err := s.db.Save(c); err != nil {
		return 0, err
	}

	return c.ID, nil

}

func (c *Chat) AddChatPartners(partners []ed25519.PublicKey) error {

	for _, newPartner := range partners {
		exist := false
		for _, existingPartner := range c.Partners {
			if hex.EncodeToString(newPartner) == hex.EncodeToString(existingPartner) {
				exist = true
			}
		}
		if exist {
			c.Partners = append(c.Partners, newPartner)
		}
	}

	return c.db.Save(&c)

}

func NewChatStorage(db storm.Node, listeners []func(event MessagePersistedEvent), km *km.KeyManager) *BoltChatStorage {
	return &BoltChatStorage{
		db:                  db,
		postPersistListener: listeners,
		km:                  km,
	}
}

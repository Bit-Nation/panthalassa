package db

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/Bit-Nation/panthalassa/crypto/aes"
	km "github.com/Bit-Nation/panthalassa/keyManager"
	bolt "github.com/coreos/bbolt"
	uuid "github.com/satori/go.uuid"
	ed25519 "golang.org/x/crypto/ed25519"
)

var (
	privateChatBucketName = []byte("private_chat")
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

type ChatMessageStorage interface {
	PersistMessageToSend(partner ed25519.PublicKey, msg Message) error
	PersistReceivedMessage(partner ed25519.PublicKey, msg Message) error
	UpdateStatus(partner ed25519.PublicKey, msgID int64, newStatus Status) error
	AllChats() ([]ed25519.PublicKey, error)
	Messages(partner ed25519.PublicKey, start int64, amount uint) (map[int64]Message, error)
	AddListener(func(e MessagePersistedEvent))
	GetMessage(partner ed25519.PublicKey, messageID int64) (*Message, error)
}

type DAppMessage struct {
	DAppPublicKey []byte                 `json:"dapp_public_key"`
	Type          string                 `json:"type"`
	Params        map[string]interface{} `json:"params"`
	ShouldSend    bool                   `json:"should_send"`
}

type Message struct {
	ID         string       `json:"message_id"`
	Version    uint         `json:"version"`
	Status     Status       `json:"status"`
	Received   bool         `json:"received"`
	DApp       *DAppMessage `json:"dapp"`
	Message    []byte       `json:"message"`
	CreatedAt  int64        `json:"created_at"`
	Sender     []byte       `json:"sender"`
	DatabaseID int64        `json:"db_id"`
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
		if len(m.DApp.DAppPublicKey) != 0 {
			return fmt.Errorf("invalid dapp public key of length %d", len(m.DApp.DAppPublicKey))
		}

	}

	// validate created at
	// must be greater then the max unix time stamp
	// in seconds since we need the micro second timestamp
	if m.CreatedAt <= 2147483647 {
		return errors.New("invalid created at - must be bigger then 2147483647")
	}

	// validate sender
	if len(m.Sender) != 32 {
		return fmt.Errorf("invalid sender of length %d", len(m.Sender))
	}

	return nil

}

type MessagePersistedEvent struct {
	Partner     ed25519.PublicKey
	Message     Message
	DBMessageID int64
}

type BoltChatMessageStorage struct {
	db                  *bolt.DB
	postPersistListener []func(event MessagePersistedEvent)
	km                  *km.KeyManager
}

func NewChatMessageStorage(db *bolt.DB, listeners []func(event MessagePersistedEvent), km *km.KeyManager) *BoltChatMessageStorage {
	return &BoltChatMessageStorage{
		db:                  db,
		postPersistListener: listeners,
		km:                  km,
	}
}

func (s *BoltChatMessageStorage) persistMessage(partner ed25519.PublicKey, msg Message) error {

	// set version of message
	msg.Version = DAppMessageVersion

	// validate message
	if err := ValidMessage(msg); err != nil {
		return err
	}

	// persist message
	return s.db.Update(func(tx *bolt.Tx) error {

		// private chat bucket
		privChatBucket, err := tx.CreateBucketIfNotExists(privateChatBucketName)
		if err != nil {
			return err
		}

		// create partner chat bucket
		partnerBucket, err := privChatBucket.CreateBucketIfNotExists(partner)
		if err != nil {
			return err
		}

		// turn created at into bytes
		// createdAtMsgID is the id used for the database
		createdAtMsgID := make([]byte, 8)
		binary.BigEndian.PutUint64(createdAtMsgID, uint64(msg.CreatedAt))

		// make sure it is not taken and adjust the time indexed timestamp
		tried := 0
		for {
			fetchedMsg := partnerBucket.Get(createdAtMsgID)
			if fetchedMsg == nil || tried == 1000 {
				break
			}
			tried++
			binary.BigEndian.PutUint64(createdAtMsgID, uint64(msg.CreatedAt+1))
		}

		// marshal message
		rawMessage, err := json.Marshal(msg)
		if err != nil {
			return err
		}

		// encrypt raw proto message
		encryptedMessage, err := s.km.AESEncrypt(rawMessage)
		if err != nil {
			return err
		}

		// marshaled encrypted message
		rawEncryptedMessage, err := encryptedMessage.Marshal()
		if err != nil {
			return err
		}

		// tell listeners that we persisted the message
		tx.OnCommit(func() {
			for _, listener := range s.postPersistListener {
				go listener(MessagePersistedEvent{
					Partner:     partner,
					Message:     msg,
					DBMessageID: int64(binary.BigEndian.Uint64(createdAtMsgID)),
				})
			}
		})

		return partnerBucket.Put(createdAtMsgID, rawEncryptedMessage)

	})
}

// fetch all chat partners
func (s *BoltChatMessageStorage) AllChats() ([]ed25519.PublicKey, error) {
	chats := []ed25519.PublicKey{}
	err := s.db.View(func(tx *bolt.Tx) error {

		// all private chats
		privateChats := tx.Bucket(privateChatBucketName)
		if privateChats == nil {
			return nil
		}

		return privateChats.ForEach(func(k, _ []byte) error {
			if len(k) == 32 {
				chats = append(chats, k)
			}
			return nil
		})
	})
	return chats, err
}

func (s *BoltChatMessageStorage) Messages(partner ed25519.PublicKey, start int64, amount uint) (map[int64]Message, error) {

	messages := map[int64]Message{}

	err := s.db.View(func(tx *bolt.Tx) error {

		// private chats
		privChatsBucket := tx.Bucket(privateChatBucketName)
		if privChatsBucket == nil {
			return nil
		}

		// partner chat bucket
		partnerBucket := privChatsBucket.Bucket(partner)
		if partnerBucket == nil {
			return nil
		}

		cursor := partnerBucket.Cursor()
		var msgID int64
		var rawMsg []byte

		// jump to position
		if start == 0 {
			key, value := cursor.Last()
			msgID = int64(binary.BigEndian.Uint64(key))
			rawMsg = value
		} else {
			startBytes := make([]byte, 8)
			binary.BigEndian.PutUint64(startBytes, uint64(start))
			key, value := cursor.Seek(startBytes)
			msgID = int64(binary.BigEndian.Uint64(key))
			rawMsg = value
		}

		// unmarshal message
		msg := Message{}
		if err := json.Unmarshal(rawMsg, &msg); err != nil {
			return err
		}
		messages[msgID] = msg

		currentAmount := amount - 1
		for {
			if currentAmount == 0 {
				break
			}
			currentAmount--
			key, rawMsg := cursor.Prev()
			if key == nil {
				break
			}
			msg := Message{}
			if err := json.Unmarshal(rawMsg, &msg); err != nil {
				return err
			}
			messages[int64(binary.BigEndian.Uint64(key))] = msg
		}

		return nil
	})

	return messages, err

}

// fetch message by it's partner and database id
// will return nil if the message doesn't exist
func (s *BoltChatMessageStorage) GetMessage(partner ed25519.PublicKey, dbID int64) (*Message, error) {
	var msg *Message
	msg = nil
	err := s.db.View(func(tx *bolt.Tx) error {

		// private chats bucket
		privateChats, err := tx.CreateBucketIfNotExists(privateChatBucketName)
		if err != nil {
			return err
		}

		// bucket with chat of partner
		partnerMessages, err := privateChats.CreateBucketIfNotExists(partner)
		if err != nil {
			return err
		}

		// turn numeric message id into byte message id
		byteMsgID := make([]byte, 8)
		binary.BigEndian.PutUint64(byteMsgID, uint64(dbID))

		// fetch encrypted message
		rawEncryptedMessage := partnerMessages.Get(byteMsgID)
		if rawEncryptedMessage == nil {
			return fmt.Errorf("coulnd't fetch message for partner: %x and message id: %d", partner, messageID)
		}

		// decrypt message
		ct, err := aes.Unmarshal(rawEncryptedMessage)
		if err != nil {
			return err
		}
		rawPlainMessage, err := s.km.AESDecrypt(ct)
		if err != nil {
			return err
		}

		// unmarshal database message
		m := Message{}
		if err := json.Unmarshal(rawPlainMessage, &m); err != nil {
			return err
		}
		msg = &m

		return nil

	})
	return msg, err
}

// add listener
func (s *BoltChatMessageStorage) AddListener(fn func(e MessagePersistedEvent)) {
	s.postPersistListener = append(s.postPersistListener, fn)
}

func (s *BoltChatMessageStorage) PersistMessageToSend(partner ed25519.PublicKey, msg Message) error {
	id, err := uuid.NewV4()
	if err != nil {
		return err
	}
	msg.ID = id.String()
	msg.Received = false
	msg.Status = StatusPersisted
	return s.persistMessage(partner, msg)
}

func (s *BoltChatMessageStorage) PersistReceivedMessage(partner ed25519.PublicKey, msg Message) error {
	msg.Status = StatusPersisted
	msg.Received = true
	return s.persistMessage(partner, msg)
}

func (s *BoltChatMessageStorage) UpdateStatus(partner ed25519.PublicKey, msgID string, newStatus Status) error {
	return errors.New("currently not implemented")
}

package db

import (
	"encoding/binary"
	"encoding/json"
	"errors"

	aes "github.com/Bit-Nation/panthalassa/crypto/aes"
	km "github.com/Bit-Nation/panthalassa/keyManager"
	bpb "github.com/Bit-Nation/protobuffers"
	bolt "github.com/coreos/bbolt"
	proto "github.com/gogo/protobuf/proto"
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
)

type ChatMessageStorage interface {
	PersistMessageToSend(partner ed25519.PublicKey, msg bpb.PlainChatMessage) error
	PersistReceivedMessage(partner ed25519.PublicKey, msg bpb.PlainChatMessage) error
	UpdateStatus(partner ed25519.PublicKey, msgID string, newStatus Status) error
	AllChats() ([]ed25519.PublicKey, error)
	Messages(partner ed25519.PublicKey, start int64, amount uint) (map[int64]Message, error)
}

type Message struct {
	ID string `json:"message_id"`
	// encrypted message
	Message aes.CipherText `json:"message"`
	Version uint32         `json:"version"`
	Status  Status         `json:"status"`
	// is a received message
	Received bool `json:"received"`
}

type PlainMessagePersistedEvent struct {
	Msg     bpb.PlainChatMessage
	Partner ed25519.PublicKey
	DBMsg   Message
	MsgID   int64
}

type BoltChatMessageStorage struct {
	db                  *bolt.DB
	postPersistListener []chan PlainMessagePersistedEvent
	km                  *km.KeyManager
}

func NewChatMessageStorage(db *bolt.DB, listeners []chan PlainMessagePersistedEvent, km *km.KeyManager) *BoltChatMessageStorage {
	return &BoltChatMessageStorage{
		db:                  db,
		postPersistListener: listeners,
		km:                  km,
	}
}

func (s *BoltChatMessageStorage) persistMessage(partner ed25519.PublicKey, msg bpb.PlainChatMessage, received bool, status Status) error {
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
		createdAt := make([]byte, 8)
		binary.BigEndian.PutUint64(createdAt, uint64(msg.CreatedAt))

		// make sure it is not taken and adjust the time indexed timestamp
		tried := 0
		for {
			fetchedMsg := partnerBucket.Get(createdAt)
			if fetchedMsg == nil || tried == 1000 {
				break
			}
			tried++
			binary.BigEndian.PutUint64(createdAt, uint64(msg.CreatedAt+1))
		}

		// marshal message
		rawProtoMsg, err := proto.Marshal(&msg)
		if err != nil {
			return err
		}

		// encrypt raw proto message
		encryptedMessage, err := s.km.AESEncrypt(rawProtoMsg)
		if err != nil {
			return err
		}

		// database message
		dbMessage := Message{
			ID:       msg.MessageID,
			Message:  encryptedMessage,
			Version:  1,
			Status:   status,
			Received: received,
		}

		// marshal DB message
		rawDBMessage, err := json.Marshal(dbMessage)
		if err != nil {
			return err
		}

		// tell listeners that we persisted the message
		tx.OnCommit(func() {
			for _, listener := range s.postPersistListener {
				listener <- PlainMessagePersistedEvent{
					Msg:     msg,
					Partner: partner,
					DBMsg:   dbMessage,
					MsgID:   int64(binary.BigEndian.Uint64(createdAt)),
				}
			}
		})

		return partnerBucket.Put(createdAt, rawDBMessage)

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

func (s *BoltChatMessageStorage) PersistMessageToSend(partner ed25519.PublicKey, msg bpb.PlainChatMessage) error {
	return s.persistMessage(partner, msg, false, StatusPersisted)
}

func (s *BoltChatMessageStorage) PersistReceivedMessage(partner ed25519.PublicKey, msg bpb.PlainChatMessage) error {
	return s.persistMessage(partner, msg, true, StatusPersisted)
}

func (s *BoltChatMessageStorage) UpdateStatus(partner ed25519.PublicKey, msgID string, newStatus Status) error {
	return errors.New("currently not implemented")
}

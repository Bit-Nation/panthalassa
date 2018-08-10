package db

import (
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"testing"

	aes "github.com/Bit-Nation/panthalassa/crypto/aes"
	bolt "github.com/coreos/bbolt"
	require "github.com/stretchr/testify/require"
	ed25519 "golang.org/x/crypto/ed25519"
)

func TestValidMessage(t *testing.T) {

	type testVector struct {
		expectedError string
		message       Message
	}

	tests := []testVector{
		testVector{
			expectedError: "invalid message id (empty string)",
			message:       Message{},
		},
		testVector{
			expectedError: "invalid version - got 0",
			message: Message{
				ID: "-",
			},
		},
		testVector{
			expectedError: "invalid status: 35939 (is not registered)",
			message: Message{
				ID:      "-",
				Version: 1,
				Status:  35939,
			},
		},
		testVector{
			expectedError: "got invalid message - dapp and message are both nil",
			message: Message{
				ID:      "-",
				Version: 1,
				Status:  100,
			},
		},
		testVector{
			expectedError: "invalid dapp public key of length 9",
			message: Message{
				ID:      "-",
				Version: 1,
				Status:  100,
				DApp: &DAppMessage{
					DAppPublicKey: []byte("too short"),
				},
			},
		},
		testVector{
			expectedError: "invalid created at - must be bigger than 2147483647",
			message: Message{
				ID:      "-",
				Version: 1,
				Status:  100,
				Message: []byte("message"),
			},
		},
		testVector{
			expectedError: "invalid sender of length 9",
			message: Message{
				ID:        "-",
				Version:   1,
				CreatedAt: 2147483648,
				Status:    100,
				Message:   []byte("message"),
				Sender:    []byte("too short"),
			},
		},
	}

	for _, v := range tests {
		require.EqualError(t, ValidMessage(v.message), v.expectedError)
	}

}

// test if message validation is really called
func TestBoltChatMessageStorage_persistMessageValidation(t *testing.T) {

	// setup
	db := createDB()
	km := createKeyManager()
	partner, _, err := ed25519.GenerateKey(rand.Reader)
	require.Nil(t, err)

	// override message validation
	oldValidMessage := ValidMessage
	ValidMessage = func(m Message) error {
		return errors.New("hello from ValidMessage mock")
	}

	storage := NewChatMessageStorage(db, []func(event MessagePersistedEvent){}, km)
	err = storage.persistMessage(partner, Message{})
	require.EqualError(t, err, "hello from ValidMessage mock")

	ValidMessage = oldValidMessage

}

func TestBoltChatMessageStorage_persistSuccess(t *testing.T) {

	// setup
	db := createDB()
	km := createKeyManager()
	partner, _, err := ed25519.GenerateKey(rand.Reader)
	require.Nil(t, err)

	msgToPersist := Message{
		ID:        "-",
		Message:   []byte("hi"),
		CreatedAt: 2147483648,
		Sender:    partner,
		Status:    StatusPersisted,
		Received:  true,
	}

	// message assertion
	messageAssertion := func(msg Message) {

		// assertions on message
		require.Equal(t, msgToPersist.ID, msg.ID)
		require.Equal(t, msgToPersist.Message, msg.Message)
		require.Equal(t, msgToPersist.Received, msg.Received)
		require.Equal(t, msgToPersist.Status, msg.Status)
		require.Equal(t, msgToPersist.CreatedAt, msg.CreatedAt)
		require.Equal(t, uint(1), msg.Version)
		require.Equal(t, msgToPersist.Sender, msg.Sender)
		require.Equal(t, msgToPersist.DApp, msg.DApp)
		require.Equal(t, int64(2147483648), msg.DatabaseID)

	}

	// listener
	calledListener := false
	listeners := []func(event MessagePersistedEvent){
		func(event MessagePersistedEvent) {
			messageAssertion(event.Message)
			require.Equal(t, int64(2147483648), event.DBMessageID)
			require.Equal(t, partner, event.Partner)
			calledListener = true
		},
	}

	// persist message
	storage := NewChatMessageStorage(db, listeners, km)
	err = storage.persistMessage(partner, msgToPersist)
	require.Nil(t, err)

	err = db.View(func(tx *bolt.Tx) error {

		// private chats bucket
		privChats := tx.Bucket(privateChatBucketName)
		require.NotNil(t, privChats)

		// private chat with partner
		partnerPrivChat := privChats.Bucket(partner)
		require.NotNil(t, partnerPrivChat)

		// id (key) for message in DB
		id := make([]byte, 8)
		binary.BigEndian.PutUint64(id, uint64(msgToPersist.CreatedAt))

		// fetch message
		message := partnerPrivChat.Get(id)
		require.NotNil(t, message)

		// message into cipher text
		encryptedCipherText, err := aes.Unmarshal(message)
		require.Nil(t, err)

		// decrypt message
		rawMessage, err := km.AESDecrypt(encryptedCipherText)
		require.Nil(t, err)

		// unmarshal message
		msg := Message{}
		require.Nil(t, json.Unmarshal(rawMessage, &msg))
		messageAssertion(msg)

		return nil
	})
	require.Nil(t, err)

}

func TestBoltChatMessageStorage_AllChats(t *testing.T) {

	// setup
	db := createDB()
	km := createKeyManager()
	partnerOne, _, err := ed25519.GenerateKey(rand.Reader)
	partnerTwo, _, err := ed25519.GenerateKey(rand.Reader)
	require.Nil(t, err)
	storage := NewChatMessageStorage(db, []func(event MessagePersistedEvent){}, km)

	require.Nil(t, storage.PersistMessageToSend(partnerOne, Message{Message: []byte("hi @partner one")}))
	require.Nil(t, storage.PersistMessageToSend(partnerTwo, Message{Message: []byte("hi @partner two")}))

	// fetch all chat partners and make sure we get the once we expect
	partners, err := storage.AllChats()
	require.Nil(t, err)
	require.Equal(t, 2, len(partners))
	require.True(t, hex.EncodeToString(partners[0]) == hex.EncodeToString(partnerOne) || hex.EncodeToString(partners[0]) == hex.EncodeToString(partnerTwo))
	require.True(t, hex.EncodeToString(partners[1]) == hex.EncodeToString(partnerOne) || hex.EncodeToString(partners[1]) == hex.EncodeToString(partnerTwo))

	// assert message for partner one
	messages, err := storage.Messages(partnerOne, 0, 10)
	require.Nil(t, err)
	for _, msg := range messages {
		require.Equal(t, []byte("hi @partner one"), msg.Message)
	}

	// assert message for partner two
	messages, err = storage.Messages(partnerTwo, 0, 10)
	require.Nil(t, err)
	for _, msg := range messages {
		require.Equal(t, []byte("hi @partner two"), msg.Message)
	}

}

func TestBoltChatMessageStorage_Messages(t *testing.T) {

	// setup
	db := createDB()
	km := createKeyManager()
	partner, _, err := ed25519.GenerateKey(rand.Reader)
	require.Nil(t, err)
	storage := NewChatMessageStorage(db, []func(event MessagePersistedEvent){}, km)

	// make sure persisted message is fetched
	require.Nil(t, storage.PersistMessageToSend(partner, Message{Message: []byte("hi")}))
	messages, err := storage.Messages(partner, 0, 10)
	require.Nil(t, err)
	require.Equal(t, []byte("hi"), messages[0].Message)

	// persist another message and fetch it
	require.Nil(t, storage.PersistMessageToSend(partner, Message{Message: []byte("another message")}))
	messages, err = storage.Messages(partner, 0, 10)
	require.Nil(t, err)
	require.Equal(t, []byte("hi"), messages[0].Message)
	require.Equal(t, []byte("another message"), messages[1].Message)

}

func TestBoltChatMessageStorage_GetMessage(t *testing.T) {

	// setup
	db := createDB()
	km := createKeyManager()
	partner, _, err := ed25519.GenerateKey(rand.Reader)
	require.Nil(t, err)
	storage := NewChatMessageStorage(db, []func(event MessagePersistedEvent){}, km)

	require.Nil(t, storage.PersistMessageToSend(partner, Message{Message: []byte("hi there")}))
	messages, err := storage.Messages(partner, 0, 10)
	require.Nil(t, err)

	msg, err := storage.GetMessage(partner, messages[0].DatabaseID)
	require.Nil(t, err)
	require.Equal(t, []byte("hi there"), msg.Message)

}

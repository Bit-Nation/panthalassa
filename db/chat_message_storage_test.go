package db

import (
	"crypto/rand"
	"encoding/binary"
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
		CreatedAt: 2147483647,
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
		require.Equal(t, int64(2147483647), msg.DatabaseID)

	}

	// listener
	calledListener := false
	listeners := []func(event MessagePersistedEvent){
		func(event MessagePersistedEvent) {
			messageAssertion(event.Message)
			require.Equal(t, int64(2147483647), event.DBMessageID)
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

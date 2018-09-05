package db

import (
	"crypto/rand"
	"encoding/hex"
	"testing"
	"time"

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

func TestChatMessages(t *testing.T) {

	storm := createStorm()
	km := createKeyManager()
	partner, _, err := ed25519.GenerateKey(rand.Reader)
	require.Nil(t, err)
	myIDKeyStr, err := km.IdentityPublicKey()
	require.Nil(t, err)
	myIDKey, err := hex.DecodeString(myIDKeyStr)
	require.Nil(t, err)

	chatStor := NewChatStorage(storm, []func(e MessagePersistedEvent){}, km)

	// create chat
	require.Nil(t, chatStor.CreateChat(partner))
	chat, err := chatStor.GetChat(partner)
	require.Nil(t, err)
	require.NotNil(t, chat)

	// first message
	err = chat.PersistMessage(Message{
		ID:        "msg one",
		Version:   1,
		Status:    StatusPersisted,
		Received:  true,
		Message:   []byte("message one"),
		CreatedAt: time.Now().UnixNano(),
		Sender:    myIDKey,
	})
	require.Nil(t, err)

	// second message
	err = chat.PersistMessage(Message{
		ID:        "msg two",
		Version:   1,
		Status:    StatusPersisted,
		Received:  true,
		Message:   []byte("message two"),
		CreatedAt: time.Now().UnixNano(),
		Sender:    myIDKey,
	})
	require.Nil(t, err)

	messages, err := chat.Messages(0, 2)
	require.Nil(t, err)

	require.Equal(t, "message one", string(messages[1].Message))
	require.Equal(t, "message two", string(messages[0].Message))

}

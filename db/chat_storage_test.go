package db

import (
	"testing"

	require "github.com/stretchr/testify/require"
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

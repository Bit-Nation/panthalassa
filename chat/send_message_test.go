package chat

import (
	"errors"
	"testing"

	db "github.com/Bit-Nation/panthalassa/db"
	bpb "github.com/Bit-Nation/protobuffers"
	x3dh "github.com/Bit-Nation/x3dh"
	require "github.com/stretchr/testify/require"
	ed25519 "golang.org/x/crypto/ed25519"
)

func TestChat_SendMessageFailToFetchPreKeyBundle(t *testing.T) {

	msgStorage := testMessageStorage{
		updateStatus: func(partner ed25519.PublicKey, msgID string, newStatus db.Status) error {
			require.Equal(t, ed25519.PublicKey{1}, partner)
			require.Equal(t, "i am the message ID", msgID)
			require.Equal(t, db.StatusFailedToSend, newStatus)
			return nil
		},
	}

	backend := testBackend{
		fetchPreKeyBundle: func(userIDPubKey ed25519.PublicKey) (x3dh.PreKeyBundle, error) {
			return nil, errors.New("i am a test error - failed to fetch pre key bundle")
		},
	}

	sharedSecretStore := testSharedSecretStorage{
		hasAny: func(key ed25519.PublicKey) (bool, error) {
			return false, nil
		},
	}

	c := Chat{
		messageDB:        &msgStorage,
		backend:          &backend,
		sharedSecStorage: &sharedSecretStore,
	}

	err := c.SendMessage(ed25519.PublicKey{1}, bpb.PlainChatMessage{MessageID: "i am the message ID"})
	require.EqualError(t, err, "i am a test error - failed to fetch pre key bundle")

}

func TestChat_SendMessageX3dhError(t *testing.T) {

	msgStorage := testMessageStorage{
		updateStatus: func(partner ed25519.PublicKey, msgID string, newStatus db.Status) error {
			require.Equal(t, ed25519.PublicKey{1}, partner)
			require.Equal(t, "i am the message ID", msgID)
			require.Equal(t, db.StatusFailedToSend, newStatus)
			return nil
		},
	}

	backend := testBackend{
		fetchPreKeyBundle: func(userIDPubKey ed25519.PublicKey) (x3dh.PreKeyBundle, error) {
			return testPreKeyBundle{validSignature: false}, nil
		},
	}

	sharedSecretStore := testSharedSecretStorage{
		hasAny: func(key ed25519.PublicKey) (bool, error) {
			return false, nil
		},
	}

	c := Chat{
		messageDB:        &msgStorage,
		backend:          &backend,
		sharedSecStorage: &sharedSecretStore,
	}

	err := c.SendMessage(ed25519.PublicKey{1}, bpb.PlainChatMessage{MessageID: "i am the message ID"})
	require.EqualError(t, err, "the signature of the received pre key bundle is invalid")

}

// test if we exit correct if we should have a shared secret
// but we don't
func TestChat_SendMessageNoSharedSecretThoWeShould(t *testing.T) {
	msgStorage := testMessageStorage{
		updateStatus: func(partner ed25519.PublicKey, msgID string, newStatus db.Status) error {
			require.Equal(t, ed25519.PublicKey{1}, partner)
			require.Equal(t, "i am the message ID", msgID)
			require.Equal(t, db.StatusFailedToSend, newStatus)
			return nil
		},
	}

	backend := testBackend{
		fetchPreKeyBundle: func(userIDPubKey ed25519.PublicKey) (x3dh.PreKeyBundle, error) {
			return testPreKeyBundle{validSignature: false}, nil
		},
	}

	sharedSecretStore := testSharedSecretStorage{
		hasAny: func(key ed25519.PublicKey) (bool, error) {
			return true, nil
		},
		getYoungest: func(key ed25519.PublicKey) (db.SharedSecret, error) {
			return db.SharedSecret{}, errors.New("no shared secret found test error")
		},
	}

	c := Chat{
		messageDB:        &msgStorage,
		backend:          &backend,
		sharedSecStorage: &sharedSecretStore,
	}

	err := c.SendMessage(ed25519.PublicKey{1}, bpb.PlainChatMessage{MessageID: "i am the message ID"})
	require.EqualError(t, err, "no shared secret found test error")
}

func TestChat_SendMessage(t *testing.T) {

}

func TestChat_SendMessageWithX3dhParameters(t *testing.T) {

}

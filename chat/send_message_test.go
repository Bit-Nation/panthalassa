package chat

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"testing"

	preKey "github.com/Bit-Nation/panthalassa/chat/prekey"
	db "github.com/Bit-Nation/panthalassa/db"
	bpb "github.com/Bit-Nation/protobuffers"
	x3dh "github.com/Bit-Nation/x3dh"
	proto "github.com/gogo/protobuf/proto"
	require "github.com/stretchr/testify/require"
	dr "github.com/tiabc/doubleratchet"
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

	backend := testBackend{}

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

// Alice sends a message to bob
func TestChat_SendMessage(t *testing.T) {

	kmAlice := createKeyManager()
	idPubKeyAliceStr, err := kmAlice.IdentityPublicKey()
	require.Nil(t, err)
	rawIdPubAliceBob, err := hex.DecodeString(idPubKeyAliceStr)
	require.Nil(t, err)

	kmBob := createKeyManager()
	idPubKeyBobStr, err := kmBob.IdentityPublicKey()
	require.Nil(t, err)
	rawIdPubKeyBob, err := hex.DecodeString(idPubKeyBobStr)
	require.Nil(t, err)

	// bob signed pre key
	curve := x3dh.NewCurve25519(rand.Reader)
	drKeyPair, err := curve.GenerateKeyPair()
	require.Nil(t, err)
	signedPreKeyBob := preKey.PreKey{}
	signedPreKeyBob.PrivateKey = drKeyPair.PrivateKey
	signedPreKeyBob.PublicKey = drKeyPair.PublicKey
	require.Nil(t, signedPreKeyBob.Sign(*kmBob))

	plainMsgToSend := bpb.PlainChatMessage{
		MessageID: "i am the message ID of the message to send",
		Message:   []byte("hi there"),
		CreatedAt: 4,
	}

	msgStorage := testMessageStorage{
		updateStatus: func(partner ed25519.PublicKey, msgID string, newStatus db.Status) error {
			require.Equal(t, hex.EncodeToString(rawIdPubKeyBob), hex.EncodeToString(partner))
			require.Equal(t, "i am the message ID of the message to send", msgID)
			require.Equal(t, db.StatusSent, newStatus)
			return nil
		},
	}

	calledBackend := false
	backend := testBackend{
		submitMessage: func(msg bpb.ChatMessage) error {

			// all of those checks relate to the if a
			// shared secret was accepted or not
			require.Equal(t, []byte(nil), msg.OneTimePreKey)
			require.Equal(t, []byte(nil), msg.SignedPreKey)
			require.Equal(t, []byte(nil), msg.EphemeralKey)
			require.Equal(t, []byte(nil), msg.EphemeralKeySignature)
			require.Equal(t, int64(0), msg.SharedSecretCreationDate)

			// make sure that the receiver and sender are correct
			require.Equal(t, rawIdPubKeyBob, msg.Receiver)
			require.Equal(t, rawIdPubAliceBob, msg.Sender)
			require.Equal(t, "i am the message ID of the message to send", string(msg.MessageID))

			// get double ratchet message from protobuf
			var dh dr.Key
			copy(dh[:], msg.Message.DoubleRatchetPK)
			drMsg := dr.Message{
				Header: dr.MessageHeader{
					DH: dh,
					N:  msg.Message.N,
					PN: msg.Message.Pn,
				},
				Ciphertext: msg.Message.CipherText,
			}

			// make sure encrypted messages is correct
			// by decrypting it and comparing the
			// original plain message with the received message
			drSession, err := dr.New([32]byte{1}, drDhPair{
				pub:  signedPreKeyBob.PublicKey,
				priv: signedPreKeyBob.PrivateKey,
			})
			decryptedRawMessage, err := drSession.RatchetDecrypt(drMsg, nil)
			require.Nil(t, err)
			plainMsg := bpb.PlainChatMessage{}
			require.Nil(t, proto.Unmarshal(decryptedRawMessage, &plainMsg))
			require.Equal(t, plainMsgToSend, plainMsg)

			calledBackend = true
			return nil
		},
	}

	sharedSecretStore := testSharedSecretStorage{
		hasAny: func(key ed25519.PublicKey) (bool, error) {
			return true, nil
		},
		getYoungest: func(key ed25519.PublicKey) (db.SharedSecret, error) {
			return db.SharedSecret{X3dhSS: x3dh.SharedSecret{1}, Accepted: true}, nil
		},
	}

	userStorage := testUserStorage{
		hasSignedPreKey: func(idKey ed25519.PublicKey) (bool, error) {
			return true, nil
		},
		getSignedPreKey: func(idKey ed25519.PublicKey) (preKey.PreKey, error) {
			// at this point we would return
			// the signed pre key of our chat partner
			return signedPreKeyBob, nil
		},
	}

	c := Chat{
		messageDB:        &msgStorage,
		backend:          &backend,
		sharedSecStorage: &sharedSecretStore,
		km:               kmAlice,
		drKeyStorage:     &dr.KeysStorageInMemory{},
		userStorage:      &userStorage,
	}

	err = c.SendMessage(rawIdPubKeyBob, plainMsgToSend)
	require.Nil(t, err)
	require.True(t, calledBackend)

}

func TestChat_SendMessageWithX3dhParameters(t *testing.T) {

}

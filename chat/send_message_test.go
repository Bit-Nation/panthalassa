package chat

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"testing"
	"time"

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
		updateStatus: func(partner ed25519.PublicKey, msgID int64, newStatus db.Status) error {
			require.Equal(t, ed25519.PublicKey{1}, partner)
			require.Equal(t, int64(2147483648), msgID)
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
		km:               createKeyManager(),
	}

	err := c.SendMessage(ed25519.PublicKey{1}, db.Message{
		ID:         "i am the message ID",
		Version:    1,
		Status:     300,
		Message:    []byte("my message"),
		CreatedAt:  2147483648,
		Sender:     make([]byte, 32),
		DatabaseID: 2147483648,
	},
	)
	require.EqualError(t, err, "i am a test error - failed to fetch pre key bundle")

}

func TestChat_SendMessageX3dhError(t *testing.T) {

	msgStorage := testMessageStorage{
		updateStatus: func(partner ed25519.PublicKey, msgID int64, newStatus db.Status) error {
			require.Equal(t, ed25519.PublicKey{1}, partner)
			require.Equal(t, int64(2147483648), msgID)
			require.Equal(t, db.StatusFailedToSend, newStatus)
			return nil
		},
	}

	backend := testBackend{
		fetchPreKeyBundle: func(userIDPubKey ed25519.PublicKey) (x3dh.PreKeyBundle, error) {
			return testPreKeyBundle{validSignature: func() (bool, error) {
				return false, nil
			}}, nil
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
		km:               createKeyManager(),
	}

	err := c.SendMessage(ed25519.PublicKey{1}, db.Message{
		ID:         "i am the message ID",
		Version:    1,
		Status:     300,
		Message:    []byte("my message"),
		CreatedAt:  2147483648,
		Sender:     make([]byte, 32),
		DatabaseID: 2147483648,
	})
	require.EqualError(t, err, "the signature of the received pre key bundle is invalid")

}

// test if we exit correct if we should have a shared secret
// but we don't
func TestChat_SendMessageNoSharedSecretThoWeShould(t *testing.T) {
	msgStorage := testMessageStorage{
		updateStatus: func(partner ed25519.PublicKey, msgID int64, newStatus db.Status) error {
			require.Equal(t, ed25519.PublicKey{1}, partner)
			require.Equal(t, int64(2147483648), msgID)
			require.Equal(t, db.StatusFailedToSend, newStatus)
			return nil
		},
	}

	backend := testBackend{}

	sharedSecretStore := testSharedSecretStorage{
		hasAny: func(key ed25519.PublicKey) (bool, error) {
			return true, nil
		},
		getYoungest: func(key ed25519.PublicKey) (*db.SharedSecret, error) {
			return nil, errors.New("no shared secret found test error")
		},
	}

	c := Chat{
		messageDB:        &msgStorage,
		backend:          &backend,
		sharedSecStorage: &sharedSecretStore,
		km:               createKeyManager(),
	}

	err := c.SendMessage(ed25519.PublicKey{1}, db.Message{
		ID:         "i am the message ID",
		Version:    1,
		Status:     300,
		Message:    []byte("my message"),
		CreatedAt:  2147483648,
		Sender:     make([]byte, 32),
		DatabaseID: 2147483648,
	})
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

	plainMsgToSend := db.Message{
		ID:         "i am the message ID of the message to send",
		Version:    1,
		Status:     300,
		Message:    []byte("my message"),
		CreatedAt:  2147483648,
		Sender:     make([]byte, 32),
		DatabaseID: 2147483648,
	}

	msgStorage := testMessageStorage{
		updateStatus: func(partner ed25519.PublicKey, msgID int64, newStatus db.Status) error {
			require.Equal(t, hex.EncodeToString(rawIdPubKeyBob), hex.EncodeToString(partner))
			require.Equal(t, int64(2147483648), msgID)
			require.Equal(t, db.StatusSent, newStatus)
			return nil
		},
	}

	sharedSecretBaseID := make([]byte, 32)
	_, err = rand.Read(sharedSecretBaseID)
	require.Nil(t, err)

	calledBackend := false
	backend := testBackend{
		submitMessages: func(messages []*bpb.ChatMessage) error {

			msg := messages[0]
			require.NotNil(t, msg)

			// all of those checks relate to the if a
			// shared secret was accepted or not
			require.Equal(t, []byte(nil), msg.OneTimePreKey)
			require.Equal(t, []byte(nil), msg.SignedPreKey)
			require.Equal(t, []byte(nil), msg.EphemeralKey)
			require.Equal(t, []byte(nil), msg.EphemeralKeySignature)

			// make sure that the receiver and sender are correct
			require.Equal(t, rawIdPubKeyBob, msg.Receiver)
			require.Equal(t, rawIdPubAliceBob, msg.Sender)
			require.Equal(t, "i am the message ID of the message to send", string(msg.MessageID))

			// make sure shared secret got attached
			require.Equal(t, hex.EncodeToString(sharedSecretBaseID), hex.EncodeToString(msg.UsedSharedSecret))

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
			drSession, err := dr.New([32]byte{1}, &drDhPair{
				x3dhPair: x3dh.KeyPair{
					PublicKey:  signedPreKeyBob.PublicKey,
					PrivateKey: signedPreKeyBob.PrivateKey,
				},
			})
			decryptedRawMessage, err := drSession.RatchetDecrypt(drMsg, nil)
			require.Nil(t, err)
			plainMsg := bpb.PlainChatMessage{}
			require.Nil(t, proto.Unmarshal(decryptedRawMessage, &plainMsg))

			// do some assertions
			require.Equal(t, plainMsgToSend.CreatedAt, plainMsg.CreatedAt)

			calledBackend = true
			return nil
		},
	}

	sharedSecretStore := testSharedSecretStorage{
		hasAny: func(key ed25519.PublicKey) (bool, error) {
			return true, nil
		},
		getYoungest: func(key ed25519.PublicKey) (*db.SharedSecret, error) {
			ss := &db.SharedSecret{
				Accepted: true,
				ID:       sharedSecretBaseID,
			}
			ss.SetX3dhSecret(x3dh.SharedSecret{1})
			return ss, nil
		},
	}

	userStorage := testUserStorage{
		getSignedPreKey: func(public ed25519.PublicKey) (*preKey.PreKey, error) {
			// at this point we would return
			// the signed pre key of our chat partner
			return &signedPreKeyBob, nil
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

	// shared secret base id
	id := make([]byte, 32)
	id[6] = 0x42

	plainMsgToSend := db.Message{
		ID:         "i am the message ID of the message to send",
		Version:    1,
		Status:     300,
		Message:    []byte("my message"),
		CreatedAt:  2147483648,
		Sender:     make([]byte, 32),
		DatabaseID: 2147483648,
	}

	msgStorage := testMessageStorage{
		updateStatus: func(partner ed25519.PublicKey, msgID int64, newStatus db.Status) error {
			require.Equal(t, hex.EncodeToString(rawIdPubKeyBob), hex.EncodeToString(partner))
			require.Equal(t, int64(2147483648), msgID)
			require.Equal(t, db.StatusSent, newStatus)
			return nil
		},
	}

	calledBackend := false
	backend := testBackend{
		submitMessages: func(messages []*bpb.ChatMessage) error {

			msg := messages[0]
			require.NotNil(t, msg)

			arrToSlice := func(arr [32]byte) []byte {
				return arr[:]
			}

			// all of those checks relate to the if a
			// shared secret was accepted or not
			require.Equal(t, arrToSlice([32]byte{3, 2, 4, 3, 2, 1}), msg.OneTimePreKey)
			require.Equal(t, arrToSlice([32]byte{4, 5, 3, 2}), msg.SignedPreKey)
			require.Equal(t, arrToSlice([32]byte{4, 3, 4}), msg.EphemeralKey)
			require.Equal(t, []byte{1, 3, 0, 3, 5}, msg.EphemeralKeySignature)

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
			drSession, err := dr.New([32]byte{1}, &drDhPair{
				x3dhPair: x3dh.KeyPair{
					PublicKey:  signedPreKeyBob.PublicKey,
					PrivateKey: signedPreKeyBob.PrivateKey,
				},
			})
			decryptedRawMessage, err := drSession.RatchetDecrypt(drMsg, nil)
			require.Nil(t, err)
			plainMsg := bpb.PlainChatMessage{}
			require.Nil(t, proto.Unmarshal(decryptedRawMessage, &plainMsg))

			// shared secret must be added since our shared secret haven't
			// been accepted
			require.Equal(t, id, plainMsg.SharedSecretBaseID)

			// make sure the shared secret creation date is the one from
			// the shared secret
			require.Equal(t, int64(4), plainMsg.SharedSecretCreationDate)

			// make sure message id is the same as the one we have chosen
			require.Equal(t, plainMsgToSend.ID, plainMsg.MessageID)

			// make sure message is the name as the one we have chosen
			require.Equal(t, plainMsgToSend.Message, plainMsg.Message)

			calledBackend = true
			return nil
		},
	}

	sharedSecretStore := testSharedSecretStorage{
		hasAny: func(key ed25519.PublicKey) (bool, error) {
			return true, nil
		},
		getYoungest: func(key ed25519.PublicKey) (*db.SharedSecret, error) {
			ss := &db.SharedSecret{
				Accepted:              false,
				CreatedAt:             time.Unix(4, 0),
				EphemeralKey:          x3dh.PublicKey{4, 3, 4},
				EphemeralKeySignature: []byte{1, 3, 0, 3, 5},
				UsedSignedPreKey:      x3dh.PublicKey{4, 5, 3, 2},
				UsedOneTimePreKey:     &x3dh.PublicKey{3, 2, 4, 3, 2, 1},
				ID:                    id,
			}
			ss.SetX3dhSecret(x3dh.SharedSecret{1})
			return ss, nil
		},
	}

	userStorage := testUserStorage{
		getSignedPreKey: func(partner ed25519.PublicKey) (*preKey.PreKey, error) {
			// at this point we would return
			// the signed pre key of our chat partner
			return &signedPreKeyBob, nil
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

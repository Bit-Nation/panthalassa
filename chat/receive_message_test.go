package chat

/**
import (
	"testing"
	"encoding/hex"

	db "github.com/Bit-Nation/panthalassa/db"
	require "github.com/stretchr/testify/require"
	queue "github.com/Bit-Nation/panthalassa/queue"
	uiapi "github.com/Bit-Nation/panthalassa/uiapi"
	backend "github.com/Bit-Nation/panthalassa/backend"
	bpb "github.com/Bit-Nation/protobuffers"
	"time"
	"fmt"
)

type chatTestBackendTransport struct {
	sendChan chan *bpb.BackendMessage
	receiveChan chan *bpb.BackendMessage
}

func (t *chatTestBackendTransport) Send(msg *bpb.BackendMessage) error {
	t.sendChan <- msg
	return nil
}

// will return the next message from the transport
func (t *chatTestBackendTransport) NextMessage() (*bpb.BackendMessage, error) {
	return <-t.receiveChan, nil
}

// close the transport
func (t *chatTestBackendTransport) Close() error {
	close(t.sendChan)
	close(t.receiveChan)
	return nil
}

func createAliceAndBob() (alice *Chat, aliceTrans *chatTestBackendTransport, bob *Chat, bobTrans *chatTestBackendTransport, err error) {

	// creates a chat. Don't forget to set the backend.
	createChat := func() (*Chat, *chatTestBackendTransport, error) {

		storm := createStorm()
		km := createKeyManager()

		chatStorage := db.NewChatStorage(storm, []func(e db.MessagePersistedEvent){}, km)
		ssStorage := db.NewBoltSharedSecretStorage(storm, km)
		drKeyStorage := db.NewBoltDRKeyStorage(storm, km)
		signedPreKeyStorage := db.NewBoltSignedPreKeyStorage(storm, km)
		oneTimePreKeyStorage := db.NewBoltOneTimePreKeyStorage(storm, km)
		userStorage := db.NewBoltUserStorage(storm)
		q := queue.New(queue.NewStorage(storm), 5, 1)

		trans := chatTestBackendTransport{
			sendChan: make(chan *bpb.BackendMessage, 10),
			receiveChan: make(chan *bpb.BackendMessage, 10),
		}

		up := uiapi.New(nil)
		b, err := backend.NewBackend(&trans, km, signedPreKeyStorage)
		if err != nil {
			return nil, nil, err
		}

		c, err := NewChat(Config{
			ChatStorage: chatStorage,
			Backend: b,
			SharedSecretDB: ssStorage,
			KM: km,
			DRKeyStorage: drKeyStorage,
			SignedPreKeyStorage: signedPreKeyStorage,
			OneTimePreKeyStorage: oneTimePreKeyStorage,
			UserStorage: userStorage,
			UiApi: up,
			Queue: q,
		})
		return c, &trans, err
	}

	// create alice
	alice, aliceTrans, err = createChat()
	if err != nil {
		return nil,nil, nil, nil, err
	}

	// create bob
	bob, bobTrans, err = createChat()
	if err != nil {
		return nil,nil, nil, nil, err
	}

	return

}
*/

/**
func TestChatBetweenAliceAndBob(t *testing.T) {

	alice, aliceTrans, bob, _, err := createAliceAndBob()
	require.Nil(t, err)
	aliceIDKeyStr, err := alice.km.IdentityPublicKey()
	require.Nil(t, err)
	aliceIDKey, err := hex.DecodeString(aliceIDKeyStr)
	require.Nil(t, err)

	bobIDKeyStr, err := bob.km.IdentityPublicKey()
	require.Nil(t, err)
	bobIDKey, err := hex.DecodeString(bobIDKeyStr)
	require.Nil(t, err)

	err = alice.SendMessage(bobIDKey, db.Message{
		ID: "my message id",
		Message: []byte("hi"),
		Version: 1,
		Status: db.StatusPersisted,
		CreatedAt: time.Now().UnixNano(),
		Sender: aliceIDKey,
	})
	require.Nil(t, err)

	go func() {
		for {
			msg := <-aliceTrans.sendChan

		}
	}()

	select {

	}

}
*/

/**
func TestByteSliceTox3dhPub(t *testing.T) {

	// to short
	_, err := byteSliceTox3dhPub([]byte{})
	require.EqualError(t, err, "got invalid x3dh public key (must have 32 bytes length)")

	// valid
	pubKey := make([]byte, 32)
	pubKey[3] = 0x33

	pub, err := byteSliceTox3dhPub(pubKey)
	require.Nil(t, err)

	require.Equal(t, hex.EncodeToString(pubKey), hex.EncodeToString(pub[:]))

}

func TestDoNotHandleOwnMessages(t *testing.T) {

	km := createKeyManager()

	c := Chat{
		km: km,
	}

	myIDKey, err := km.IdentityPublicKey()
	require.Nil(t, err)
	myIDRawKey, err := hex.DecodeString(myIDKey)
	require.Nil(t, err)

	// the sender must be a valid ed25519 public key
	err = c.handleReceivedMessage(&bpb.ChatMessage{
		Sender:           myIDRawKey,
		UsedSharedSecret: make([]byte, 32),
	})

	require.EqualError(t, err, "in can't handle messages I created my self - this is non sense")

}

func TestSenderTooShort(t *testing.T) {

	km := createKeyManager()

	c := Chat{
		km: km,
	}

	// the sender must be a valid ed25519 public key
	err := c.handleReceivedMessage(&bpb.ChatMessage{
		Sender:           []byte("too short"),
		UsedSharedSecret: make([]byte, 32),
	})

	require.EqualError(t, err, "sender public key too short")

}

func TestDoubleRatchetPKTooShort(t *testing.T) {

	km := createKeyManager()

	c := Chat{
		km: km,
	}

	// the double ratchet key must be 32 bytes long
	err := c.handleReceivedMessage(&bpb.ChatMessage{
		Sender:           make([]byte, 32),
		Message:          &bpb.DoubleRatchetMsg{},
		UsedSharedSecret: make([]byte, 32),
	})

	require.EqualError(t, err, "got invalid double ratchet public key - must have a length of 32")

}

func TestChatInitEphemeralKeySignatureValidation(t *testing.T) {

	sender, _, err := ed25519.GenerateKey(rand.Reader)
	require.Nil(t, err)

	km := createKeyManager()
	stormDB := createStorm()
	chatStorage := db.NewChatStorage(stormDB, nil, km)

	c := Chat{
		chatStorage: chatStorage,
		signedPreKeyStorage: &testSignedPreKeyStore{
			get: func(publicKey x3dh.PublicKey) (*x3dh.PrivateKey, error) {
				return &x3dh.PrivateKey{}, nil
			},
		},
		km: km,
	}

	err = c.handleReceivedMessage(&bpb.ChatMessage{
		Sender: sender,
		Message: &bpb.DoubleRatchetMsg{
			DoubleRatchetPK: make([]byte, 32),
		},
		SignedPreKey: make([]byte, 32),
		// definitely an invalid EphemeralKey
		EphemeralKey: make([]byte, 32),
		// definitely an invalid signature
		EphemeralKeySignature: []byte("invalid signature"),
		UsedSharedSecret:      make([]byte, 32),
	})

	require.EqualError(t, err, "aborted chat initialization - invalid ephemeral key")

}

func TestChatInitChatIDKeySignatureValidation(t *testing.T) {

	sender, senderPriv, err := ed25519.GenerateKey(rand.Reader)
	require.Nil(t, err)

	km := createKeyManager()
	stormDB := createStorm()
	chatStorage := db.NewChatStorage(stormDB, nil, km)

	c := Chat{
		chatStorage: chatStorage,
		signedPreKeyStorage: &testSignedPreKeyStore{
			get: func(publicKey x3dh.PublicKey) (*x3dh.PrivateKey, error) {
				return &x3dh.PrivateKey{}, nil
			},
		},
		km: createKeyManager(),
	}

	// mock ephemeralKey
	ephemeralKey := make([]byte, 32)
	ephemeralKey[3] = 0x32

	err = c.handleReceivedMessage(&bpb.ChatMessage{
		Sender: sender,
		Message: &bpb.DoubleRatchetMsg{
			DoubleRatchetPK: make([]byte, 32),
		},
		SignedPreKey: make([]byte, 32),
		// a ephemeralKey to pass the check for
		// if this message is a valid chat init message
		EphemeralKey: ephemeralKey,
		// a signature of the ephemeralKey to make sure the check pass
		EphemeralKeySignature: ed25519.Sign(senderPriv, ephemeralKey),
		// definitely an invalid EphemeralKey
		SenderChatIDKey: make([]byte, 32),
		// definitely an invalid signature
		SenderChatIDKeySignature: []byte("invalid signature"),
		UsedSharedSecret:         make([]byte, 32),
	})

	require.EqualError(t, err, "aborted chat initialization - invalid chat id key")

}

// test case for handling out of order messages
// in the case we already derived a shared secret
// from the same init params of a prev message
func TestChatInitDecryptMessageWhenSecretExist(t *testing.T) {

	curve25519 := x3dh.NewCurve25519(rand.Reader)
	sharedSecret := [32]byte{3, 4, 1, 3, 4, 9, 22}

	// alice would like to send a message to bob

	// create sender id key pair
	senderPub, senderPriv, err := ed25519.GenerateKey(rand.Reader)
	require.Nil(t, err)

	// create bob signed pre key
	bobSignedPreKey, err := curve25519.GenerateKeyPair()
	require.Nil(t, err)

	// create double ratchet session with alice
	var bobDRKey dr.Key
	copy(bobDRKey[:], bobSignedPreKey.PublicKey[:])

	aliceSession, err := dr.NewWithRemoteKey(sharedSecret, bobDRKey, dr.WithKeysStorage(&dr.KeysStorageInMemory{}))
	require.Nil(t, err)

	// marshal message
	rawPlainMsg, err := proto.Marshal(&bpb.PlainChatMessage{Message: []byte("hi bob"), MessageID: "my-message-id", Version: 1})
	require.Nil(t, err)

	// double ratchet message for bob
	drMessage := aliceSession.RatchetEncrypt(rawPlainMsg, nil)

	// bob key manager
	bobKm := createKeyManager()

	// bob id pub key
	//bobIDPubKeyStr, err := bobKm.IdentityPublicKey()
	require.Nil(t, err)
	//bobIDPubKey, err := hex.DecodeString(bobIDPubKeyStr)
	require.Nil(t, err)

	msg := &bpb.ChatMessage{
		Message: &bpb.DoubleRatchetMsg{
			DoubleRatchetPK: drMessage.Header.DH[:],
			N:               drMessage.Header.N,
			Pn:              drMessage.Header.PN,
			CipherText:      drMessage.Ciphertext,
		},
		SenderChatIDKey:          make([]byte, 32),
		SenderChatIDKeySignature: ed25519.Sign(senderPriv, make([]byte, 32)),
		EphemeralKey:             make([]byte, 32),
		EphemeralKeySignature:    ed25519.Sign(senderPriv, make([]byte, 32)),
		SignedPreKey:             bobSignedPreKey.PublicKey[:],
		Sender:                   senderPub,
		UsedSharedSecret:         make([]byte, 32),
	}

	km := createKeyManager()
	stormDB := createStorm()
	chatStorage := db.NewChatStorage(stormDB, nil, km)

	c := Chat{
		chatStorage: chatStorage,
		signedPreKeyStorage: &testSignedPreKeyStore{
			get: func(publicKey x3dh.PublicKey) (*x3dh.PrivateKey, error) {
				require.Equal(t, bobSignedPreKey.PublicKey, publicKey)
				return &bobSignedPreKey.PrivateKey, nil
			},
		},
		sharedSecStorage: &testSharedSecretStorage{
			get: func(key ed25519.PublicKey, sharedSecretID []byte) (*db.SharedSecret, error) {
				ss := &db.SharedSecret{}
				ss.SetX3dhSecret(sharedSecret)
				return ss, nil
			},
		},
		km: bobKm,
	}

	err = c.handleReceivedMessage(msg)
	require.Nil(t, err)

}

func TestChatInitSharedSecretAgreementAndMsgPersistence(t *testing.T) {

	curve25519 := x3dh.NewCurve25519(rand.Reader)

	bobSignedPreKeyPair, err := curve25519.GenerateKeyPair()
	require.Nil(t, err)

	bobChatIDKeyPair, err := curve25519.GenerateKeyPair()
	require.Nil(t, err)

	aliceIDKeyPair, err := curve25519.GenerateKeyPair()
	require.Nil(t, err)

	aliceX3dh := x3dh.New(&curve25519, sha256.New(), "proto", aliceIDKeyPair)
	aliceInitializedProto, err := aliceX3dh.CalculateSecret(testPreKeyBundle{
		identityKey:  bobChatIDKeyPair.PublicKey,
		signedPreKey: bobSignedPreKeyPair.PublicKey,
		validSignature: func() (bool, error) {
			return true, nil
		},
	})
	require.Nil(t, err)

	// create sender id key pair
	senderPub, senderPriv, err := ed25519.GenerateKey(rand.Reader)
	require.Nil(t, err)

	// create double ratchet session with alice
	var bobDRKey dr.Key
	copy(bobDRKey[:], bobSignedPreKeyPair.PublicKey[:])

	var sharedSec dr.Key
	copy(sharedSec[:], aliceInitializedProto.SharedSecret[:])

	aliceSession, err := dr.NewWithRemoteKey(sharedSec, bobDRKey, dr.WithKeysStorage(&dr.KeysStorageInMemory{}))
	require.Nil(t, err)

	sharedSecretIDRef := make([]byte, 32)
	sharedSecretIDRef[3] = 0x34

	// marshal message
	rawPlainMsg, err := proto.Marshal(&bpb.PlainChatMessage{
		Message:                  []byte("hi bob"),
		SharedSecretBaseID:       sharedSecretIDRef,
		SharedSecretCreationDate: 4,
		CreatedAt:                444,
		MessageID:                "message-id",
		Version:                  1,
	})
	require.Nil(t, err)

	// double ratchet message for bob
	drMessage := aliceSession.RatchetEncrypt(rawPlainMsg, nil)

	msg := &bpb.ChatMessage{
		Message: &bpb.DoubleRatchetMsg{
			DoubleRatchetPK: drMessage.Header.DH[:],
			N:               drMessage.Header.N,
			Pn:              drMessage.Header.PN,
			CipherText:      drMessage.Ciphertext,
		},
		SenderChatIDKey:          aliceIDKeyPair.PublicKey[:],
		SenderChatIDKeySignature: ed25519.Sign(senderPriv, aliceIDKeyPair.PublicKey[:]),
		EphemeralKey:             aliceInitializedProto.EphemeralKey[:],
		EphemeralKeySignature:    ed25519.Sign(senderPriv, aliceInitializedProto.EphemeralKey[:]),
		SignedPreKey:             bobSignedPreKeyPair.PublicKey[:],
		Sender:                   senderPub,
		UsedSharedSecret:         make([]byte, 32),
	}

	bobX3dh := x3dh.New(&curve25519, sha256.New(), "proto", bobChatIDKeyPair)

	km := createKeyManager()

	bobIDPubRaw, err := km.IdentityPublicKey()
	require.Nil(t, err)

	bobIDPub, err := hex.DecodeString(bobIDPubRaw)
	require.Nil(t, err)

	// chat storage
	stormDB := createStorm()
	chatStorage := db.NewChatStorage(stormDB, nil, km)

	sharedSecChan := make(chan *db.SharedSecret, 1)
	c := Chat{
		chatStorage: chatStorage,
		km:          km,
		signedPreKeyStorage: &testSignedPreKeyStore{
			get: func(publicKey x3dh.PublicKey) (*x3dh.PrivateKey, error) {
				// the public key must be bob's signed pre key
				// the matching private key must be returned
				require.Equal(t, bobSignedPreKeyPair.PublicKey, publicKey)
				return &bobSignedPreKeyPair.PrivateKey, nil
			},
		},
		sharedSecStorage: &testSharedSecretStorage{
			put: func(sharedSecret db.SharedSecret) error {
				// must be true since we received a chat init message
				require.True(t, sharedSecret.Accepted)

				// shared secret id
				sharedSecretID, err := sharedSecretID(senderPub, bobIDPub, sharedSecretIDRef)
				require.Nil(t, err)
				require.Equal(t, hex.EncodeToString(sharedSecretID), hex.EncodeToString(sharedSecret.ID))

				ss := sharedSecret.GetX3dhSecret()
				require.Equal(t, hex.EncodeToString(sharedSec[:]), hex.EncodeToString(ss[:]))
				require.Equal(t, x3dh.PublicKey{}, sharedSecret.EphemeralKey)
				require.Equal(t, []byte(nil), sharedSecret.EphemeralKeySignature)
				require.Nil(t, sharedSecret.UsedOneTimePreKey)
				require.Nil(t, sharedSecret.DestroyAt)
				sharedSecChan <- &sharedSecret
				return nil
			},
			get: func(key ed25519.PublicKey, sharedSecretID []byte) (*db.SharedSecret, error) {
				return nil, nil
			},
		},
		oneTimePreKeyStorage: &testOneTimePreKeyStorage{
			cut: func(pubKey []byte) (*x3dh.PrivateKey, error) {
				return nil, nil
			},
		},
		x3dh:         &bobX3dh,
		drKeyStorage: &dr.KeysStorageInMemory{},
	}

	err = c.handleReceivedMessage(msg)
	require.Nil(t, err)

}

func TestChatHandleInvalidShortSharedSecretID(t *testing.T) {

	km := createKeyManager()
	stormDB := createStorm()
	chatStorage := db.NewChatStorage(stormDB, nil, km)

	c := Chat{
		chatStorage: chatStorage,
		km:          km,
	}

	pub, _, err := ed25519.GenerateKey(rand.Reader)
	require.Nil(t, err)

	err = c.handleReceivedMessage(&bpb.ChatMessage{
		Sender: pub,
		Message: &bpb.DoubleRatchetMsg{
			DoubleRatchetPK: make([]byte, 32),
		},
	})

	require.EqualError(t, err, "message is not a chat initialisation message but don't contain information about which shared secret has been used")

}

func TestChatHandleNoSharedSecret(t *testing.T) {

	pub, _, err := ed25519.GenerateKey(rand.Reader)
	require.Nil(t, err)

	usedSharedSecret := make([]byte, 32)
	usedSharedSecret[3] = 0x33

	km := createKeyManager()
	stormDB := createStorm()
	chatStorage := db.NewChatStorage(stormDB, nil, km)

	c := Chat{
		chatStorage: chatStorage,
		sharedSecStorage: &testSharedSecretStorage{
			get: func(key ed25519.PublicKey, sharedSecretID []byte) (*db.SharedSecret, error) {
				require.Equal(t, pub, key)
				require.Equal(t, usedSharedSecret, sharedSecretID)
				return nil, nil
			},
		},
		km: km,
	}

	err = c.handleReceivedMessage(&bpb.ChatMessage{
		Sender:           pub,
		UsedSharedSecret: usedSharedSecret,
		Message: &bpb.DoubleRatchetMsg{
			DoubleRatchetPK: make([]byte, 32),
		},
	})

}

func TestChatHandleDecryptSuccessfullyAndAcceptSharedSecret(t *testing.T) {

	curve25519 := x3dh.NewCurve25519(rand.Reader)

	bobSignedPreKeyPair, err := curve25519.GenerateKeyPair()
	require.Nil(t, err)

	bobChatIDKeyPair, err := curve25519.GenerateKeyPair()
	require.Nil(t, err)

	aliceIDKeyPair, err := curve25519.GenerateKeyPair()
	require.Nil(t, err)

	aliceX3dh := x3dh.New(&curve25519, sha256.New(), "proto", aliceIDKeyPair)
	aliceInitializedProto, err := aliceX3dh.CalculateSecret(testPreKeyBundle{
		identityKey:  bobChatIDKeyPair.PublicKey,
		signedPreKey: bobSignedPreKeyPair.PublicKey,
		validSignature: func() (bool, error) {
			return true, nil
		},
	})
	require.Nil(t, err)

	// create sender id key pair
	senderPub, _, err := ed25519.GenerateKey(rand.Reader)
	require.Nil(t, err)

	// create double ratchet session with alice
	var bobDRKey dr.Key
	copy(bobDRKey[:], bobSignedPreKeyPair.PublicKey[:])

	var sharedSec dr.Key
	copy(sharedSec[:], aliceInitializedProto.SharedSecret[:])

	aliceSession, err := dr.NewWithRemoteKey(sharedSec, bobDRKey, dr.WithKeysStorage(&dr.KeysStorageInMemory{}))
	require.Nil(t, err)

	sharedSecretIDRef := make([]byte, 32)
	sharedSecretIDRef[3] = 0x34

	// marshal message
	rawPlainMsg, err := proto.Marshal(&bpb.PlainChatMessage{
		Message:            []byte("hi bob"),
		SharedSecretBaseID: sharedSecretIDRef,
		MessageID:          "message-id",
		Version:            1,
	})
	require.Nil(t, err)

	// double ratchet message for bob
	drMessage := aliceSession.RatchetEncrypt(rawPlainMsg, nil)

	bobX3dh := x3dh.New(&curve25519, sha256.New(), "proto", bobChatIDKeyPair)

	km := createKeyManager()

	bobIDPubRaw, err := km.IdentityPublicKey()
	require.Nil(t, err)

	bobIDPub, err := hex.DecodeString(bobIDPubRaw)
	require.Nil(t, err)

	expectedSharedSecretID, err := sharedSecretID(senderPub, bobIDPub, sharedSecretIDRef)

	msg := &bpb.ChatMessage{
		Message: &bpb.DoubleRatchetMsg{
			DoubleRatchetPK: drMessage.Header.DH[:],
			N:               drMessage.Header.N,
			Pn:              drMessage.Header.PN,
			CipherText:      drMessage.Ciphertext,
		},
		Sender:           senderPub,
		UsedSharedSecret: expectedSharedSecretID,
	}

	// chat storage
	stormDB := createStorm()
	chatStorage := db.NewChatStorage(stormDB, nil, km)

	c := Chat{
		chatStorage: chatStorage,
		km:          km,
		signedPreKeyStorage: &testSignedPreKeyStore{
			get: func(publicKey x3dh.PublicKey) (*x3dh.PrivateKey, error) {
				// the public key must be bob's signed pre key
				// the matching private key must be returned
				require.Equal(t, bobSignedPreKeyPair.PublicKey, publicKey)
				return &bobSignedPreKeyPair.PrivateKey, nil
			},
			all: func() ([]*x3dh.KeyPair, error) {
				return []*x3dh.KeyPair{
					&bobSignedPreKeyPair,
				}, nil
			},
		},
		sharedSecStorage: &testSharedSecretStorage{
			get: func(chatPartner ed25519.PublicKey, ssID []byte) (*db.SharedSecret, error) {
				// check if chat partner is really bob
				require.Equal(t, senderPub, chatPartner)
				// check id
				require.Nil(t, err)
				require.Equal(t, hex.EncodeToString(expectedSharedSecretID), hex.EncodeToString(ssID))
				// we return the secret that we share with alice
				ss := &db.SharedSecret{}
				ss.SetX3dhSecret(aliceInitializedProto.SharedSecret)
				return ss, nil
			},
			accept: func(ss db.SharedSecret) error {
				require.Equal(t, aliceInitializedProto.SharedSecret, ss.GetX3dhSecret())
				return nil
			},
		},
		x3dh:         &bobX3dh,
		drKeyStorage: &dr.KeysStorageInMemory{},
	}

	err = c.handleReceivedMessage(msg)
	require.Nil(t, err)

}
*/

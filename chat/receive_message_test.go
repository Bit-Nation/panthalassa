package chat

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"testing"

	db "github.com/Bit-Nation/panthalassa/db"
	bpb "github.com/Bit-Nation/protobuffers"
	x3dh "github.com/Bit-Nation/x3dh"
	proto "github.com/gogo/protobuf/proto"
	require "github.com/stretchr/testify/require"
	dr "github.com/tiabc/doubleratchet"
	ed25519 "golang.org/x/crypto/ed25519"
)

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

func TestSenderTooShort(t *testing.T) {

	c := Chat{}

	// the sender must be a valid ed25519 public key
	err := c.handleReceivedMessage(&bpb.ChatMessage{
		Sender:           []byte("too short"),
		UsedSharedSecret: make([]byte, 32),
	})

	require.EqualError(t, err, "sender public key too short")

}

func TestDoubleRatchetPKTooShort(t *testing.T) {

	c := Chat{}

	// the double ratchet key must be 32 bytes long
	err := c.handleReceivedMessage(&bpb.ChatMessage{
		Sender:           make([]byte, 32),
		Message:          &bpb.DoubleRatchedMsg{},
		UsedSharedSecret: make([]byte, 32),
	})

	require.EqualError(t, err, "got invalid double ratchet public key - must have a length of 32")

}

func TestChatInitEphemeralKeySignatureValidation(t *testing.T) {

	sender, _, err := ed25519.GenerateKey(rand.Reader)
	require.Nil(t, err)

	c := Chat{
		signedPreKeyStorage: &testSignedPreKeyStore{
			get: func(publicKey x3dh.PublicKey) (*x3dh.PrivateKey, error) {
				return &x3dh.PrivateKey{}, nil
			},
		},
	}

	err = c.handleReceivedMessage(&bpb.ChatMessage{
		Sender: sender,
		Message: &bpb.DoubleRatchedMsg{
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

	c := Chat{
		signedPreKeyStorage: &testSignedPreKeyStore{
			get: func(publicKey x3dh.PublicKey) (*x3dh.PrivateKey, error) {
				return &x3dh.PrivateKey{}, nil
			},
		},
	}

	// mock ephemeralKey
	ephemeralKey := make([]byte, 32)
	ephemeralKey[3] = 0x32

	err = c.handleReceivedMessage(&bpb.ChatMessage{
		Sender: sender,
		Message: &bpb.DoubleRatchedMsg{
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
	rawPlainMsg, err := proto.Marshal(&bpb.PlainChatMessage{Message: []byte("hi bob")})
	require.Nil(t, err)

	// double ratchet message for bob
	drMessage := aliceSession.RatchetEncrypt(rawPlainMsg, nil)

	msg := &bpb.ChatMessage{
		Message: &bpb.DoubleRatchedMsg{
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

	c := Chat{
		signedPreKeyStorage: &testSignedPreKeyStore{
			get: func(publicKey x3dh.PublicKey) (*x3dh.PrivateKey, error) {
				require.Equal(t, bobSignedPreKey.PublicKey, publicKey)
				return &bobSignedPreKey.PrivateKey, nil
			},
		},
		sharedSecStorage: &testSharedSecretStorage{
			secretForChatInitMsg: func(givenMsg *bpb.ChatMessage) (*db.SharedSecret, error) {
				require.Equal(t, msg, givenMsg)
				return &db.SharedSecret{
					X3dhSS: sharedSecret,
				}, nil
			},
		},
		messageDB: &testMessageStorage{
			persistReceivedMessage: func(partner ed25519.PublicKey, msg bpb.PlainChatMessage) error {
				require.Equal(t, senderPub, partner)
				require.Equal(t, "hi bob", string(msg.Message))
				return nil
			},
		},
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
	})
	require.Nil(t, err)

	// double ratchet message for bob
	drMessage := aliceSession.RatchetEncrypt(rawPlainMsg, nil)

	msg := &bpb.ChatMessage{
		Message: &bpb.DoubleRatchedMsg{
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

	c := Chat{
		km: km,
		signedPreKeyStorage: &testSignedPreKeyStore{
			get: func(publicKey x3dh.PublicKey) (*x3dh.PrivateKey, error) {
				// the public key must be bob's signed pre key
				// the matching private key must be returned
				require.Equal(t, bobSignedPreKeyPair.PublicKey, publicKey)
				return &bobSignedPreKeyPair.PrivateKey, nil
			},
		},
		sharedSecStorage: &testSharedSecretStorage{
			secretForChatInitMsg: func(givenMsg *bpb.ChatMessage) (*db.SharedSecret, error) {
				require.Equal(t, msg, givenMsg)
				return nil, nil
			},
			put: func(key ed25519.PublicKey, sharedSecret db.SharedSecret) error {
				// must be true since we received a chat init message
				require.True(t, sharedSecret.Accepted)

				// shared secret id
				sharedSecretID, err := sharedSecretID(senderPub, bobIDPub, sharedSecretIDRef)
				require.Nil(t, err)
				require.Equal(t, hex.EncodeToString(sharedSecretID), hex.EncodeToString(sharedSecret.ID))

				// shared secret id init params
				sharedSecretIDInitParams, err := sharedSecretInitID(msg.Sender, bobIDPub, *msg)
				require.Nil(t, err)
				require.Equal(t, hex.EncodeToString(sharedSecretIDInitParams), hex.EncodeToString(sharedSecret.IDInitParams))

				require.Equal(t, hex.EncodeToString(sharedSec[:]), hex.EncodeToString(sharedSecret.X3dhSS[:]))
				require.Equal(t, x3dh.PublicKey{}, sharedSecret.EphemeralKey)
				require.Equal(t, []byte(nil), sharedSecret.EphemeralKeySignature)
				require.Nil(t, sharedSecret.UsedOneTimePreKey)
				require.Nil(t, sharedSecret.DestroyAt)
				return nil
			},
		},
		messageDB: &testMessageStorage{
			persistReceivedMessage: func(partner ed25519.PublicKey, msg bpb.PlainChatMessage) error {
				require.Equal(t, senderPub, partner)
				require.Equal(t, "hi bob", string(msg.Message))
				require.Equal(t, sharedSecretIDRef, msg.SharedSecretBaseID)
				return nil
			},
		},
		oneTimePreKeyStorage: &testOneTimePreKeyStorage{
			cut: func(pubKey ed25519.PublicKey) (*x3dh.PrivateKey, error) {
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

	c := Chat{}

	pub, _, err := ed25519.GenerateKey(rand.Reader)
	require.Nil(t, err)

	err = c.handleReceivedMessage(&bpb.ChatMessage{
		Sender: pub,
		Message: &bpb.DoubleRatchedMsg{
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

	c := Chat{
		sharedSecStorage: &testSharedSecretStorage{
			get: func(key ed25519.PublicKey, sharedSecretID []byte) (*db.SharedSecret, error) {
				require.Equal(t, pub, key)
				require.Equal(t, usedSharedSecret, sharedSecretID)
				return nil, nil
			},
		},
	}

	err = c.handleReceivedMessage(&bpb.ChatMessage{
		Sender:           pub,
		UsedSharedSecret: usedSharedSecret,
		Message: &bpb.DoubleRatchedMsg{
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
		Message: &bpb.DoubleRatchedMsg{
			DoubleRatchetPK: drMessage.Header.DH[:],
			N:               drMessage.Header.N,
			Pn:              drMessage.Header.PN,
			CipherText:      drMessage.Ciphertext,
		},
		Sender:           senderPub,
		UsedSharedSecret: expectedSharedSecretID,
	}

	c := Chat{
		km: km,
		signedPreKeyStorage: &testSignedPreKeyStore{
			get: func(publicKey x3dh.PublicKey) (*x3dh.PrivateKey, error) {
				// the public key must be bob's signed pre key
				// the matching private key must be returned
				require.Equal(t, bobSignedPreKeyPair.PublicKey, publicKey)
				return &bobSignedPreKeyPair.PrivateKey, nil
			},
			all: func() []*x3dh.KeyPair {
				return []*x3dh.KeyPair{
					&bobSignedPreKeyPair,
				}
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
				return &db.SharedSecret{
					X3dhSS: aliceInitializedProto.SharedSecret,
				}, nil
			},
			accept: func(sharedSec *db.SharedSecret) error {
				require.Equal(t, aliceInitializedProto.SharedSecret, sharedSec.X3dhSS)
				return nil
			},
		},
		messageDB: &testMessageStorage{
			persistReceivedMessage: func(partner ed25519.PublicKey, msg bpb.PlainChatMessage) error {
				require.Equal(t, senderPub, partner)
				require.Equal(t, "hi bob", string(msg.Message))
				require.Equal(t, sharedSecretIDRef, msg.SharedSecretBaseID)
				return nil
			},
		},
		x3dh:         &bobX3dh,
		drKeyStorage: &dr.KeysStorageInMemory{},
	}

	err = c.handleReceivedMessage(msg)
	require.Nil(t, err)

}

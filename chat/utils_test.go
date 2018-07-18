package chat

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"testing"

	preKey "github.com/Bit-Nation/panthalassa/chat/prekey"
	bpb "github.com/Bit-Nation/protobuffers"
	x3dh "github.com/Bit-Nation/x3dh"
	mh "github.com/multiformats/go-multihash"
	require "github.com/stretchr/testify/require"
	ed25519 "golang.org/x/crypto/ed25519"
)

func TestChatRefreshSignedPreKeyBackendError(t *testing.T) {

	backend := testBackend{
		fetchSignedPreKey: func(userIdPubKey ed25519.PublicKey) (preKey.PreKey, error) {
			require.Equal(t, ed25519.PublicKey{1, 2}, userIdPubKey)
			return preKey.PreKey{}, errors.New("i am a test error")
		},
	}

	c := Chat{
		backend: &backend,
	}

	err := c.refreshSignedPreKey(ed25519.PublicKey{1, 2})
	require.EqualError(t, err, "i am a test error")

}

func TestChatRefreshSignedPreKeyInvalidSignature(t *testing.T) {

	km := createKeyManager()

	// create key paird
	c25519 := x3dh.NewCurve25519(rand.Reader)
	keyPair, err := c25519.GenerateKeyPair()
	require.Nil(t, err)

	// create pre key
	k := preKey.PreKey{}
	k.PublicKey = keyPair.PublicKey
	k.PrivateKey = keyPair.PrivateKey
	require.Nil(t, k.Sign(*km))

	// fake backend
	backend := testBackend{
		fetchSignedPreKey: func(userIdPubKey ed25519.PublicKey) (preKey.PreKey, error) {
			return k, nil
		},
	}

	// chat
	c := Chat{
		backend: &backend,
	}

	wrongIDKey := [32]byte{1, 2}

	// refresh signed pre key
	// should fail since we return a wrong
	// signed pre key
	err = c.refreshSignedPreKey(wrongIDKey[:])
	require.EqualError(t, err, "signed pre key signature is invalid")

}

func TestChatRefreshSignedSuccess(t *testing.T) {

	km := createKeyManager()

	idPubKeyStr, err := km.IdentityPublicKey()
	require.Nil(t, err)

	idPubKey, err := hex.DecodeString(idPubKeyStr)
	require.Nil(t, err)

	// create key paird
	c25519 := x3dh.NewCurve25519(rand.Reader)
	keyPair, err := c25519.GenerateKeyPair()
	require.Nil(t, err)

	// create pre key
	k := preKey.PreKey{}
	k.PublicKey = keyPair.PublicKey
	k.PrivateKey = keyPair.PrivateKey
	require.Nil(t, k.Sign(*km))

	// fake backend
	backend := testBackend{
		fetchSignedPreKey: func(userIdPubKey ed25519.PublicKey) (preKey.PreKey, error) {
			return k, nil
		},
	}

	calledStore := false
	userStorage := testUserStorage{
		putSignedPreKey: func(idKey ed25519.PublicKey, key preKey.PreKey) error {
			require.Equal(t, hex.EncodeToString(idPubKey), hex.EncodeToString(idKey))
			require.Equal(t, k, key)
			calledStore = true
			return nil
		},
	}

	// chat
	c := Chat{
		backend:     &backend,
		userStorage: &userStorage,
	}

	err = c.refreshSignedPreKey(idPubKey)
	require.Nil(t, err)
	require.True(t, calledStore)

}

func TestHashChatMessage(t *testing.T) {

	msg := bpb.ChatMessage{
		OneTimePreKey:            []byte("one time pre key"),
		SignedPreKey:             []byte("signed pre key"),
		EphemeralKey:             []byte("ephemeral key"),
		EphemeralKeySignature:    []byte("ephemeral key signature"),
		SenderChatIDKey:          []byte("sender chat id key"),
		SenderChatIDKeySignature: []byte("sender chat id key signature"),
		SharedSecretCreationDate: 34564,
		Message: &bpb.DoubleRatchedMsg{
			DoubleRatchetPK: []byte("double ratchet pk"),
			N:               uint32(3),
			Pn:              uint32(4544),
			CipherText:      []byte("cipher text"),
		},
		Receiver:         []byte("receiver"),
		Sender:           []byte("sender"),
		MessageID:        []byte("message id"),
		UsedSharedSecret: []byte("used shared secret"),
	}

	// manual hash
	manHash := bytes.NewBuffer(nil)

	// first block of metadata
	_, err := manHash.Write([]byte("one time pre key"))
	require.Nil(t, err)
	_, err = manHash.Write([]byte("signed pre key"))
	require.Nil(t, err)
	_, err = manHash.Write([]byte("ephemeral key"))
	require.Nil(t, err)
	_, err = manHash.Write([]byte("ephemeral key signature"))
	require.Nil(t, err)
	_, err = manHash.Write([]byte("sender chat id key"))
	require.Nil(t, err)
	_, err = manHash.Write([]byte("sender chat id key signature"))
	require.Nil(t, err)

	// write time stamp
	ssCreationDate := make([]byte, 8)
	binary.BigEndian.PutUint64(ssCreationDate, uint64(34564))
	_, err = manHash.Write(ssCreationDate)
	require.Nil(t, err)

	// write message data
	_, err = manHash.Write([]byte("double ratchet pk"))
	require.Nil(t, err)

	n := make([]byte, 4)
	binary.BigEndian.PutUint32(n, uint32(3))
	_, err = manHash.Write(n)
	require.Nil(t, err)

	pn := make([]byte, 4)
	binary.BigEndian.PutUint32(pn, uint32(4544))
	_, err = manHash.Write(pn)
	require.Nil(t, err)

	_, err = manHash.Write([]byte("cipher text"))
	require.Nil(t, err)

	// second block of metadata
	_, err = manHash.Write([]byte("receiver"))
	require.Nil(t, err)
	_, err = manHash.Write([]byte("sender"))
	require.Nil(t, err)
	_, err = manHash.Write([]byte("message id"))
	require.Nil(t, err)
	_, err = manHash.Write([]byte("used shared secret"))
	require.Nil(t, err)

	mhManHash, err := mh.Sum(manHash.Bytes(), mh.SHA3_256, -1)
	require.Nil(t, err)

	msgHash, err := hashChatMessage(msg)
	require.Nil(t, err)

	require.Equal(t, mhManHash, msgHash)

}

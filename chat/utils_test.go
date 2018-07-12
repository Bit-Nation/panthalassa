package chat

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"testing"

	preKey "github.com/Bit-Nation/panthalassa/chat/prekey"
	x3dh "github.com/Bit-Nation/x3dh"
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

package chat

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"testing"

	keyManager "github.com/Bit-Nation/panthalassa/keyManager"
	keyStore "github.com/Bit-Nation/panthalassa/keyStore"
	mnemonic "github.com/Bit-Nation/panthalassa/mnemonic"
	x3dh "github.com/Bit-Nation/x3dh"
	require "github.com/stretchr/testify/require"
	dr "github.com/tiabc/doubleratchet"
)

func TestChat_HandleInitialMessage(t *testing.T) {

	// chat alice
	mneAlice, err := mnemonic.New()
	require.Nil(t, err)

	ksAlice, err := keyStore.NewFromMnemonic(mneAlice)
	require.Nil(t, err)

	kmAlice := keyManager.CreateFromKeyStore(ksAlice)
	IDKeyPairAlice, err := kmAlice.ChatIdKeyPair()
	require.Nil(t, err)

	curveAlice := x3dh.NewCurve25519(rand.Reader)

	xAlice := x3dh.New(&curveAlice, sha256.New(), "testing", IDKeyPairAlice)

	chatAlice := Chat{
		doubleRachetKeyStore: &dr.KeysStorageInMemory{},
		x3dh:                 xAlice,
		km:                   kmAlice,
	}

	// chat bob
	mneBob, err := mnemonic.New()
	require.Nil(t, err)

	ksBob, err := keyStore.NewFromMnemonic(mneBob)
	require.Nil(t, err)

	kmBob := keyManager.CreateFromKeyStore(ksBob)
	IDKeyPairBob, err := kmBob.ChatIdKeyPair()
	require.Nil(t, err)

	curveBob := x3dh.NewCurve25519(rand.Reader)

	xBob := x3dh.New(&curveBob, sha256.New(), "testing", IDKeyPairBob)

	chatBob := Chat{
		doubleRachetKeyStore: &dr.KeysStorageInMemory{},
		x3dh:                 xBob,
		km:                   kmBob,
	}

	bobPreKeyBundle, err := chatBob.NewPreKeyBundle()
	require.Nil(t, err)

	idPublicKeyStrAlice, err := chatAlice.km.IdentityPublicKey()
	require.Nil(t, err)

	idPublicKeyAlice, err := hex.DecodeString(idPublicKeyStrAlice)
	require.Nil(t, err)

	msgFromAlice, initializedProtocol, err := chatAlice.InitializeChat(idPublicKeyAlice, bobPreKeyBundle)

	sharedSecret, err := chatBob.HandleInitialMessage(msgFromAlice, PreKeyBundlePrivate{
		OneTimePreKey: bobPreKeyBundle.PrivatePart.OneTimePreKey,
		SignedPreKey:  bobPreKeyBundle.PrivatePart.SignedPreKey,
	})
	require.Nil(t, err)

	require.Equal(t, initializedProtocol.SharedSecret, sharedSecret)

	// decrypt the message
	plainMsg, err := chatBob.DecryptMessage(sharedSecret, msgFromAlice)
	require.Nil(t, err)
	require.Equal(t, "hi", plainMsg)

	plainMsg, err = chatBob.DecryptMessage(sharedSecret, msgFromAlice)
	require.Nil(t, err)
	require.Equal(t, "hi", plainMsg)

	plainMsg, err = chatBob.DecryptMessage(sharedSecret, msgFromAlice)
	require.Nil(t, err)
	require.Equal(t, "hi", plainMsg)
}

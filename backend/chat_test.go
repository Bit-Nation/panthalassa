package backend

import (
	"crypto/rand"
	"encoding/hex"
	"testing"

	preKey "github.com/Bit-Nation/panthalassa/chat/prekey"
	keyManager "github.com/Bit-Nation/panthalassa/keyManager"
	keyStore "github.com/Bit-Nation/panthalassa/keyStore"
	mnemonic "github.com/Bit-Nation/panthalassa/mnemonic"
	profile "github.com/Bit-Nation/panthalassa/profile"
	bpb "github.com/Bit-Nation/protobuffers"
	x3dh "github.com/Bit-Nation/x3dh"
	require "github.com/stretchr/testify/require"
	ed25519 "golang.org/x/crypto/ed25519"
)

func TestBackend_FetchSignedPreKey(t *testing.T) {

	// key manager setup
	mne, err := mnemonic.New()
	require.Nil(t, err)
	ks, err := keyStore.NewFromMnemonic(mne)
	require.Nil(t, err)
	km := keyManager.CreateFromKeyStore(ks)

	// target person signed pre key
	c25519 := x3dh.NewCurve25519(rand.Reader)
	signingKeyPair, err := c25519.GenerateKeyPair()
	require.Nil(t, err)

	// signed pre key
	signedPreKey := preKey.PreKey{}
	signedPreKey.PublicKey = signingKeyPair.PublicKey
	require.Nil(t, signedPreKey.Sign(*km))
	signedPreProto, err := signedPreKey.ToProtobuf()
	require.Nil(t, err)

	// identity key
	identityKey, err := km.IdentityPublicKey()
	require.Nil(t, err)
	rawIdentityKey, err := hex.DecodeString(identityKey)
	require.Nil(t, err)

	// transport
	transport := testTransport{}
	reqIDChan := make(chan string, 1)
	transport.send = func(msg *bpb.BackendMessage) error {
		require.Equal(t, identityKey, hex.EncodeToString(msg.Request.SignedPreKey))
		reqIDChan <- msg.RequestID
		// send response with signed protobuf back
		return nil
	}
	transport.nextMessage = func() (*bpb.BackendMessage, error) {
		return &bpb.BackendMessage{
			RequestID: <-reqIDChan,
			Response: &bpb.BackendMessage_Response{
				SignedPreKey: &signedPreProto,
			},
		}, nil
	}

	b, err := NewBackend(&transport, nil)
	b.authenticate <- true
	require.Nil(t, err)

	// fetched signed pre key
	fetchedSignedPreKey, err := b.FetchSignedPreKey(rawIdentityKey)
	require.Nil(t, err)

	valid, err := fetchedSignedPreKey.VerifySignature(rawIdentityKey)
	require.Nil(t, err)
	require.True(t, valid)

	// make sure we got what we expected
	require.Equal(t, hex.EncodeToString(signedPreKey.PublicKey[:]), hex.EncodeToString(fetchedSignedPreKey.PublicKey[:]))

}

func TestBackend_FetchSignedPreKeyInvalidSignature(t *testing.T) {

	// key manager setup
	mne, err := mnemonic.New()
	require.Nil(t, err)
	ks, err := keyStore.NewFromMnemonic(mne)
	require.Nil(t, err)
	km := keyManager.CreateFromKeyStore(ks)

	// target person signed pre key
	c25519 := x3dh.NewCurve25519(rand.Reader)
	signingKeyPair, err := c25519.GenerateKeyPair()
	require.Nil(t, err)

	// signed pre key
	signedPreKey := preKey.PreKey{}
	signedPreKey.PublicKey = signingKeyPair.PublicKey
	require.Nil(t, signedPreKey.Sign(*km))
	signedPreProto, err := signedPreKey.ToProtobuf()
	require.Nil(t, err)

	// identity key
	identityKey, err := km.IdentityPublicKey()
	require.Nil(t, err)

	// transport
	transport := testTransport{}
	reqIDChan := make(chan string, 1)
	transport.send = func(msg *bpb.BackendMessage) error {
		// should be not equal since we testing what happens if the
		// returned signed pre key is not the one the client ask for
		require.NotEqual(t, identityKey, hex.EncodeToString(msg.Request.SignedPreKey))
		reqIDChan <- msg.RequestID
		return nil
	}
	transport.nextMessage = func() (*bpb.BackendMessage, error) {
		return &bpb.BackendMessage{
			RequestID: <-reqIDChan,
			Response: &bpb.BackendMessage_Response{
				SignedPreKey: &signedPreProto,
			},
		}, nil
	}

	b, err := NewBackend(&transport, nil)
	b.authenticate <- true
	require.Nil(t, err)

	// the chat partner of which we would like to receive the signed pre key
	chatPartner, _, _ := ed25519.GenerateKey(rand.Reader)

	// fetched signed pre key
	// chatPartner Must not be the same as the signer of the signed pre key
	// since we want to test what happens if we receive an signed pre key
	// of the wrong person
	fetchedSignedPreKey, err := b.FetchSignedPreKey(chatPartner)
	require.EqualError(t, err, "invalid signed pre key signature")

	valid, err := fetchedSignedPreKey.VerifySignature(chatPartner)
	require.EqualError(t, err, "got invalid identity key public key")
	require.False(t, valid)

}

func TestBackend_FetchPreKeyBundle(t *testing.T) {

	// key manager setup
	mne, err := mnemonic.New()
	require.Nil(t, err)
	ks, err := keyStore.NewFromMnemonic(mne)
	require.Nil(t, err)
	km := keyManager.CreateFromKeyStore(ks)

	// target person signed pre key
	c25519 := x3dh.NewCurve25519(rand.Reader)
	signingKeyPair, err := c25519.GenerateKeyPair()
	require.Nil(t, err)

	// signed pre key
	signedPreKey := preKey.PreKey{}
	signedPreKey.PublicKey = signingKeyPair.PublicKey
	require.Nil(t, signedPreKey.Sign(*km))
	signedPreProto, err := signedPreKey.ToProtobuf()
	require.Nil(t, err)

	// identity key
	identityKey, err := km.IdentityPublicKey()
	require.Nil(t, err)
	rawIdentityKey, err := hex.DecodeString(identityKey)
	require.Nil(t, err)

	// profile
	prof, err := profile.SignProfile("Florian", "earth", "base64", *km)
	require.Nil(t, err)
	protoProf, err := prof.ToProtobuf()
	require.Nil(t, err)

	// transport
	transport := testTransport{}
	reqIDChan := make(chan string)
	transport.send = func(msg *bpb.BackendMessage) error {
		require.Equal(t, identityKey, hex.EncodeToString(msg.Request.PreKeyBundle))
		// send response with signed protobuf back
		reqIDChan <- msg.RequestID
		return nil
	}
	transport.nextMessage = func() (*bpb.BackendMessage, error) {
		return &bpb.BackendMessage{
			RequestID: <-reqIDChan,
			Response: &bpb.BackendMessage_Response{
				PreKeyBundle: &bpb.BackendMessage_PreKeyBundle{
					SignedPreKey: &signedPreProto,
					Profile:      protoProf,
				},
			},
		}, nil
	}

	b, err := NewBackend(&transport, nil)
	b.authenticate <- true
	require.Nil(t, err)

	fetchedSignedPreKey, err := b.FetchPreKeyBundle(rawIdentityKey)
	require.Nil(t, err)

	valid, err := fetchedSignedPreKey.ValidSignature()
	require.Nil(t, err)
	require.True(t, valid)

}

func TestBackend_SubmitMessage(t *testing.T) {

	// transport
	transport := testTransport{}
	reqIDChan := make(chan string, 1)
	transport.send = func(msg *bpb.BackendMessage) error {
		require.Equal(t, 2, len(msg.Request.Messages))
		reqIDChan <- msg.RequestID
		return nil
	}
	transport.nextMessage = func() (*bpb.BackendMessage, error) {
		return &bpb.BackendMessage{
			RequestID: <-reqIDChan,
			Response:  &bpb.BackendMessage_Response{},
		}, nil
	}

	b, err := NewBackend(&transport, nil)
	b.authenticate <- true
	require.Nil(t, err)

	err = b.SubmitMessages([]*bpb.ChatMessage{
		&bpb.ChatMessage{},
		&bpb.ChatMessage{},
	})
	require.Nil(t, err)

}

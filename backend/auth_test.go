package backend

import (
	"testing"

	keyManager "github.com/Bit-Nation/panthalassa/keyManager"
	keyStore "github.com/Bit-Nation/panthalassa/keyStore"
	mnemonic "github.com/Bit-Nation/panthalassa/mnemonic"
	bpb "github.com/Bit-Nation/protobuffers"
	require "github.com/stretchr/testify/require"
	ed25519 "golang.org/x/crypto/ed25519"
)

func TestBackend_Auth(t *testing.T) {

	// key manager setup
	mne, err := mnemonic.New()
	require.Nil(t, err)
	ks, err := keyStore.NewFromMnemonic(mne)
	require.Nil(t, err)
	km := keyManager.CreateFromKeyStore(ks)

	// transport
	transport := testTransport{}
	transport.send = func(msg *bpb.BackendMessage) error {

		auth := msg.Response.Auth

		valid := ed25519.Verify(auth.IdentityPublicKey, auth.ToSign, auth.Signature)
		require.True(t, valid)

		return nil
	}

	_, err = NewBackend(&transport, km)
	require.Nil(t, err)
	err = transport.onMessage(&bpb.BackendMessage{
		RequestID: "i_am_the_request_id",
		Request: &bpb.BackendMessage_Request{
			Auth: &bpb.BackendMessage_Auth{
				ToSign: []byte{1, 2, 3, 4},
			},
		},
	})
	require.Nil(t, err)

}

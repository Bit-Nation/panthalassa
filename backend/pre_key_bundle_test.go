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
)

func TestPreKeyBundleFromProto(t *testing.T) {

	// curve 25519
	c25519 := x3dh.NewCurve25519(rand.Reader)

	// mnemonic
	mne, err := mnemonic.New()
	require.Nil(t, err)

	// key store
	ks, err := keyStore.NewFromMnemonic(mne)
	require.Nil(t, err)

	// key manager
	km := keyManager.CreateFromKeyStore(ks)

	// identity key
	idPubKeyRaw, err := km.IdentityPublicKey()
	require.Nil(t, err)
	idPubKey, err := hex.DecodeString(idPubKeyRaw)
	require.Nil(t, err)

	// profile
	prof, err := profile.SignProfile("Florian", "Earth", "base64...", *km)
	require.Nil(t, err)

	// profile protobuf
	protoProf, err := prof.ToProtobuf()
	require.Nil(t, err)

	// signed pre key
	signedPreKeyPair, err := c25519.GenerateKeyPair()
	require.Nil(t, err)
	signedPreKey := preKey.PreKey{}
	signedPreKey.PublicKey = signedPreKeyPair.PublicKey
	require.Nil(t, signedPreKey.Sign(*km))

	// one time pre key
	oneTimePreKeyPair, err := c25519.GenerateKeyPair()
	require.Nil(t, err)
	oneTimePreKey := preKey.PreKey{}
	oneTimePreKey.PublicKey = oneTimePreKeyPair.PublicKey
	require.Nil(t, oneTimePreKey.Sign(*km))

	// proto pre key bundle
	protoPreKeyBundle := bpb.BackendMessage_PreKeyBundle{
		Profile: protoProf,
		SignedPreKey: func() *bpb.PreKey {
			k, err := signedPreKey.ToProtobuf()
			require.Nil(t, err)
			return &k
		}(),
		OneTimePreKey: func() *bpb.PreKey {
			k, err := oneTimePreKey.ToProtobuf()
			require.Nil(t, err)
			return &k
		}(),
	}

	// pre key bundle
	preKeyBundle, err := PreKeyBundleFromProto(idPubKey, &protoPreKeyBundle)
	require.Nil(t, err)

	// validate pre key bundle signatures
	valid, err := preKeyBundle.ValidSignature()
	require.Nil(t, err)
	require.True(t, valid)

	// make sure keys are correct
	require.Equal(t, prof.Information.ChatIDKey, preKeyBundle.IdentityKey())
	require.Equal(t, signedPreKeyPair.PublicKey, preKeyBundle.SignedPreKey())
	require.Equal(t, &oneTimePreKeyPair.PublicKey, preKeyBundle.OneTimePreKey())

}

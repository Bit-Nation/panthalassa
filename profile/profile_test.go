package profile

import (
	"testing"

	km "github.com/Bit-Nation/panthalassa/keyManager"
	ks "github.com/Bit-Nation/panthalassa/keyStore"
	mnemonic "github.com/Bit-Nation/panthalassa/mnemonic"
	require "github.com/stretchr/testify/require"
)

// @todo add test's which compute the expected signature as well and check it with the signature returned in the profile
// @todo add test's which call "SignaturesValid" with invalid data
// @todo add test's for unmarshal

func TestSignProfile(t *testing.T) {

	// create test mnemonic
	mne, err := mnemonic.New()
	require.Nil(t, err)

	// create key store
	store, err := ks.NewFromMnemonic(mne)
	require.Nil(t, err)

	// open key manger with created keystore
	keyManager := km.CreateFromKeyStore(store)

	// create profile
	prof, err := SignProfile("Florian", "Earth", "base64", *keyManager)
	require.Nil(t, err)

	// basic check's
	require.Equal(t, "Florian", prof.Information.Name)
	require.Equal(t, "Earth", prof.Information.Location)
	require.Equal(t, "base64", prof.Information.Image)

	// validate profile
	valid, err := prof.SignaturesValid()
	require.Nil(t, err)
	require.True(t, valid)

}

func TestSignProfileKeyStore(t *testing.T) {

	// create test mnemonic
	mne, err := mnemonic.New()
	require.Nil(t, err)

	// create key store from mnemonic
	keyStore, err := ks.NewFromMnemonic(mne)
	require.Nil(t, err)

	// create key manager
	keyManager := km.CreateFromKeyStore(keyStore)

	// export key manager store
	keyManagerStore, err := keyManager.Export("pw", "pw")
	require.Nil(t, err)

	profile, err := SignWithKeyManagerStore("Florian", "Earth", "base64", keyManagerStore, "pw")
	require.Nil(t, err)

	valid, err := profile.SignaturesValid()
	require.Nil(t, err)
	require.True(t, valid)

	require.Equal(t, "Florian", profile.Information.Name)
	require.Equal(t, "Earth", profile.Information.Location)
	require.Equal(t, "base64", profile.Information.Image)

}

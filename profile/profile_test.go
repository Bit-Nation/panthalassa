package profile

import (
	"encoding/hex"
	"testing"
	"time"

	km "github.com/Bit-Nation/panthalassa/keyManager"
	ks "github.com/Bit-Nation/panthalassa/keyStore"
	mnemonic "github.com/Bit-Nation/panthalassa/mnemonic"
	require "github.com/stretchr/testify/require"
)

// @todo add test's which compute the expected signature as well and check it with the signature returned in the profile
// @todo add test's which call "SignaturesValid" with invalid data
// @todo add test's for unmarshal

func TestProfile(t *testing.T) {

	// create test mnemonic
	mne, err := mnemonic.FromString("warrior come shuffle soccer dragon cube embody labor display junk metal left chef drive venue home maximum lounge brush scheme return liquid again chaos")
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

	// validate profile Signatures
	valid, err := prof.SignaturesValid()
	require.Nil(t, err)
	require.True(t, valid)

	// test protobuf transformation
	pp, err := prof.ToProtobuf()
	require.Nil(t, err)
	require.Equal(t, pp.EthereumPubKey, prof.Information.EthereumPubKey)
	require.Equal(t, pp.IdentityPubKey, prof.Information.IdentityPubKey)
	require.Equal(t, pp.Name, prof.Information.Name)
	require.Equal(t, pp.Image, prof.Information.Image)
	require.Equal(t, pp.Location, prof.Information.Location)
	require.Equal(t, hex.EncodeToString(pp.ChatIdentityPubKey), hex.EncodeToString(prof.Information.ChatIDKey[:]))
	require.Equal(t, pp.Timestamp, prof.Information.Timestamp.Format(TimeFormat))
	require.Equal(t, uint8(pp.Version), prof.Information.Version)
	require.Equal(t, pp.EthereumKeySignature, prof.Signatures.EthereumKey)
	require.Equal(t, pp.IdentityKeySignature, prof.Signatures.IdentityKey)

	// test protobuf to profile
	prof, err = ProtobufToProfile(pp)
	require.Nil(t, err)
	require.Equal(t, pp.EthereumPubKey, prof.Information.EthereumPubKey)
	require.Equal(t, pp.IdentityPubKey, prof.Information.IdentityPubKey)
	require.Equal(t, pp.Name, prof.Information.Name)
	require.Equal(t, pp.Image, prof.Information.Image)
	require.Equal(t, pp.Location, prof.Information.Location)
	require.Equal(t, hex.EncodeToString(pp.ChatIdentityPubKey), hex.EncodeToString(prof.Information.ChatIDKey[:]))
	require.Equal(t, pp.Timestamp, prof.Information.Timestamp.Format(TimeFormat))
	require.Equal(t, uint8(pp.Version), prof.Information.Version)
	require.Equal(t, pp.EthereumKeySignature, prof.Signatures.EthereumKey)
	require.Equal(t, pp.IdentityKeySignature, prof.Signatures.IdentityKey)

	// test IdentityPublicKey
	key, err := prof.IdentityPublicKey()
	require.Nil(t, err)
	require.Equal(t, hex.EncodeToString(key), hex.EncodeToString(pp.IdentityPubKey))

	// test EthereumPublicKey
	key, err = prof.EthereumPublicKey()
	require.Nil(t, err)
	require.Equal(t, hex.EncodeToString(key), hex.EncodeToString(pp.EthereumPubKey))

	addr, err := prof.EthereumAddress()
	require.Nil(t, err)
	require.Equal(t, "0xb6E578c7535863c075e2A1411eF0Fa0c6E12416A", addr.String())
}

func TestHash(t *testing.T) {

	timeStamp, err := time.Parse(TimeFormat, "Mon Jan  1 00:00:00 UTC 0001")
	require.Nil(t, err)

	// create profile
	prof := Profile{
		Information: Information{
			Name:           "Florian",
			Location:       "Earth",
			Image:          "base64",
			IdentityPubKey: []byte{3, 4, 5},
			EthereumPubKey: []byte{1, 3, 4},
			ChatIDKey:      [32]byte{1, 0, 3},
			Timestamp:      timeStamp,
			Version:        3,
		},
	}

	h, err := prof.Hash()
	require.Nil(t, err)
	require.Equal(t, "12207877aff128d20141ff49d90fc9ca2eb4c47536d52d5753c926f137aab4903ced", h.String())

}

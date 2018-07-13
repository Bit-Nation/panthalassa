package prekey

import (
	"encoding/hex"
	"testing"
	"time"
  
	keyManager "github.com/Bit-Nation/panthalassa/keyManager"
	keyStore "github.com/Bit-Nation/panthalassa/keyStore"
	mnemonic "github.com/Bit-Nation/panthalassa/mnemonic"
	require "github.com/stretchr/testify/require"
)

func TestHashInvalidIdPublicKey(t *testing.T) {
	k := PreKey{}
	_, err := k.hash()
	require.EqualError(t, err, "got invalid identity key public key")
}

func TestHashInvalidPublicKey(t *testing.T) {
	k := PreKey{}
	k.identityPublicKey = [32]byte{1}
	_, err := k.hash()
	require.EqualError(t, err, "got invalid pre key public key")
}

func TestSign(t *testing.T) {

	mne, err := mnemonic.New()
	require.Nil(t, err)

	ks, err := keyStore.NewFromMnemonic(mne)
	require.Nil(t, err)

	km := keyManager.CreateFromKeyStore(ks)

	k := PreKey{}
	k.PublicKey = [32]byte{1}

	require.Nil(t, k.Sign(*km))

	valid, err := k.VerifySignature(k.identityPublicKey[:])
	require.Nil(t, err)
	require.True(t, valid)

}

func TestPreKey_ToProtobufAndBack(t *testing.T) {

	now := time.Now()

	k := PreKey{
		identityPublicKey: [32]byte{1},
		signature:         []byte{2},
		time:              now,
	}
	k.PublicKey = [32]byte{3}

	pp, err := k.ToProtobuf()
	require.Nil(t, err)

	require.Equal(t, hex.EncodeToString(k.identityPublicKey[:]), hex.EncodeToString(pp.IdentityKey))
	require.Equal(t, k.signature, pp.IdentityKeySignature)
	require.Equal(t, k.time.Unix(), pp.TimeStamp)
	require.Equal(t, hex.EncodeToString(k.PublicKey[:]), hex.EncodeToString(pp.Key))

	k, err = FromProtoBuf(pp)
	require.Nil(t, err)

	require.Equal(t, hex.EncodeToString(pp.IdentityKey), hex.EncodeToString(k.identityPublicKey[:]))
	require.Equal(t, pp.IdentityKeySignature, k.signature)
	require.Equal(t, pp.TimeStamp, k.time.Unix())
	require.Equal(t, hex.EncodeToString(pp.Key), hex.EncodeToString(k.PublicKey[:]))

}

func TestPreKey_OlderThen(t *testing.T) {
	k := PreKey{
		time: time.Now().Truncate(time.Second * 10),
	}
	require.False(t, k.OlderThen(time.Second*5))
}

package chat

import (
	"testing"
	"time"
	"github.com/tiabc/doubleratchet"
	"github.com/Bit-Nation/panthalassa/mnemonic"
	"github.com/Bit-Nation/panthalassa/keyStore"
	"github.com/Bit-Nation/panthalassa/keyManager"
	"github.com/stretchr/testify/require"
	"encoding/hex"
)

func TestMessage_SignVerify(t *testing.T) {
	
	mne, err := mnemonic.New()
	require.Nil(t, err)
	
	ks, err := keyStore.NewFromMnemonic(mne)
	require.Nil(t, err)
	
	km := keyManager.CreateFromKeyStore(ks)
	
	pubKeyStr, err := km.IdentityPublicKey()
	require.Nil(t, err)
	
	pubKeyRaw, err := hex.DecodeString(pubKeyStr)
	require.Nil(t, err)
	
	m := Message {
		Type: "HI",
		SendAt: time.Now(),
		AdditionalData: map[string]string{
			"key": "value",
		},
		DoubleratchetMessage: doubleratchet.Message{},
		IDPubKey: pubKeyRaw,
	}
	
	require.Nil(t, m.Sign(km))
	valid, err := m.VerifySignature()
	require.Nil(t, err)
	require.True(t, valid)
	
}
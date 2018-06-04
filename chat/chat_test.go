package chat

import (
	"github.com/Bit-Nation/panthalassa/keyManager"
	"github.com/Bit-Nation/panthalassa/keyStore"
	"github.com/Bit-Nation/panthalassa/mnemonic"
	"github.com/Bit-Nation/x3dh"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestEncryptDecryptX3DHSecretTest(t *testing.T) {

	mne, err := mnemonic.New()
	require.Nil(t, err)

	ks, err := keyStore.NewFromMnemonic(mne)
	require.Nil(t, err)

	km := keyManager.CreateFromKeyStore(ks)

	s := x3dh.SharedSecret{
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}

	// encrypt
	encryptedSecret, err := EncryptX3DHSecret(s, km)
	require.Nil(t, err)

	// decrypt
	decryptedSecret, err := DecryptX3DHSecret(encryptedSecret, km)
	require.Nil(t, err)

	require.Equal(t, s, decryptedSecret)

}

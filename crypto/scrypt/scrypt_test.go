package scrypt

import (
	"errors"
	"testing"

	"github.com/Bit-Nation/panthalassa/crypto/aes"
	"github.com/stretchr/testify/require"
)

func TestScryptValueKeyCreation(t *testing.T) {

	sV, err := makeScryptKey([]byte("a"))
	require.Nil(t, err)

	require.Equal(t, n, sV.N)
	require.Equal(t, p, sV.P)
	require.Equal(t, r, sV.R)
	require.Equal(t, saltLength, len(sV.Salt))
	require.Equal(t, 32, len(sV.key))

}

func TestSuccessEncryptAndDecrypt(t *testing.T) {

	value := []byte("i am the value")
	key := []byte("password")

	//create cipher text
	cipherText, err := NewCipherText(value, key)
	require.Nil(t, err)

	// mock cfb decrypt function
	// it is not supposed to get called
	cfbDecrypt = func(cipherText aes.CipherText, key aes.Secret) (aes.PlainText, error) {
		t.Error("I should not be called. CFBDecrypt should only be used with legacy scrypt cipher texts")
		t.FailNow()
		return []byte{}, nil
	}

	// Version one of the scrypt cipher text uses
	// aes cipher text version two which uses aes ctr
	require.Equal(t, uint8(1), cipherText.Version)
	require.Equal(t, uint8(2), cipherText.CipherText.Version)

	//decrypt cipher text
	plainText, err := DecryptCipherText(cipherText, key)
	require.Nil(t, err)
	require.Equal(t, string(value), string(plainText))

	// reset cfb decrypt function
	cfbDecrypt = aes.CFBDecrypt
}

func TestFailDecryption(t *testing.T) {

	value := []byte("i am the value")
	key := []byte("password")

	//create cipher text
	ethKey, err := NewCipherText(value, key)
	require.Nil(t, err)

	//decrypt cipher text
	plainText, err := DecryptCipherText(ethKey, []byte("word_password"))
	require.Error(t, errors.New("cipher: message authentication failed"), err)
	require.Equal(t, "", string(plainText))

}

// check if legacy decryption works as expected
func TestDecryptOldVersion(t *testing.T) {

	value := []byte("i am the value")
	key := []byte("password")

	ct, err := NewCipherText(value, key)
	require.Nil(t, err)

	// set to old version (deprecated)
	ct.Version = uint8(0)

	// mock the cfbDecrypt function
	cfbDecrypt = func(cipherText aes.CipherText, key aes.Secret) (aes.PlainText, error) {
		// there is no need to check parameters
		// the cbf function will only be called
		// in the case the we are using an old
		// scrypt cipher text
		return []byte("I am the mock plain text"), nil
	}

	plainText, err := DecryptCipherText(ct, key)
	require.Nil(t, err)

	require.Equal(t, string("I am the mock plain text"), string(plainText))

	// reset cfb decrypt function
	cfbDecrypt = aes.CFBDecrypt
}

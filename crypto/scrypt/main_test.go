package scrypt

import (
	"errors"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestScryptValueKeyCreation(t *testing.T) {

	keyLength := 100

	sV, e := Scrypt("a", keyLength)

	if e != nil {
		t.Error(e)
	}

	if sV.N != ScryptN {
		t.Errorf("Expected N to be: %d got %d", ScryptN, sV.N)
	}

	if sV.P != ScryptP {
		t.Errorf("Expected P to be: %d got %d", ScryptP, sV.P)
	}

	if sV.R != ScryptR {
		t.Errorf("Expected R to be: %d got %d", ScryptR, sV.R)
	}

	if len(sV.Salt) != ScryptSaltLength {
		t.Errorf("Expect Salt to have lenght: %d got %d", ScryptSaltLength, len(sV.Salt))
	}

	if len(sV.key) != keyLength {
		t.Errorf("key lenght must match the length of: %d", keyLength)
	}

}

func TestScryptCipherText_Export(t *testing.T) {

	s := ScryptCipherText{}

	encryptedValue, e := s.Export()

	require.Nil(t, e)

	require.Equal(t, "{\"cipher_text\":\"\",\"scrypt_key\":{\"n\":0,\"r\":0,\"p\":0,\"key_len\":0,\"salt\":null}}", encryptedValue)

}

func TestSuccessEncryptAndDecrypt(t *testing.T) {

	value := "i am the value"
	key := "password"

	//create cipher text
	cipherText, err := NewCipherText(value, key)
	require.Nil(t, err)

	//decrypt cipher text
	plainText, err := DecryptCipherText(cipherText, key)
	require.Nil(t, err)
	require.Equal(t, value, plainText)

}

func TestFailDecryption(t *testing.T) {

	value := "i am the value"
	key := "password"

	//create cipher text
	ethKey, err := NewCipherText(value, key)
	require.Nil(t, err)

	//decrypt cipher text
	plainText, err := DecryptCipherText(ethKey, "passwordddd")
	require.Error(t, errors.New("cipher: message authentication failed"), err)
	require.Equal(t, "", plainText)

}

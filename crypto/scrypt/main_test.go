package scrypt

import (
	"errors"
	"testing"

	require "github.com/stretchr/testify/require"
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

func TestScryptCipherText_Export(t *testing.T) {

	s := CipherText{}

	encryptedValue, e := s.Export()

	require.Nil(t, e)

	require.Equal(t, "{\"cipher_text\":{\"iv\":null,\"cipher_text\":null,\"mac\":null,\"v\":0},\"scrypt_key\":{\"n\":0,\"r\":0,\"p\":0,\"key_len\":0,\"salt\":null}}", encryptedValue)

}

func TestSuccessEncryptAndDecrypt(t *testing.T) {

	value := []byte("i am the value")
	key := []byte("password")

	//create cipher text
	cipherText, err := NewCipherText(value, key)
	require.Nil(t, err)

	//decrypt cipher text
	plainText, err := DecryptCipherText(cipherText, key)
	require.Nil(t, err)
	require.Equal(t, string(value), string(plainText))

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

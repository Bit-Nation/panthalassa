package aes

import (
	"errors"
	"github.com/stretchr/testify/require"
	"testing"
)

//Test the encrypt and decrypt function in one batch
func TestSuccessEncryptDecrypt(t *testing.T) {

	var secret, value string

	secret = "11111111111111111111111111111111"
	value = "I am the value"

	//Encrypt
	cipherText, e := Encrypt(value, secret)
	require.Nil(t, e)

	//Decrypt
	res, err := Decrypt(cipherText, secret)
	require.Nil(t, err)

	//Decrypted value must match the given value
	require.Equal(t, value, res)

}

func TestFailedDecryption(t *testing.T) {

	var secret, value string

	secret = "11111111111111111111111111111111"
	value = "I am the value"

	//Encrypt
	cipherText, e := Encrypt(value, secret)
	require.Nil(t, e)

	//Decrypt
	res, err := Decrypt(cipherText, "11111111111111111111111111111112")
	require.Equal(t, "", res)
	require.Error(t, errors.New("cipher: message authentication failed"), err)

}

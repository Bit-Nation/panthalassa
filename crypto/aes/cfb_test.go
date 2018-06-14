package aes

import (
	"testing"

	require "github.com/stretchr/testify/require"
)

//Test the encrypt and decrypt function in one batch
func CFBTestSuccessEncryptDecrypt(t *testing.T) {

	secret := Secret{
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}
	value := []byte("I am the value")

	//Encrypt
	cipherText, e := CFBEncrypt(value, secret)
	require.Nil(t, e)

	//Decrypt
	res, err := CFBDecrypt(cipherText, secret)
	require.Nil(t, err)

	//Decrypted value must match the given value
	require.Equal(t, string(value), string(res))

}

func CFBTestFailedDecryption(t *testing.T) {

	secret := Secret{
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}
	value := []byte("I am the plain text")

	// encrypt
	cipherText, e := CFBEncrypt(value, secret)
	require.Nil(t, e)

	// change last byte to fail on decryption
	secret[31] = 0x01

	// decrypt
	plainText, err := CFBDecrypt(cipherText, secret)
	require.Equal(t, PlainText{}, plainText)
	require.EqualError(t, err, "invalid key - message authentication failed")

}

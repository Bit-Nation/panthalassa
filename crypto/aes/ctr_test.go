package aes

import (
	"bytes"
	"encoding/hex"
	"testing"

	require "github.com/stretchr/testify/require"
)

// test vectors taken from here: https://nvlpubs.nist.gov/nistpubs/Legacy/SP/nistspecialpublication800-38a.pdf page 57

type testVector struct {
	iv         string
	plainText  string
	cipherText string
	mac        string
}

func TestCTREncryptDecrypt(t *testing.T) {

	secret := Secret{
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x10, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}
	value := []byte("I am the value")

	// encrypt
	ct, err := CTREncrypt(value, secret)
	require.Nil(t, err)

	require.Equal(t, uint8(2), ct.Version)

	// decrypt
	plainText, err := CTRDecrypt(ct, secret)
	require.Nil(t, err)

	// make sure we got what we expect
	require.Equal(t, string(value), string(plainText))

}

func TestCTREncryptDecryptFail(t *testing.T) {

	secret := Secret{
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x10, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}
	value := []byte("I am the value")

	// encrypt
	ct, err := CTREncrypt(value, secret)
	require.Nil(t, err)
	require.Equal(t, uint8(2), ct.Version)

	// change byte to make sure it fails
	secret[3] = 0x10

	// expect decrypt to fail
	plainText, err := CTRDecrypt(ct, secret)
	require.EqualError(t, MacError, err.Error())
	require.Equal(t, PlainText{}, plainText)

}

func TestCTREncrypt(t *testing.T) {

	// encryption key
	encryptionKey, err := hex.DecodeString("603deb1015ca71be2b73aef0857d77811f352c073b6108d72d9810a30914dff4")
	require.Nil(t, err)

	// aes ctr test vectors
	testVectors := []testVector{
		testVector{
			iv:         "f0f1f2f3f4f5f6f7f8f9fafbfcfdfeff",
			plainText:  "6bc1bee22e409f96e93d7e117393172a",
			cipherText: "601ec313775789a5b7a7f504bbf3d228",
		},
		testVector{
			iv:         "f0f1f2f3f4f5f6f7f8f9fafbfcfdff00",
			plainText:  "ae2d8a571e03ac9c9eb76fac45af8e51",
			cipherText: "f443e3ca4d62b59aca84e990cacaf5c5",
		},
		testVector{
			iv:         "f0f1f2f3f4f5f6f7f8f9fafbfcfdff01",
			plainText:  "30c81c46a35ce411e5fbc1191a0a52ef",
			cipherText: "2b0930daa23de94ce87017ba2d84988d",
		},
		testVector{
			iv:         "f0f1f2f3f4f5f6f7f8f9fafbfcfdff02",
			plainText:  "f69f2445df4f9b17ad2b417be66c3710",
			cipherText: "dfc9c58db67aada613c2dd08457941a6",
		},
	}

	for _, v := range testVectors {

		// initialisation vector
		iv, err := hex.DecodeString(v.iv)
		require.Nil(t, err)
		ctrRandReader = bytes.NewReader(iv)

		// plain text
		plainText, err := hex.DecodeString(v.plainText)
		require.Nil(t, err)

		// cipher text
		cipherText, err := hex.DecodeString(v.cipherText)
		require.Nil(t, err)

		// create secret key
		secKey := Secret{}
		copy(secKey[:], encryptionKey)

		// encrypt
		ct, err := CTREncrypt(plainText, secKey)
		require.Nil(t, err)

		require.Equal(t, hex.EncodeToString(cipherText), hex.EncodeToString(ct.CipherText))

	}

}

func TestCTRDecrypt(t *testing.T) {

	// encryption key
	encryptionKey, err := hex.DecodeString("603deb1015ca71be2b73aef0857d77811f352c073b6108d72d9810a30914dff4")
	require.Nil(t, err)

	// aes ctr test vectors
	testVectors := []testVector{
		testVector{
			iv:         "f0f1f2f3f4f5f6f7f8f9fafbfcfdfeff",
			plainText:  "6bc1bee22e409f96e93d7e117393172a",
			cipherText: "601ec313775789a5b7a7f504bbf3d228",
			mac:        "3465d5a7141a0f265ce4148f7b2c00461167662e2892ac15ca49a13474e9759c",
		},
		testVector{
			iv:         "f0f1f2f3f4f5f6f7f8f9fafbfcfdff00",
			plainText:  "ae2d8a571e03ac9c9eb76fac45af8e51",
			cipherText: "f443e3ca4d62b59aca84e990cacaf5c5",
			mac:        "31493d770e783d72d0504497cc53c7504be3cc9a2f06e6f41485b38692b1b34c",
		},
		testVector{
			iv:         "f0f1f2f3f4f5f6f7f8f9fafbfcfdff01",
			plainText:  "30c81c46a35ce411e5fbc1191a0a52ef",
			cipherText: "2b0930daa23de94ce87017ba2d84988d",
			mac:        "dff4a5f39af148fb68523150c39098d988b43392b37ed1b38fd58525b8432dac",
		},
		testVector{
			iv:         "f0f1f2f3f4f5f6f7f8f9fafbfcfdff02",
			plainText:  "f69f2445df4f9b17ad2b417be66c3710",
			cipherText: "dfc9c58db67aada613c2dd08457941a6",
			mac:        "858aba8f3a3141c8e7317128687c3d70069d9b7287271bc565d704d6c6298911",
		},
	}

	for _, v := range testVectors {

		// initialisation vector
		iv, err := hex.DecodeString(v.iv)
		require.Nil(t, err)
		ctrRandReader = bytes.NewReader(iv)

		// mac
		mac, err := hex.DecodeString(v.mac)
		require.Nil(t, err)

		// cipher text
		cipherText, err := hex.DecodeString(v.cipherText)
		require.Nil(t, err)

		// create secret key
		secKey := Secret{}
		copy(secKey[:], encryptionKey)

		// encrypt
		plain, err := CTRDecrypt(CipherText{
			IV:         iv,
			Mac:        mac,
			CipherText: cipherText,
			Version:    uint8(2),
		}, secKey)
		require.Nil(t, err)

		require.Equal(t, v.plainText, hex.EncodeToString(plain))

	}

}

package aes

import (
	"bytes"
	"encoding/hex"
	"testing"

	require "github.com/stretchr/testify/require"
)

// Test the encrypt and decrypt function in one batch
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

func TestCFBEncrypt(t *testing.T) {

	// encryption key
	encryptionKey, err := hex.DecodeString("603deb1015ca71be2b73aef0857d77811f352c073b6108d72d9810a30914dff4")
	require.Nil(t, err)

	// aes cfb test vectors
	testVectors := []testVector{
		testVector{
			iv:         "000102030405060708090a0b0c0d0e0f",
			plainText:  "6bc1bee22e409f96e93d7e117393172a",
			cipherText: "dc7e84bfda79164b7ecd8486985d3860",
		},
		testVector{
			iv:         "dc7e84bfda79164b7ecd8486985d3860",
			plainText:  "ae2d8a571e03ac9c9eb76fac45af8e51",
			cipherText: "39ffed143b28b1c832113c6331e5407b",
		},
		testVector{
			iv:         "39ffed143b28b1c832113c6331e5407b",
			plainText:  "30c81c46a35ce411e5fbc1191a0a52ef",
			cipherText: "df10132415e54b92a13ed0a8267ae2f9",
		},
		testVector{
			iv:         "df10132415e54b92a13ed0a8267ae2f9",
			plainText:  "f69f2445df4f9b17ad2b417be66c3710",
			cipherText: "75a385741ab9cef82031623d55b1e471",
		},
	}

	for _, v := range testVectors {

		// initialisation vector
		iv, err := hex.DecodeString(v.iv)
		require.Nil(t, err)
		cfbRandReader = bytes.NewReader(iv)

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
		ct, err := CFBEncrypt(plainText, secKey)
		require.Nil(t, err)

		require.Equal(t, hex.EncodeToString(cipherText), hex.EncodeToString(ct.CipherText))

	}

}

func TestCFBDecrypt(t *testing.T) {

	// encryption key
	encryptionKey, err := hex.DecodeString("603deb1015ca71be2b73aef0857d77811f352c073b6108d72d9810a30914dff4")
	require.Nil(t, err)

	// aes cfb test vectors
	testVectors := []testVector{
		testVector{
			iv:         "000102030405060708090a0b0c0d0e0f",
			plainText:  "6bc1bee22e409f96e93d7e117393172a",
			cipherText: "dc7e84bfda79164b7ecd8486985d3860",
			mac:        "43517b9531fdb8565ce2a6af1cbb24930cb7ec2b3685701cfac6731bb861cb60",
		},
		testVector{
			iv:         "dc7e84bfda79164b7ecd8486985d3860",
			plainText:  "ae2d8a571e03ac9c9eb76fac45af8e51",
			cipherText: "39ffed143b28b1c832113c6331e5407b",
			mac:        "de571950b74a72055810fad70421676aeadd651d6bd3ba6ad0348a7c7c9ebc8e",
		},
		testVector{
			iv:         "39ffed143b28b1c832113c6331e5407b",
			plainText:  "30c81c46a35ce411e5fbc1191a0a52ef",
			cipherText: "df10132415e54b92a13ed0a8267ae2f9",
			mac:        "6ce26e57380b71dec45177dda9d1d322e1dceba40b462b6c15990a43bb1b2ec2",
		},
		testVector{
			iv:         "df10132415e54b92a13ed0a8267ae2f9",
			plainText:  "f69f2445df4f9b17ad2b417be66c3710",
			cipherText: "75a385741ab9cef82031623d55b1e471",
			mac:        "73fb351bba3a8299d2271ff489dc4ed2412e69653b9877bc19a6b408c51e83da",
		},
	}

	for _, v := range testVectors {

		// initialisation vector
		iv, err := hex.DecodeString(v.iv)
		require.Nil(t, err)
		cfbRandReader = bytes.NewReader(iv)

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
		plain, err := CFBDecrypt(CipherText{
			IV:         iv,
			Mac:        mac,
			CipherText: cipherText,
		}, secKey)
		require.Nil(t, err)

		require.Equal(t, v.plainText, hex.EncodeToString(plain))

	}

}

package aes

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"io"
)

var ctrRandReader io.Reader = rand.Reader

// encrypt plain text by given key using AES CTR 256
func CTREncrypt(plainText PlainText, secret Secret) (CipherText, error) {

	// block
	block, err := aes.NewCipher(secret[:])
	if err != nil {
		return CipherText{}, err
	}

	// initialisation vector
	iv := make([]byte, 16)
	_, err = io.ReadFull(ctrRandReader, iv)
	if err != nil {
		return CipherText{}, err
	}

	// create cipher text
	cipherText := make([]byte, len(plainText))

	// encrypt
	stream := cipher.NewCTR(block, iv)
	stream.XORKeyStream(cipherText, plainText)

	// create mac
	mac := hmac.New(sha256.New, secret[:])
	if _, err := mac.Write(cipherText); err != nil {
		return CipherText{}, err
	}

	ct := CipherText{
		IV:         iv,
		CipherText: cipherText,
		Mac:        mac.Sum(nil),
		Version:    2,
	}

	return ct, err
}

// decrypt cipher text by given key
func CTRDecrypt(cipherText CipherText, key Secret) (PlainText, error) {

	// create block
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return PlainText{}, err
	}

	// validate key
	valid, err := cipherText.ValidMAC(key)
	if err != nil {
		return PlainText{}, err
	}
	if !valid {
		return PlainText{}, MacError
	}

	// decrypt
	plainText := make(PlainText, len(cipherText.CipherText))
	stream := cipher.NewCTR(block, cipherText.IV)
	stream.XORKeyStream(plainText, cipherText.CipherText)

	return plainText, nil

}

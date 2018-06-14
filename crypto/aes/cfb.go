package aes

//Taken from this example: https://gist.github.com/cannium/c167a19030f2a3c6adbb5a5174bea3ff

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"io"
)

// encrypt a plain text with given secret
func CFBEncrypt(plainText PlainText, key Secret) (CipherText, error) {

	block, err := aes.NewCipher(key[:])
	if err != nil {
		return CipherText{}, err
	}

	cipherText := make([]byte, len(plainText))
	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return CipherText{}, err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(cipherText, plainText)

	mac := hmac.New(sha256.New, key[:])
	if _, err := mac.Write(cipherText); err != nil {
		return CipherText{}, err
	}

	ct := CipherText{
		IV:         iv,
		CipherText: cipherText,
		Mac:        mac.Sum(nil),
		Version:    1,
	}

	return ct, nil

}

// decrypt a cipher text with given secret
func CFBDecrypt(cipherText CipherText, key Secret) (PlainText, error) {

	// create block
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return PlainText{}, err
	}

	valid, err := cipherText.ValidMAC(key)
	if err != nil {
		return PlainText{}, err
	}
	if !valid {
		return PlainText{}, errors.New("invalid key - message authentication failed")
	}

	stream := cipher.NewCFBDecrypter(block, cipherText.IV)

	cc := cipherText.CipherText

	// XORKeyStream can work in-place if the two arguments are the same.
	stream.XORKeyStream(cc, cc)

	return cc, nil
}

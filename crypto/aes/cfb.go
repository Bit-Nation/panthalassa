package aes

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"io"
)

var cfbRandReader io.Reader = rand.Reader

// encrypt a plain text with given secret
// Deprecated: use CTREncrypt instead
func CFBEncrypt(plainText PlainText, key Secret) (CipherText, error) {

	// create block
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return CipherText{}, err
	}

	// create IV
	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(cfbRandReader, iv); err != nil {
		return CipherText{}, err
	}

	// create cfb encrypter
	stream := cipher.NewCFBEncrypter(block, iv)
	cipherText := make([]byte, len(plainText))
	stream.XORKeyStream(cipherText, plainText)

	// create mac
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
// Deprecated: use CFBDecrypt instead
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
		return PlainText{}, MacError
	}

	stream := cipher.NewCFBDecrypter(block, cipherText.IV)

	cc := cipherText.CipherText

	// XORKeyStream can work in-place if the two arguments are the same.
	stream.XORKeyStream(cc, cc)

	return cc, nil
}

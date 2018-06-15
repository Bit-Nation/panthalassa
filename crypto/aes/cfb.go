package aes

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
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

	// cipher text
	ct := CipherText{
		IV:         iv,
		CipherText: cipherText,
		Version:    1,
	}

	// create mac
	mac, err := vOneMac(ct, key)
	if err != nil {
		return CipherText{}, err
	}
	ct.Mac = mac

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

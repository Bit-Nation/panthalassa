package aes

//Taken from this example: https://gist.github.com/cannium/c167a19030f2a3c6adbb5a5174bea3ff

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"io"
)

type PlainText []byte

type CipherText struct {
	IV         []byte `json:"iv"`
	CipherText []byte `json:"cipher_text"`
	Mac        []byte `json:"mac"`
	Version    uint8  `json:"v"`
}

// verify MAC of cipher text
func (c CipherText) ValidMAC(s Secret) (bool, error) {
	h := hmac.New(sha256.New, s[:])
	if _, err := h.Write(c.CipherText); err != nil {
		return false, err
	}
	return hmac.Equal(c.Mac, h.Sum(nil)), nil
}

func (c CipherText) Marshal() ([]byte, error) {
	return json.Marshal(c)
}

type Secret [32]byte

// encrypt a plain text with given secret
func Encrypt(plainText PlainText, key Secret) (CipherText, error) {

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
func Decrypt(cipherText CipherText, key Secret) (PlainText, error) {

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

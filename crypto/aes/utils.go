package aes

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/json"
	"errors"
)

type PlainText []byte
type Secret [32]byte

type CipherText struct {
	IV         []byte `json:"iv"`
	CipherText []byte `json:"cipher_text"`
	Mac        []byte `json:"mac"`
	Version    uint8  `json:"v"`
}

var MacError = errors.New("invalid key - message authentication failed")

// verify MAC of cipher text
func (c CipherText) ValidMAC(s Secret) (bool, error) {
	h := hmac.New(sha256.New, s[:])
	if _, err := h.Write(c.CipherText); err != nil {
		return false, err
	}
	return hmac.Equal(c.Mac, h.Sum(nil)), nil
}

// marshal cipher text
func (c CipherText) Marshal() ([]byte, error) {
	return json.Marshal(c)
}

// unmarshal cipher text
func Unmarshal(rawCipherText []byte) (CipherText, error) {
	var ct CipherText
	if err := json.Unmarshal(rawCipherText, &ct); err != nil {
		return CipherText{}, err
	}
	return ct, nil
}

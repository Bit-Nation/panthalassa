package aes

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
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

// create version one MAC
// based on cipher text
var vOneMac = func(ct CipherText, secret Secret) ([]byte, error) {
	if ct.Version != 1 {
		return nil, errors.New("cipher text must be of version one")
	}
	h := hmac.New(sha256.New, secret[:])
	if _, err := h.Write(ct.CipherText); err != nil {
		return nil, err
	}
	return h.Sum(nil), nil
}

// create version two of MAC
// cipher text + IV + Version
var vTwoMac = func(ct CipherText, secret Secret) ([]byte, error) {
	if ct.Version != 2 {
		return nil, errors.New("cipher text must be of version two")
	}
	h := hmac.New(sha256.New, secret[:])
	if _, err := h.Write(ct.CipherText); err != nil {
		return nil, err
	}
	if _, err := h.Write(ct.IV); err != nil {
		return nil, err
	}
	if _, err := h.Write([]byte{ct.Version}); err != nil {
		return nil, err
	}
	return h.Sum(nil), nil
}

// verify MAC of cipher text
func (c CipherText) ValidMAC(s Secret) (bool, error) {

	switch c.Version {
	// this has only been used for version 1 of the cipher text
	// version one only used the cipher text as the hmac message
	case uint8(1):
		mac, err := vOneMac(c, s)
		if err != nil {
			return false, err
		}
		return hmac.Equal(c.Mac, mac), nil
	// version two use the cipher text + IV as the message
	case uint8(2):
		mac, err := vTwoMac(c, s)
		if err != nil {
			return false, err
		}
		return hmac.Equal(c.Mac, mac), nil
	default:
		return false, errors.New(fmt.Sprintf("failed to verify MAC since we don't know how to handle version: %d", c.Version))
	}

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

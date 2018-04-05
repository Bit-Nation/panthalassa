package bip39

import (
	"crypto/rand"
	bip39 "github.com/tyler-smith/go-bip39"
)

//Wordlist from bip30 package
var WordList = bip39.WordList

func newMnemonic(entropy []byte) (string, error) {
	return bip39.NewMnemonic(entropy)
}

func NewMnemonic() (string, error) {
	entropy := make([]byte, 32)
	if _, err := rand.Read(entropy); err != nil {
		return "", err
	}
	m, err := newMnemonic(entropy)
	if err != nil {
		return "", err
	}
	return m, nil
}

//Generate new seed of mnemonic and password
func NewSeed(mnemonic string, password string) []byte {
	return bip39.NewSeed(mnemonic, password)
}

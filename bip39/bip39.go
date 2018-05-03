package bip39

import (
	"crypto/rand"
	"errors"
	bip39 "github.com/tyler-smith/go-bip39"
	"strings"
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
func NewSeed(mnemonic string, password string) ([]byte, error) {

	if !ValidMnemonic(mnemonic) {
		return []byte(""), errors.New("got invalid mnemonic")
	}

	return newSeed(mnemonic, password), nil

}

func newSeed(mnemonic, password string) []byte {
	return bip39.NewSeed(mnemonic, password)
}

//Check if an mnemonic is valid or not
func ValidMnemonic(mnemonic string) bool {

	words := strings.Split(mnemonic, " ")

	if len(words) != 24 {
		return false
	}

	for _, word := range words {

		exist := false

		for _, bip39Word := range EnglishWordList {

			if word == bip39Word {
				exist = true
				break
			}

		}

		if !exist {
			return false
		}

	}

	return true

}

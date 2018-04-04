package bip39

import (
	"crypto/rand"
	bip39 "github.com/tyler-smith/go-bip39"
	"fmt"
)

//@todo add test's (for the imported lib as well)
func NewMnemonic() (string, error) {
	entropy := make([]byte, 32)
	if _, err := rand.Read(entropy); err != nil {
		return "", err
	}
	fmt.Println(entropy)
	m, err := bip39.NewMnemonic(entropy)
	if err != nil {
		return "", err
	}
	return m, nil
}

//@todo add test's (for the imported lib as well)
func NewSeed(mnemonic string, password string) []byte {
	return bip39.NewSeed(mnemonic, password)
}
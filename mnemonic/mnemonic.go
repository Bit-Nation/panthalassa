package mnemonic

import (
	"crypto/rand"
	"errors"

	bip39 "github.com/tyler-smith/go-bip39"
)

type Mnemonic struct {
	mnemonic string
}

//Create Mnemonic from string
//@todo should also check if = 24 word's long
func FromString(mnemonic string) (Mnemonic, error) {

	if !bip39.IsMnemonicValid(mnemonic) {
		return Mnemonic{}, errors.New("got invalid mnemonic")
	}

	return Mnemonic{
		mnemonic: mnemonic,
	}, nil
}

//Mnemonic factory
func New() (Mnemonic, error) {

	//Secure random numbers
	entropy := make([]byte, 32)
	if _, err := rand.Read(entropy); err != nil {
		return Mnemonic{}, err
	}

	//Mnemonic as word string
	m, err := bip39.NewMnemonic(entropy)
	if err != nil {
		return Mnemonic{}, err
	}

	return FromString(m)
}

//Generate new seed of mnemonic and password
func (m Mnemonic) NewSeed(password string) ([]byte, error) {

	if !bip39.IsMnemonicValid(m.mnemonic) {
		panic("Got invalid mnemonic. This shouldn't happen since we check it in the FromString method")
	}

	return bip39.NewSeed(m.mnemonic, password), nil

}

func (m Mnemonic) String() string {

	if !bip39.IsMnemonicValid(m.mnemonic) {
		panic("Got invalid mnemonic. This shouldn't happen since we check it in the FromString method")
	}

	return m.mnemonic
}

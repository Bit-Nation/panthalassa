package rootKey

import (
	"errors"
	bip39 "github.com/Bit-Nation/panthalassa/bip39"
	crypto "github.com/Bit-Nation/panthalassa/crypto"
	"strings"
	"unicode"
)

type RootKey struct {
	words [24]string
}

//Get the plain root key with white separation
func (rK RootKey) GetWithWhiteSpace() string {
	return strings.Join(rK.words[:], " ")
}

//Get the plain root key with comma separation
func (rK RootKey) getWithComma() string {
	return strings.Join(rK.words[:], ",")
}

//Export the root key with given password
func (rK RootKey) Export(pw string) (string, error) {

	//mnemonic as string separated with white space
	mnemonicStr := rK.getWithComma()

	//Encrypt the mnemonic with the password
	cipherText, err := crypto.NewScryptCipherText(pw, mnemonicStr)

	if err != nil {
		return "", err
	}

	return cipherText, nil

}

//Create a new RootKey instance
func NewRootKey() (RootKey, error) {

	mnemonic, err := bip39.NewMnemonic()

	if err != nil {
		return RootKey{}, err
	}

	words := strings.Split(mnemonic, " ")

	//strip white space from words,
	//in case there is a white space left
	for _, word := range words {
		mnemonic = strings.Map(func(r rune) rune {
			if unicode.IsSpace(r) == true {
				return -1
			}
			return r
		}, word)
	}

	rk := RootKey{}

	copy(rk.words[:], words)

	return rk, nil
}

//Create an RootKey instance form the encryptedRootKey and the password
func RootKeyFromCipherText(encryptedRootKey string, pw string) (RootKey, error) {

	mnemonic, err := crypto.DecryptScryptCipherText(pw, encryptedRootKey)

	if err != nil {
		return RootKey{}, err
	}

	words := strings.Split(mnemonic, ",")

	if len(words) != 24 {
		return RootKey{}, errors.New("Given encryptedRootKey doesn't seem to be valid")
	}

	rk := RootKey{}
	copy(rk.words[:], words)

	return rk, nil

}

package panthalassa

import (
	"errors"

	cid "github.com/Bit-Nation/panthalassa/crypto/cid"
	scrypt "github.com/Bit-Nation/panthalassa/crypto/scrypt"
	"github.com/Bit-Nation/panthalassa/keyManager"
	"github.com/Bit-Nation/panthalassa/keyStore"
	mnemonic "github.com/Bit-Nation/panthalassa/mnemonic"
	"github.com/tyler-smith/go-bip39"
	"strings"
)

//Encrypt's data
//A password and a data string is required
//A key is derived from the password with scrypt
//The derived key is then used to encrypt the data
//with AES256
//Returned is the whole cipher text included with the scrypt parameters
func ScryptEncrypt(data, pw, pwConfirm string) (string, error) {

	if pw != pwConfirm {
		return "", errors.New("password mismatch")
	}

	return scrypt.NewCipherText(data, pw)
}

//Decrypt scrypt cipher text
//Need's a string value like the one returned from ScryptEncrypt
func ScryptDecrypt(data, pw string) (string, error) {
	return scrypt.DecryptCipherText(data, pw)
}

//Creates an new set of encrypted account key's
func NewAccountKeys(pw, pwConfirm string) (string, error) {

	//Create mnemonic
	mn, err := mnemonic.New()
	if err != nil {
		return "", err
	}

	//Create KeyStore
	ks, err := keyStore.NewFromMnemonic(mn)
	if err != nil {
		return "", err
	}

	km := keyManager.CreateFromKeyStore(ks)
	return km.Export(pw, pwConfirm)
}

//Create new account store from mnemonic
//This can e.g. be used in case you need to recover your account
func NewAccountKeysFromMnemonic(mne, pw, pwConfirm string) (string, error) {

	//Create mnemonic
	mn, err := mnemonic.FromString(mne)
	if err != nil {
		return "", err
	}

	//Create key store from mnemonic
	ks, err := keyStore.NewFromMnemonic(mn)
	if err != nil {
		return "", err
	}

	//Create keyManager
	km := keyManager.CreateFromKeyStore(ks)
	return km.Export(pw, pwConfirm)
}

//Get the CID of a value with
//sha3 512 as a base64 string
func CIDSha512(value string) (string, error) {
	return cid.Sha512(value)
}

//Get the CID of a value with
//sha3 256 as a base64 string
func CIDSha256(value string) (string, error) {
	return cid.Sha256(value)
}

//Check if CID is valid
func IsValidCID(c string) bool {
	return cid.IsValidCid(c)
}

//Check if mnemonic is valid
func IsValidMnemonic(mne string) bool {

	words := strings.Split(mne, " ")

	if len(words) != 24 {
		return false
	}

	return bip39.IsMnemonicValid(mne)
}

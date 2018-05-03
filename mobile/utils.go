package panthalassa

import (
	"errors"

	"github.com/Bit-Nation/panthalassa/crypto"
	"github.com/Bit-Nation/panthalassa/keyManager"
	"github.com/Bit-Nation/panthalassa/keyStore"
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

	return crypto.NewScryptCipherText(pw, data)
}

//Decrypt scrypt cipher text
//Need's a string value like the one returned from ScryptEncrypt
func ScryptDecrypt(data, pw string) (string, error) {
	return crypto.NewScryptCipherText(pw, data)
}

//Creates an new set of encrypted account key's
func NewAccountKeys(pw, pwConfirm string) (string, error) {
	ks, err := keyStore.NewKeyStoreFactory()
	if err != nil {
		return "", err
	}
	km := keyManager.CreateFromKeyStore(ks)
	return km.Export(pw, pwConfirm)
}

//Create new account store from mnemonic
//This can e.g. be used in case you need to recover your account
func NewAccountKeysFromMnemonic(mnemonic, pw, pwConfirm string) (string, error) {
	ks, err := keyStore.NewFromMnemonic(mnemonic)
	if err != nil {
		return "", err
	}
	km := keyManager.CreateFromKeyStore(ks)
	return km.Export(pw, pwConfirm)
}

//Get the CID of a value with
//sha3 512 as a base64 string
func CIDSha512(value string) (string, error) {
	return crypto.CIDSha512(value)
}

//Get the CID of a value with
//sha3 256 as a base64 string
func CIDSha256(value string) (string, error) {
	return crypto.CIDSha256(value)
}

//Check if CID is valid
func IsValidCID(cid string) bool {
	return crypto.IsValidCid(cid)
}

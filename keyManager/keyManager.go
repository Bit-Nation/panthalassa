package keyManager

import (
	"encoding/json"
	"errors"
	crypto "github.com/Bit-Nation/panthalassa/crypto"
	ks "github.com/Bit-Nation/panthalassa/keyStore"
)

type KeyManager struct {
	keyStore *ks.KeyStore
	account  accountKeyStore
}

type accountKeyStore struct {
	Password          string `json:"password"`
	EncryptedKeyStore string `json:"encrypted_key_store"`
	Version           uint8  `json:"version"`
}

//Open encrypted keystore with password
func OpenWithPassword(encryptedAccount, pw string) (*KeyManager, error) {

	//unmarshal encrypted account
	var acc accountKeyStore
	if err := json.Unmarshal([]byte(encryptedAccount), &acc); err != nil {
		return &KeyManager{}, err
	}

	//Decrypt key store
	jsonKeyStore, err := crypto.DecryptScryptCipherText(pw, acc.EncryptedKeyStore)
	if err != nil {
		return &KeyManager{}, err
	}

	//unmarshal key store
	keyStore, err := ks.FromJson(jsonKeyStore)
	if err != nil {
		return &KeyManager{}, err
	}

	return &KeyManager{
		keyStore: keyStore,
		account:  acc,
	}, nil

}

//Open account with mnemonic.
//This should only be used as a backup
func OpenWithMnemonic(account, mnemonic string) (*KeyManager, error) {

	//unmarshal encrypted account
	var acc accountKeyStore
	if err := json.Unmarshal([]byte(account), &acc); err != nil {
		return &KeyManager{}, err
	}

	//decrypt password with mnemonic
	pw, err := crypto.DecryptScryptCipherText(mnemonic, acc.Password)
	if err != nil {
		return &KeyManager{}, err
	}

	return OpenWithPassword(account, pw)

}

//Export the account
func (km KeyManager) Export(pw, pwConfirm string) (string, error) {

	//Exit if password's are not equal
	if pw != pwConfirm {
		return "", errors.New("password miss match")
	}

	//Marshal the keystore
	keyStore, err := km.keyStore.Marshal()
	if err != nil {
		return "", err
	}

	//encrypt key store with password
	encryptedKeyStore, err := crypto.NewScryptCipherText(pw, string(keyStore))
	if err != nil {
		return "", err
	}

	//encrypt password with mnemonic
	encryptedPassword, err := crypto.NewScryptCipherText(km.keyStore.GetMnemonic(), pw)

	//Marshal account
	acc, err := json.Marshal(accountKeyStore{
		Password:          encryptedPassword,
		EncryptedKeyStore: encryptedKeyStore,
		Version:           1,
	})

	return string(acc), err

}

//Get ethereum private key
func (km KeyManager) GetEthereumPrivateKey() (string, error) {
	return km.keyStore.GetKey("eth_private_key")
}

//Create new key manager from key store
func CreateFromKeyStore(ks *ks.KeyStore) *KeyManager {
	return &KeyManager{
		keyStore: ks,
	}
}

package keyManager

import (
	"errors"
	crypto "github.com/Bit-Nation/panthalassa/crypto"
	ks "github.com/Bit-Nation/panthalassa/keyStore"
)

type KeyManager struct {
	keyStore *ks.KeyStore
}

//Open an encrypted key store file
func Open(encryptedKeyStore, pw string) (*KeyManager, error) {

	//Decrypt key store
	jsonKeyStore, err := crypto.DecryptScryptCipherText(pw, encryptedKeyStore)
	if err != nil {
		return &KeyManager{}, nil
	}

	//transform json key store string to KeyStore
	keyStore, err := ks.FromJson(jsonKeyStore)

	if err != nil {
		return &KeyManager{}, nil
	}

	//
	return &KeyManager{
		keyStore: keyStore,
	}, nil

}

//Export the key store
func (km KeyManager) Export(pw, pwConfirm string) (string, error) {

	//Exit if password's are not equal
	if pw != pwConfirm {
		return "", errors.New("password miss match")
	}

	keyStoreJson, err := km.keyStore.Marshal()
	if err != nil {
		return "", err
	}
	return crypto.NewScryptCipherText(pw, string(keyStoreJson))
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

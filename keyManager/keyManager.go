package keyManager

import (
	rk "github.com/Bit-Nation/panthalassa/rootKey"
)

type KeyManager struct {
	rootKey rk.RootKey
}

//Create a new key manager from the encrypted root key
func NewKeyManager(encryptedRootKey, pw string) (*KeyManager, error) {

	rootKey, err := rk.RootKeyFromCipherText(encryptedRootKey, pw)

	if err != nil {
		return &KeyManager{}, nil
	}

	keyManager := KeyManager{
		rootKey: rootKey,
	}

	return &keyManager, nil
}

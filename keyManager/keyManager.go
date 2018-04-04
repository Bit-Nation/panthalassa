package keyManager

import (
	"encoding/hex"
	bip32 "github.com/Bit-Nation/panthalassa/bip32"
	rk "github.com/Bit-Nation/panthalassa/rootKey"
)

type KeyManager struct {
	rootKey rk.RootKey
	ethKey  bip32.Key
}

//Encode the eth bip32 key to a hex string
func (kM KeyManager) GetEthPrivateKey() string {
	return hex.EncodeToString(kM.ethKey.Key)
}

//Create a new key manager from the encrypted root key
func NewKeyManager(encryptedRootKey, pw string) (*KeyManager, error) {

	rootKey, err := rk.RootKeyFromCipherText(encryptedRootKey, pw)

	if err != nil {
		return &KeyManager{}, nil
	}

	masterKey, err := bip32.NewMasterKey([]byte(rootKey.GetWithWhiteSpace()))

	if err != nil {
		return &KeyManager{}, err
	}

	//Create ethereum key as specified
	ethKey, err := masterKey.NewChildKey(2 ^ 31 + 100)
	if err != nil {
		return &KeyManager{}, err
	}
	ethKey, err = ethKey.NewChildKey(2 ^ 31 + 10)
	if err != nil {
		return &KeyManager{}, err
	}

	keyManager := KeyManager{
		rootKey: rootKey,
		ethKey:  *ethKey,
	}

	return &keyManager, nil
}

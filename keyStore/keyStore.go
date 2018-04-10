package keyStore

import (
	jsonUtil "encoding/json"
	"errors"
	bip32 "github.com/Bit-Nation/panthalassa/bip32"
	bip39 "github.com/Bit-Nation/panthalassa/bip39"
)

//Ethereum private key validation rule
var ethPrivateKeyValidation = func(store KeyStore) error {

	//derive ethereum key
	seed := bip39.NewSeed(store.mnemonic, "pangea")
	key, err := bip32.NewMasterKey(seed)
	if err != nil {
		return err
	}
	key, err = bip32.Derive("m/100H/10H", key)
	if err != nil {
		return err
	}

	return nil
}

//A set of validation rule for each key store version
var validationRules = map[uint8][]func(ks KeyStore) error{
	1: []func(store KeyStore) error{
		ethPrivateKeyValidation,
	},
}

type KeyStore struct {
	mnemonic string            `json:mnemonic`
	keys     map[string]string `json:keys`
	version  uint8             `json:version`
}

func (ks KeyStore) GetKey(key string) (string, error) {

	_, exist := ks.keys[key]
	if exist == false {
		return "", errors.New("key does not exist")
	}

	return ks.keys[key], nil
}

//Validate the key store
func (ks KeyStore) validate() error {

	//Check if validation rules are present
	_, exist := validationRules[ks.version]
	if exist == false {
		return errors.New("couldn't find validation rules")
	}

	//Validate the key store
	for _, f := range validationRules[ks.version] {
		if f(ks) != nil {
			return f(ks)
		}
	}

	return nil
}

//Marshal key store
func (ks KeyStore) Marshal() ([]byte, error) {
	return jsonUtil.Marshal(ks)
}

//Convert json keystore to object
func FromJson(json string) (*KeyStore, error) {
	var keyStore KeyStore
	err := jsonUtil.Unmarshal([]byte(json), &keyStore)
	if err != nil {
		return &keyStore, err
	}

	if err := keyStore.validate(); err != nil {
		return &KeyStore{}, err
	}

	return &KeyStore{}, nil
}

func NewMnemonic(mnemonic string) {

}

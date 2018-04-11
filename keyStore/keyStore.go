package keyStore

import (
	"encoding/hex"
	jsonUtil "encoding/json"
	"errors"
	bip32 "github.com/Bit-Nation/panthalassa/bip32"
	bip39 "github.com/Bit-Nation/panthalassa/bip39"
)

//Ethereum private key validation rule
var ethPrivateKeyValidation = func(store KeyStore) error {

	//derive seed used for coins
	seed := bip39.NewSeed(store.mnemonic, "coins")
	//make it the master key
	key, err := bip32.NewMasterKey(seed)
	if err != nil {
		return err
	}
	//Derive the ethereum key as per spec
	key, err = bip32.Derive("m/100H/10H", key)
	if err != nil {
		return err
	}

	//compare key's
	hexKey, err := store.GetKey("eth_private_key")
	if err != nil {
		return err
	}
	if hex.EncodeToString(key.Key) != hexKey {
		return errors.New("derivation mismatch - ethereum private key from storage and derived one doesn't match")
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
	mnemonic string            `json:"mnemonic"`
	keys     map[string]string `json:"keys"`
	version  uint8             `json:"version"`
}

//Only used for json marshalling
type jsonKeyStore struct {
	Mnemonic string            `json:"mnemonic"`
	Keys     map[string]string `json:"keys"`
	Version  uint8             `json:"version"`
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
	for _, vR := range validationRules[ks.version] {
		if err := vR(ks); err != nil {
			return err
		}
	}

	return nil
}

//Marshal key store
func (ks KeyStore) Marshal() ([]byte, error) {

	if err := ks.validate(); err != nil {
		return []byte{}, err
	}

	jsk := jsonKeyStore{
		Mnemonic: ks.mnemonic,
		Keys:     ks.keys,
		Version:  ks.version,
	}

	return jsonUtil.Marshal(jsk)
}

//Convert json keystore to object
func FromJson(json string) (*KeyStore, error) {
	var jsk jsonKeyStore
	err := jsonUtil.Unmarshal([]byte(json), &jsk)
	if err != nil {
		return &KeyStore{}, err
	}

	//Create keystore form parsed json
	ks := KeyStore{
		mnemonic: jsk.Mnemonic,
		keys:     jsk.Keys,
		version:  jsk.Version,
	}

	//Exit on invalid key store
	if err := ks.validate(); err != nil {
		return &KeyStore{}, err
	}

	return &ks, nil
}

func NewMnemonic(mnemonic string) {

}

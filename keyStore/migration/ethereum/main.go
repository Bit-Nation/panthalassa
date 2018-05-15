package ethereum

import (
	"encoding/hex"
	"errors"

	bip32 "github.com/Bit-Nation/panthalassa/bip32"
	mnemonic "github.com/Bit-Nation/panthalassa/mnemonic"
)

const DerivationPath = "m/100H/10H"
const KeyStoreKey = "ethereum_private_key"

type Migration struct{}

func (m Migration) Up(mnemonic mnemonic.Mnemonic, keys map[string]string) (map[string]string, error) {

	//Get the coin seed
	seed, err := mnemonic.NewSeed("coins")
	if err != nil {
		return keys, err
	}

	//Create master key from coin seed
	masterKey, err := bip32.NewMasterKey(seed)
	if err != nil {
		return keys, err
	}

	//Derive the ethereum private key
	k, err := bip32.Derive(DerivationPath, masterKey)
	if err != nil {
		return keys, err
	}

	//private key
	privateKey := hex.EncodeToString(k.Key)

	//Exit with error if value is not like we expect it
	if value, exist := keys[KeyStoreKey]; exist && value != privateKey {
		return keys, errors.New("private key already exist BUT does not match derived private key")
	}

	keys[KeyStoreKey] = privateKey

	return keys, nil
}

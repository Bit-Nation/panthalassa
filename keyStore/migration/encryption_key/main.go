package encryption_key

import (
	"encoding/hex"
	"errors"
	"github.com/Bit-Nation/panthalassa/mnemonic"
)

const BIP39Password = "encryption_key"

type Migration struct{}

func (m Migration) Up(mnemonic mnemonic.Mnemonic, keys map[string]string) (map[string]string, error) {

	//Derive key's
	key, err := mnemonic.NewSeed(BIP39Password)
	if err != nil {
		return keys, err
	}
	key = key[:32]

	//check derived key with already existing key
	if existingKey, exist := keys[BIP39Password]; exist && existingKey != hex.EncodeToString(key) {
		return keys, errors.New("migration - derived key miss match with existing key")
	}

	keys[BIP39Password] = hex.EncodeToString(key)

	return keys, nil

}

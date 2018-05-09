package identity_ed_25519

import (
	"bytes"
	"encoding/hex"
	"errors"

	mnemonic "github.com/Bit-Nation/panthalassa/mnemonic"
	bip39 "github.com/tyler-smith/go-bip39"
	ed25519 "golang.org/x/crypto/ed25519"
)

const Bip39Password = "identity"
const Ed25519PrivateKey = "ed_25519_private_key"
const Ed25519PublicKey = "ed_25519_public_key"

type Migration struct{}

func (m Migration) Up(mnemonic mnemonic.Mnemonic, keys map[string]string) (map[string]string, error) {
	seed := bip39.NewSeed(mnemonic.String(), Bip39Password)

	//Create ed25519 key pair's
	edPub, edPriv, err := ed25519.GenerateKey(bytes.NewReader(seed))
	if err != nil {
		return keys, err
	}

	//Set private key
	if value, exist := keys[Ed25519PrivateKey]; exist && value != hex.EncodeToString(edPriv) {
		return keys, errors.New("migration - ed 25519 derivation miss match")
	}
	keys[Ed25519PrivateKey] = hex.EncodeToString(edPriv)

	//Set public key
	if value, exist := keys[Ed25519PublicKey]; exist && value != hex.EncodeToString(edPub) {
		return keys, errors.New("migration - ed 25519 derivation miss match")
	}
	keys[Ed25519PublicKey] = hex.EncodeToString(edPub)

	return keys, nil
}

package chat

import (
	"bytes"
	"encoding/hex"
	"errors"

	mnemonic "github.com/Bit-Nation/panthalassa/mnemonic"
	x3dh "github.com/Bit-Nation/x3dh"
)

type Migration struct{}

const MigrationPrivPrefix = "chat_identity_curve25519_private_key"
const MigrationPubPrefix = "chat_identity_curve25519_public_key"

func (m Migration) Up(mne mnemonic.Mnemonic, keys map[string]string) (map[string]string, error) {

	// get seed for signal private key
	s, err := mne.NewSeed("signal")
	if err != nil {
		return nil, err
	}

	// create curve and use the seed as random source
	c := x3dh.NewCurve25519(bytes.NewReader(s))
	keyPair, err := c.GenerateKeyPair()
	if err != nil {
		return nil, err
	}

	// validate if private key is still the same if already present
	currentPriv, exist := keys[MigrationPrivPrefix]
	if currentPriv != hex.EncodeToString(keyPair.PrivateKey[:]) && exist {
		return nil, errors.New("migration (chat) derivation miss match of private key")
	}

	// validate if public key is still the same if already present
	currentPub, exist := keys[MigrationPubPrefix]
	if currentPub != hex.EncodeToString(keyPair.PublicKey[:]) && exist {
		return nil, errors.New("migration (chat) derivation miss match of public key")
	}

	keys[MigrationPrivPrefix] = hex.EncodeToString(keyPair.PrivateKey[:])
	keys[MigrationPubPrefix] = hex.EncodeToString(keyPair.PublicKey[:])

	return keys, nil
}

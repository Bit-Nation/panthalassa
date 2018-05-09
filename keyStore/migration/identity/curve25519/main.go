package curve25519

import (
	"errors"

	"encoding/hex"
	identity "github.com/Bit-Nation/panthalassa/keyStore/migration/identity"
	mnemonic "github.com/Bit-Nation/panthalassa/mnemonic"
	extra "github.com/agl/ed25519/extra25519"
)

type Migration struct{}

func (m *Migration) Up(mnemonic mnemonic.Mnemonic, keys map[string]string) (map[string]string, error) {

	err := errors.New("you need to run the ed25519 migration's first")

	if _, exist := keys[identity.Ed25519PrivateKey]; !exist {
		return keys, err
	}

	if _, exist := keys[identity.Ed25519PublicKey]; !exist {
		return keys, err
	}

	//Get private key
	ed25519Priv, err := hex.DecodeString(keys[identity.Ed25519PrivateKey])
	if err != nil {
		return keys, err
	}
	if len(ed25519Priv) != 64 {
		return keys, errors.New("decoded private key MUST be 64 bytes long")
	}
	ed25519PrivByte := [64]byte{}
	copy(ed25519PrivByte[:], ed25519Priv)

	//Get public key
	ed25519Pub, err := hex.DecodeString(keys[identity.Ed25519PublicKey])
	if err != nil {
		return keys, err
	}
	if len(ed25519Priv) != 64 {
		return keys, errors.New("decoded private key MUST be 32 bytes long")
	}
	ed25519PubByte := [32]byte{}
	copy(ed25519PubByte[:], ed25519Pub)

	//private key
	privKey := [32]byte{}
	extra.PrivateKeyToCurve25519(&privKey, &ed25519PrivByte)
	privKeyStr := hex.EncodeToString(privKey[:])
	if value, exist := keys[identity.Curve25519PrivateKey]; exist && value != privKeyStr {
		return keys, errors.New("migration - curve25519 public key derivation miss match")
	}
	keys[identity.Curve25519PrivateKey] = privKeyStr

	//public key
	pubKey := [32]byte{}
	if success := extra.PublicKeyToCurve25519(&pubKey, &ed25519PubByte); !success {
		return keys, errors.New("migration - failed to convert ed25519 pub key to curve25519 pub key")
	}
	pubKeyString := hex.EncodeToString(pubKey[:])
	keys[identity.Curve25519PublicKey] = pubKeyString

	return keys, nil

}

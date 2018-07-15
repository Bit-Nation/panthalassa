package chat

import (
	"errors"

	x3dh "github.com/Bit-Nation/x3dh"
	dr "github.com/tiabc/doubleratchet"
	ed25519 "golang.org/x/crypto/ed25519"
)

type drDhPair struct {
	x3dhPair x3dh.KeyPair
}

func (p *drDhPair) PrivateKey() dr.Key {
	var k dr.Key
	copy(k[:], p.x3dhPair.PrivateKey[:])
	return k
}

func (p *drDhPair) PublicKey() dr.Key {
	var k dr.Key
	copy(k[:], p.x3dhPair.PublicKey[:])
	return k
}

// update the local signed pre key for a given id
func (c *Chat) refreshSignedPreKey(idPubKey ed25519.PublicKey) error {

	// fetch signed pre key bundle
	signedPreKey, err := c.backend.FetchSignedPreKey(idPubKey)
	if err != nil {
		return err
	}

	// verify signature of signed pre key bundle
	validSig, err := signedPreKey.VerifySignature(idPubKey)
	if err != nil {
		return err
	}
	if !validSig {
		return errors.New("signed pre key signature is invalid")
	}

	// check if signed pre key didn't expire
	expired := signedPreKey.OlderThan(SignedPreKeyValidTimeFrame)
	if expired {
		return errors.New("signed pre key expired")
	}

	return c.userStorage.PutSignedPreKey(idPubKey, signedPreKey)

}

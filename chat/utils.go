package chat

import (
	"errors"

	ed25519 "golang.org/x/crypto/ed25519"
)

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
	expired := signedPreKey.OlderThen(SignedPreKeyValidTimeFrame)
	if expired {
		return errors.New("signed pre key expired")
	}

	return c.userStorage.PutSignedPreKey(idPubKey, signedPreKey)

}

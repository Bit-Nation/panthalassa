package chat

import (
	"bytes"
	"errors"

	bpb "github.com/Bit-Nation/protobuffers"
	x3dh "github.com/Bit-Nation/x3dh"
	mh "github.com/multiformats/go-multihash"
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

// generate shared secret id
func sharedSecretID(sender, receiver ed25519.PublicKey, sharedSecretID []byte) (mh.Multihash, error) {
	b := bytes.NewBuffer(sender)
	if _, err := b.Write(receiver); err != nil {
		return nil, err
	}
	if _, err := b.Write(sharedSecretID); err != nil {
		return nil, err
	}
	return mh.Sum(b.Bytes(), mh.SHA3_256, -1)
}

// create identification
func sharedSecretInitID(sender, receiver ed25519.PublicKey, msg bpb.ChatMessage) (mh.Multihash, error) {
	if len(sender) != 32 {
		return nil, errors.New("sender must be 32 bytes long")
	}
	if len(receiver) != 32 {
		return nil, errors.New("receiver must be 32 bytes long")
	}
	b := bytes.Buffer{}
	if _, err := b.Write(sender); err != nil {
		return nil, err
	}
	if _, err := b.Write(receiver); err != nil {
		return nil, err
	}
	if _, err := b.Write(msg.SenderChatIDKey); err != nil {
		return nil, err
	}
	if _, err := b.Write(msg.SignedPreKey); err != nil {
		return nil, err
	}
	if _, err := b.Write(msg.OneTimePreKey); err != nil {
		return nil, err
	}
	return mh.Sum(b.Bytes(), mh.SHA3_256, -1)
}

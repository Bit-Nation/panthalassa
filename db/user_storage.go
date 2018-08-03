package db

import (
	preKey "github.com/Bit-Nation/panthalassa/chat/prekey"
	ed25519 "golang.org/x/crypto/ed25519"
)

// user storage store meta data about users
type UserStorage interface {
	// don't forget to verify the signature when implementing this
	GetSignedPreKey(idKey ed25519.PublicKey) (preKey.PreKey, error)
	HasSignedPreKey(idKey ed25519.PublicKey) (bool, error)
	PutSignedPreKey(idKey ed25519.PublicKey, key preKey.PreKey) error
}

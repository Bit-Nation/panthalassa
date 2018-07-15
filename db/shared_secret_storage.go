package db

import (
	"time"

	x3dh "github.com/Bit-Nation/x3dh"
	ed25519 "golang.org/x/crypto/ed25519"
)

type SharedSecret struct {
	X3dhSS                x3dh.SharedSecret
	Accepted              bool
	CreatedAt             time.Time
	DestroyAt             time.Time
	EphemeralKey          x3dh.PublicKey
	EphemeralKeySignature []byte
	UsedSignedPreKey      x3dh.PublicKey
	UsedOneTimePreKey     *x3dh.PublicKey
}

type SharedSecretStorage interface {
	HasAny(key ed25519.PublicKey) (bool, error)
	// must return an error if no shared secret found
	GetYoungest(key ed25519.PublicKey) (SharedSecret, error)
	Put(key ed25519.PublicKey, proto x3dh.InitializedProtocol) error
}

package db

import (
	x3dh "github.com/Bit-Nation/x3dh"
	ed25519 "golang.org/x/crypto/ed25519"
)

type OneTimePreKeyStorage interface {
	Cut(pubKey ed25519.PublicKey) (*x3dh.PrivateKey, error)
}

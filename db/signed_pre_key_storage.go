package db

import (
	x3dh "github.com/Bit-Nation/x3dh"
)

type SignedPreKeyStorage interface {
	// check if there is a active signed pre key
	HasActive() (bool, error)
	// get the current active pre key
	GetActive() (x3dh.KeyPair, error)
	// persist the signed pre key
	// @todo don't forget to give the option to register a listener in the put function
	Put(signedPreKey x3dh.KeyPair) error
}

package bip32

import (
	bip32 "github.com/tyler-smith/go-bip32"
)

type Key = bip32.Key

//Create new MasterKey as specified in bip32
//where child key's can be derived
func NewMasterKey(seed []byte) (*bip32.Key, error) {
	return bip32.NewMasterKey(seed)
}

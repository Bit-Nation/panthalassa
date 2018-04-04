package bip32

import (
	b32 "github.com/tyler-smith/go-bip32"
)

type Key = b32.Key

func NewMasterKey(seed []byte) (*b32.Key, error) {
	return b32.NewMasterKey(seed)
}

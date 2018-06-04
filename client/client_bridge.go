package client

import (
	x3dh "github.com/Bit-Nation/x3dh"
)

type Client interface {
	FetchSignedPreKey() x3dh.KeyPair
}

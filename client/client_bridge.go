package client

import (
	x3dh "github.com/Bit-Nation/x3dh"
)

type Bridge interface {
	// will fetch the signed pre key
	// the signed pre key will be encrypted. So make sure to decrypt it.
	FetchSignedPreKey() (string, error)
}

type Client struct {
	bridge Bridge
}

func (c *Client) FetchSignedPreKey() x3dh.KeyPair {
	return x3dh.KeyPair{}
}

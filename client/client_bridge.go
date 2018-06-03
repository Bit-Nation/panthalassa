package client

import (
	x3dh "github.com/Bit-Nation/x3dh"
)

type ClientBridge interface {
	// will fetch the signed pre key
	// the signed pre key will be encrypted. So make sure to decrypt it.
	FetchSignedPreKey() (string, error)
}

type Client struct {
	bridge ClientBridge
}

func (c *Client) FetchSignedPreKey() x3dh.KeyPair {
	return x3dh.KeyPair{}
}

func (c *Client) InitializedChat() error {

}

package mesh

import (
	"crypto/rand"
	"testing"

	lp2pCrypto "github.com/libp2p/go-libp2p-crypto"
	require "github.com/stretchr/testify/require"
)

func TestNewPrivateKey(t *testing.T) {

	priv, pubKey, err := lp2pCrypto.GenerateEd25519Key(rand.Reader)
	require.Nil(t, err)

	network, _, err := New(priv, "-")
	hostPubKey, err := network.host.ID().ExtractPublicKey()
	require.Nil(t, err)

	require.True(t, pubKey.Equals(hostPubKey))
}

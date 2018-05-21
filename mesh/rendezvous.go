package mesh

import (
	cid "github.com/ipfs/go-cid"
	lp2pCrypto "github.com/libp2p/go-libp2p-crypto"
	mh "github.com/multiformats/go-multihash"
)

const (
	pangea = "pangea"
)

type RendezvousKey struct {
	// RK
	key []byte
}

func (r *RendezvousKey) Derive(key string) (*cid.Cid, error) {
	final := append(r.key, []byte(key)...)
	return cid.NewPrefixV1(cid.Raw, mh.SHA3_256).Sum(final)
}

// generate profile key dht key for given public key
// RK || "profile" || pubKey
func (r RendezvousKey) Profile(pubKey lp2pCrypto.PubKey) (*cid.Cid, error) {
	rawPubKey, err := pubKey.Bytes()
	if err != nil {
		return nil, err
	}

	key := append(r.key, []byte("profile")...)
	return cid.NewPrefixV1(cid.Raw, mh.SHA3_256).Sum(append(key, rawPubKey...))
}

func NewRendezvousKey(seed string) RendezvousKey {
	return RendezvousKey{
		key: []byte(seed),
	}
}

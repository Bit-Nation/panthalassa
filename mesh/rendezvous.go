package mesh

import (
	cid "github.com/ipfs/go-cid"
	mh "github.com/multiformats/go-multihash"
)

const (
	pangea = "pangea"
)

type RendezvousKey struct {
	key []byte
}

func (r *RendezvousKey) Derive(key string) (*cid.Cid, error) {
	final := append(r.key, []byte(key)...)
	return cid.NewPrefixV1(cid.Raw, mh.SHA3_256).Sum(final)
}

func NewRendezvousKey(seed string) RendezvousKey {
	return RendezvousKey{
		key: []byte(seed),
	}
}

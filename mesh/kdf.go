package mesh

import (
	mh "github.com/multiformats/go-multihash"
	cid "github.com/ipfs/go-cid"
)

//Generate the rendezvous content id
//@todo chose something that involves time or a version on pangea
func rendezvousKey(seed string) (*cid.Cid, error) {

	return cid.NewPrefixV1(cid.Raw, mh.SHA3_256).Sum([]byte(seed))

}

package mesh

import (
	"gx/ipfs/QmZyZDi491cCNTLfAhwcaDii2Kg4pwKRkhqQzURGDvY6ua/go-multihash"
	"gx/ipfs/QmcZfnkapfECQGcLZaf9B79NRg7cRa9EnZh4LSbkCzwNvY/go-cid"
)

//Generate the rendezvous content id
//@todo chose something that involves time or a version on pangea
func rendezvousKey(seed string) (*cid.Cid, error) {

	return cid.NewPrefixV1(cid.Raw, multihash.SHA3_256).Sum([]byte(seed))

}

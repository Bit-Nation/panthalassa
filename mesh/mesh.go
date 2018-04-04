package mesh

import (
	"context"
	cid "github.com/ipfs/go-cid"
	datastore "github.com/ipfs/go-datastore"
	floodsub "github.com/libp2p/go-floodsub"
	libp2p "github.com/libp2p/go-libp2p"
	host "github.com/libp2p/go-libp2p-host"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	pstore "github.com/libp2p/go-libp2p-peerstore"
	ma "github.com/multiformats/go-multiaddr"
)

var bootstrapPeers = []string{
	"/ip4/104.131.131.82/tcp/4001/ipfs/QmaCpDMGvV2BGHeYERUEnRQAwe3N8SzbUtfsmvsqQLuvuJ",
	"/ip4/104.236.76.40/tcp/4001/ipfs/QmSoLV4Bbm51jM9C4gDYZQ9Cy3U6aXMJDAbzgu2fzaDs64",
	"/ip4/104.236.176.52/tcp/4001/ipfs/QmSoLnSGccFuZQJzRadHn95W2CrSFmZuTdDWP8HXaHca9z",
	"/ip4/104.236.179.241/tcp/4001/ipfs/QmSoLPppuBtQSGwKDZT2M73ULpjvfd3aZ6ha4oFGL1KrGM",
	"/ip4/162.243.248.213/tcp/4001/ipfs/QmSoLueR4xBeUbY9WZ9xGUUxunbKWcrNFTDAadQJmocnWm",
	"/ip4/128.199.219.111/tcp/4001/ipfs/QmSoLSafTMBsPKadTEgaXctDQVcqN88CNLHXMkTNwMKPnu",
	"/ip4/178.62.158.247/tcp/4001/ipfs/QmSoLer265NRgSp2LA3dPaeykiS1J6DifTC88f5uVQKNAd",
	"/ip4/178.62.61.185/tcp/4001/ipfs/QmSoLMeWqB7YGVLJN3pNLQpmmEk35v6wYtsMGLzSr5QBU3",
	"/ip4/104.236.151.122/tcp/4001/ipfs/QmSoLju6m7xTh3DuokvT3886QRYqxAzb1kShaanJgW36yx",
}

func meshConfig(cfg *libp2p.Config) error {
	// Create a multiaddress that listens on a random port on all interfaces
	addr, err := ma.NewMultiaddr("/ip4/0.0.0.0/tcp/0")
	if err != nil {
		return err
	}

	cfg.ListenAddrs = []ma.Multiaddr{addr}
	cfg.Peerstore = pstore.NewPeerstore()
	cfg.Muxer = libp2p.DefaultMuxer()
	return nil
}

type Mesh struct {
	dht           *dht.IpfsDHT
	host          host.Host
	started       bool
	ctx           context.Context
	floodSub      *floodsub.PubSub
	rendezvousKey *cid.Cid
	close         *chan struct{}
}

//Create a new instance of the mesh network
func NewMesh(rendezvousSeed string) (Mesh, error) {

	//Create rendezvous key
	rK, err := rendezvousKey(rendezvousSeed)

	if err != nil {
		return Mesh{}, err
	}

	//Mesh network instance
	m := Mesh{
		rendezvousKey: rK,
	}

	//Context
	m.ctx = context.Background()

	//Create host
	h, err := libp2p.New(m.ctx, libp2p.Defaults)
	m.host = h

	//Return on host error
	if err != nil {
		return Mesh{}, nil
	}

	//Create floodsub
	floodSub, err := floodsub.NewFloodSub(m.ctx, h)

	if err != nil {
		return Mesh{}, nil
	}

	m.floodSub = floodSub

	//Create close chan
	c := make(chan struct{})
	m.close = &c

	//Create dht
	//@todo use real data store
	m.dht = dht.NewDHTClient(m.ctx, h, datastore.NewMapDatastore())

	return m, nil
}

//Stop the mesh network
func (m *Mesh) Stop() {

	*m.close <- struct{}{}

}

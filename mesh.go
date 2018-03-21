package panthalassa

import (
	"context"
	"fmt"
	"gx/ipfs/QmNh1kGFFdsPu79KNSaL4NUKUPb4Eiz4KHdMtFY6664RDp/go-libp2p"
	host "gx/ipfs/QmNmJZL7FQySMtE2BQuLMuZg2EB2CLEunJJUSVSc9YnnbV/go-libp2p-host"
	"gx/ipfs/QmPpegoMqhAEqjncrzArm7KVWAkCm78rqL2DPuNjhPrshg/go-datastore"
	"gx/ipfs/QmQViVWBHbU6HmYjXcdNq7tVASCNgdg64ZGcauuDkLCivW/go-ipfs-addr"
	"gx/ipfs/QmVSep2WwKcXxMonPASsAJ3nZVjfVMKgMcaSigxKnUWpJv/go-libp2p-kad-dht"
	"gx/ipfs/QmWWQ2Txc2c6tqjsBpzg5Ar652cHPGNsQQp2SejkNmkUMb/go-multiaddr"
	"gx/ipfs/QmXauCuJzmzapetmC6W4TuDJLL1yFFrVzSHoWv8YdbmnxH/go-libp2p-peerstore"
	pstore "gx/ipfs/QmXauCuJzmzapetmC6W4TuDJLL1yFFrVzSHoWv8YdbmnxH/go-libp2p-peerstore"
	"gx/ipfs/QmcZfnkapfECQGcLZaf9B79NRg7cRa9EnZh4LSbkCzwNvY/go-cid"
	floodsub "gx/ipfs/QmctbcXMMhxTjm5ybWpjMwDmabB39ANuhB5QNn8jpD4JTv/go-libp2p-floodsub"
	"time"
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
	addr, err := multiaddr.NewMultiaddr("/ip4/0.0.0.0/tcp/0")
	if err != nil {
		return err
	}

	cfg.ListenAddrs = []multiaddr.Multiaddr{addr}
	cfg.Peerstore = pstore.NewPeerstore()
	cfg.Muxer = libp2p.DefaultMuxer()
	return nil
}

type Mesh struct {
	dht           *dht.IpfsDHT
	host          host.Host
	logger        CliLogger
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
		logger:        NewCliLogger(),
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

//Initial start of the mesh network and connect to bootstrapping nodes
func (m *Mesh) Start(cb func(error)) {

	//Connect to bootstrapping nodes
	for _, addr := range bootstrapPeers {
		iAddr, err := ipfsaddr.ParseString(addr)

		if err != nil {
			cb(err)
			return
		}

		pInfo, err := peerstore.InfoFromP2pAddr(iAddr.Multiaddr())

		if err != nil {
			cb(err)
			return
		}

		if err := m.host.Connect(m.ctx, *pInfo); err != nil {
			cb(err)
			return
		}

		m.logger.Info(fmt.Sprintf("connected to peer: %s", pInfo.ID.String()))
	}

	//Announce to the network that we are a member of bitnation
	tCtx, _ := context.WithTimeout(m.ctx, time.Second*10)
	//@todo why do we need the content timeout?
	if err := m.dht.Provide(tCtx, m.rendezvousKey, true); err != nil {
		cb(err)
		return
	}

	//Continue searching for peer's
	go func() {
		for {

			m.logger.Info("Search for peer's")

			//Find other bitnation peer's
			peers, err := m.dht.FindProviders(m.ctx, m.rendezvousKey)

			m.logger.Info(fmt.Sprintf("Found: %d peer's", len(peers)))

			//Connect to discovered nodes
			for _, peer := range peers {

				if peer.ID == m.host.ID() {
					continue
				}

				//@todo maybe check here if already connected to peer
				if err := m.host.Connect(m.ctx, peer); err != nil {
					m.logger.Error(err.Error())
				}

				m.logger.Info(fmt.Sprintf("connected to peer: %s", peer.ID.String()))
			}

			if err != nil {
				panic(err)
			}

			time.Sleep(10 * time.Second)
		}
	}()

	cb(nil)

	//Wait for the close. Blocking here.
	<-*m.close

	//Close host
	if err := m.host.Close(); err != nil {
		cb(err)
		return
	}

	//Close DHT
	if err := m.dht.Close(); err != nil {
		cb(err)
		return
	}

}

//Stop the mesh network
func (m *Mesh) Stop() {

	*m.close <- struct{}{}

}

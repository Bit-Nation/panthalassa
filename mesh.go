package panthalassa

import (
	"context"
	"fmt"
	"gx/ipfs/QmNh1kGFFdsPu79KNSaL4NUKUPb4Eiz4KHdMtFY6664RDp/go-libp2p"
	host "gx/ipfs/QmNmJZL7FQySMtE2BQuLMuZg2EB2CLEunJJUSVSc9YnnbV/go-libp2p-host"
	"gx/ipfs/QmPpegoMqhAEqjncrzArm7KVWAkCm78rqL2DPuNjhPrshg/go-datastore"
	"gx/ipfs/QmQViVWBHbU6HmYjXcdNq7tVASCNgdg64ZGcauuDkLCivW/go-ipfs-addr"
	"gx/ipfs/QmVSep2WwKcXxMonPASsAJ3nZVjfVMKgMcaSigxKnUWpJv/go-libp2p-kad-dht"
	ma "gx/ipfs/QmWWQ2Txc2c6tqjsBpzg5Ar652cHPGNsQQp2SejkNmkUMb/go-multiaddr"
	"gx/ipfs/QmXauCuJzmzapetmC6W4TuDJLL1yFFrVzSHoWv8YdbmnxH/go-libp2p-peerstore"
	pstore "gx/ipfs/QmXauCuJzmzapetmC6W4TuDJLL1yFFrVzSHoWv8YdbmnxH/go-libp2p-peerstore"
	"time"
	floodsub "gx/ipfs/QmctbcXMMhxTjm5ybWpjMwDmabB39ANuhB5QNn8jpD4JTv/go-libp2p-floodsub"
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
	dht     *dht.IpfsDHT
	host    host.Host
	logger  CliLogger
	started bool
}

func NewMesh() Mesh {
	return Mesh{
		logger: NewCliLogger(),
	}
}

//Initial start of the mesh network
func (m *Mesh) Start() error {

	if m.started == true {
		return nil
	}

	ctx := context.Background()

	//Create host
	h, err := libp2p.New(ctx, meshConfig)

	//Exit on error
	if err != nil {
		return err
	}

	m.host = h

	for _, addr := range bootstrapPeers {
		
		iAddr, err := ipfsaddr.ParseString(addr)

		if err != nil {
			//@todo we might now want to return here since one of the bootstrapping nodes could be offline
			return err
		}

		pInfo, err := peerstore.InfoFromP2pAddr(iAddr.Multiaddr())

		if err != nil {
			//@todo we might now want to return here since one of the bootstrapping nodes could be offline
			return err
		}

		tCtx, _ := context.WithTimeout(ctx, time.Second*7)
		if err := h.Connect(tCtx, *pInfo); err != nil {
			//@todo we might now want to return here since one of the bootstrapping nodes could be offline
			return err
		}

		m.logger.Info(fmt.Sprintf("connected to peer: %s", pInfo.ID.String()))

	}
	
	m.dht = dht.NewDHTClient(ctx, h, datastore.NewMapDatastore())
	m.started = true

	return nil
}

//Shutdown the mesh network
func (m *Mesh) Stop() error {

	if m.started == false {
		m.logger.Info("Please start the mesh first")
		return nil
	}

	if e := m.host.Close(); e != nil {
		m.logger.Debug("Closed the host")
		return e
	}

	if e := m.dht.Close(); e != nil {
		m.logger.Debug("Closed the DHT")
		return e
	}

	m.started = false

	return nil

}

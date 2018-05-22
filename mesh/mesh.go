package mesh

import (
	"context"
	"fmt"
	"time"

	deviceApi "github.com/Bit-Nation/panthalassa/api/device"
	bootstrap "github.com/florianlenz/go-libp2p-bootstrap"
	ds "github.com/ipfs/go-datastore"
	log "github.com/ipfs/go-log"
	lp2p "github.com/libp2p/go-libp2p"
	lp2pCrypto "github.com/libp2p/go-libp2p-crypto"
	host "github.com/libp2p/go-libp2p-host"
	dht "github.com/libp2p/go-libp2p-kad-dht"
)

var bootstrapPeers = []string{
	"/ip4/104.131.131.82/tcp/4001/ipfs/QmaCpDMGvV2BGHeYERUEnRQAwe3N8SzbUtfsmvsqQLuvuJ",
	"/ip4/104.236.76.40/tcp/4001/ipfs/QmSoLV4Bbm51jM9C4gDYZQ9Cy3U6aXMJDAbzgu2fzaDs64",
	"/ip4/104.236.176.52/tcp/4001/ipfs/QmSoLnSGccFuZQJzRadHn95W2CrSFmZuTdDWP8HXaHca9z",
	"/ip4/104.236.179.241/tcp/4001/ipfs/QmSoLPppuBtQSGwKDZT2M73ULpjvfd3aZ6ha4oFGL1KrGM",
	"/ip4/162.243.248.213/tcp/4001/ipfs/QmSoLueR4xBeUbY9WZ9xGUUxunbKWcrNFTDAadQJmocnWm",
	"/ip4/128.199.219.111/tcp/4001/ipfs/QmSoLSafTMBsPKadTEgaXctDQVcqN88CNLHXMkTNwMKPnu",
}

var logger = log.Logger("mesh")

func New(meshPk lp2pCrypto.PrivKey, api *deviceApi.Api, rendezvousKey, signedProfile string) (*Network, <-chan error, error) {

	//Create host
	h, err := lp2p.New(context.Background(), func(cfg *lp2p.Config) error {

		cfg.PeerKey = meshPk
		cfg.DisableSecio = false

		return lp2p.Defaults(cfg)

	})
	if err != nil {
		return nil, nil, err
	}

	//Create rendezvous key
	rk := NewRendezvousKey(rendezvousKey)

	//error report channel
	errReport := make(chan error)

	//Bootstrapping
	b, err := bootstrap.New(h, bootstrap.Config{
		BootstrapPeers:    bootstrapPeers,
		MinPeers:          5,
		BootstrapInterval: time.Second * 5,
		HardBootstrap:     time.Second * 120,
	})

	//start bootstrapping
	go func(b *bootstrap.Bootstrap) {
		logger.Debug("Start bootstrapping")
		if err := b.Start(context.Background()); err != nil {
			logger.Error(err)
		}
		logger.Debug("Finished bootstrapping")
	}(b)

	//Create DHT and bootstrap it
	d := dht.NewDHT(context.Background(), h, ds.NewMapDatastore())
	go func(dht *dht.IpfsDHT, key RendezvousKey, h host.Host) {
		logger.Debug("Start DHT bootstrapping")
		if err := dht.Bootstrap(context.Background()); err != nil {
			errReport <- err
		}
		logger.Debug("Finished DHT bootstrapping")

	}(d, rk, h)

	// put profile in the dht
	pubKey, err := h.ID().ExtractPublicKey()
	if err != nil {
		return nil, nil, err
	}
	cid, err := rk.Profile(pubKey)
	if err != nil {
		return nil, nil, err
	}
	logger.Info(fmt.Sprintf("put my profile (%s) with key (%s) in the DHT", signedProfile, cid.String()))
	d.PutValue(context.Background(), cid.String(), []byte(signedProfile))

	return &Network{
		Host:      h,
		Bootstrap: b,
		Dht:       d,
	}, errReport, nil

}

type Network struct {
	Host          host.Host
	Bootstrap     *bootstrap.Bootstrap
	Dht           *dht.IpfsDHT
	RendezvousKey RendezvousKey
}

func (n *Network) Close() error {

	var err error

	if e := n.Host.Close(); e != nil {
		logger.Error(err)
		err = e
	}

	if e := n.Dht.Close(); e != nil {
		logger.Error(err)
		err = e
	}

	return err
}

//Bootstrap manual
func (n *Network) BootstrapManual(ctx context.Context) error {
	return n.Bootstrap.Bootstrap(ctx)
}

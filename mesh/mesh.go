package mesh

import (
	"context"
	"time"

	bootstrap "github.com/florianlenz/go-libp2p-bootstrap"
	log "github.com/ipfs/go-log"
	lp2p "github.com/libp2p/go-libp2p"
	lp2pCrypto "github.com/libp2p/go-libp2p-crypto"
	host "github.com/libp2p/go-libp2p-host"
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

func New(meshPk *lp2pCrypto.Ed25519PrivateKey) (*Network, error) {

	//Create host
	h, err := lp2p.New(context.Background(), func(cfg *lp2p.Config) error {

		cfg.PeerKey = meshPk
		cfg.DisableSecio = false

		return lp2p.Defaults(cfg)

	})
	if err != nil {
		return nil, err
	}

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

	return &Network{
		host: h,
	}, nil

}

type Network struct {
	host      host.Host
	bootstrap bootstrap.Bootstrap
}

func (n *Network) Close() error {

	var err error

	err = n.host.Close()

	return err
}

//Bootstrap manual
func (n *Network) BootstrapManual(ctx context.Context) error {
	return n.bootstrap.Bootstrap(ctx)
}

package p2p

import (
	"context"

	log "github.com/ipfs/go-log"
	lp2p "github.com/libp2p/go-libp2p"
	host "github.com/libp2p/go-libp2p-host"
)

var logger = log.Logger("network")

func New() (*Network, error) {

	//Create host
	h, err := lp2p.New(context.Background(), func(cfg *lp2p.Config) error {
		if err := lp2p.Defaults(cfg); err != nil {
			return err
		}
		cfg.DisableSecio = false
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &Network{
		Host: h,
	}, nil

}

type Network struct {
	Host host.Host
}

func (n *Network) Close() error {

	var err error

	if e := n.Host.Close(); e != nil {
		logger.Error(err)
		err = e
	}

	return err
}

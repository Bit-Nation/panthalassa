package p2p

import (
	"context"

	log "github.com/ipfs/go-log"
	lp2p "github.com/libp2p/go-libp2p"
	host "github.com/libp2p/go-libp2p-host"
	mplex "github.com/whyrusleeping/go-smux-multiplex"
	msmux "github.com/whyrusleeping/go-smux-multistream"
	yamux "github.com/whyrusleeping/go-smux-yamux"
)

var logger = log.Logger("network")

func New() (*Network, error) {

	//Create host
	h, err := lp2p.New(context.Background(), func(cfg *lp2p.Config) error {
		if err := lp2p.Defaults(cfg); err != nil {
			return err
		}
		cfg.DisableSecio = false

		// add muxer
		tpt := msmux.NewBlankTransport()
		tpt.AddTransport("/yamux/1.0.0", yamux.DefaultTransport)
		// @todo mplex is registered twice till this https://github.com/libp2p/go-libp2p/commit/77b7d8f06f6639fcff5414525257ec18808b6112 is released
		tpt.AddTransport("/mplex/6.3.0", mplex.DefaultTransport)
		tpt.AddTransport("/mplex/6.7.0", mplex.DefaultTransport)
		cfg.Muxer = tpt

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

package mesh

/*
@todo there is an problem with the bootstrapping bundle. Since we don't need the mesh network now we will comment it
import (
	"context"
	bootstrap "github.com/florianlenz/go-libp2p-bootstrap"
	cid "github.com/ipfs/go-cid"
	libp2p "github.com/libp2p/go-libp2p"
	host "github.com/libp2p/go-libp2p-host"
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

//configuration for mesh network
func meshConfig(cfg *libp2p.Config) error {

	libp2p.Defaults(cfg)
	return nil
}

type Mesh struct {
	host          host.Host
	started       bool
	ctx           context.Context
	rendezvousKey *cid.Cid
	bootstrap     *bootstrap.Bootstrap
	close         *chan struct{}
}

//Create a new instance of the mesh network
func NewMesh(rendezvousSeed string) (*Mesh, error) {

	//Create rendezvous key
	rK, err := rendezvousKey(rendezvousSeed)
	if err != nil {
		return &Mesh{}, err
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
	if err != nil {
		return &Mesh{}, nil
	}

	//Create close chan
	c := make(chan struct{})
	m.close = &c

	//Add bootstrap to host
	cfg := bootstrap.Config{
		BootstrapPeers:    bootstrapPeers,
		MinPeers:          6,
		BootstrapInterval: time.Second * 5,
		HardBootstrap:     time.Second * 100,
	}
	err, bs := bootstrap.NewBootstrap(h, cfg)
	if err != nil {
		return &Mesh{}, err
	}
	m.bootstrap = bs

	return &m, nil
}

//Start mesh network
func (m *Mesh) Start(cb func(err error)) {

	//Bootstrap start
	if err := m.bootstrap.Start(); err != nil {
		cb(err)
		return
	}

	cb(nil)

	//Hang till closed
	<-*m.close
}

//Stop the mesh network
func (m *Mesh) Stop() {

	//Stop mesh network
	m.bootstrap.Stop()

	*m.close <- struct{}{}

}
*/

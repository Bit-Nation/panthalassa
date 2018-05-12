package mesh

import (
	"context"

	libp2p "gx/ipfs/QmNh1kGFFdsPu79KNSaL4NUKUPb4Eiz4KHdMtFY6664RDp/go-libp2p"
	host "gx/ipfs/QmNmJZL7FQySMtE2BQuLMuZg2EB2CLEunJJUSVSc9YnnbV/go-libp2p-host"
)

//configuration for mesh network
func meshConfig(cfg *libp2p.Config) error {
	return libp2p.Defaults(cfg)
}

type Mesh struct {
	host  host.Host
	close chan struct{}
}

//Create a new instance of the mesh network
func NewMesh(rendezvousSeed string) (*Mesh, error) {

	m := Mesh{
		close: make(chan struct{}, 1),
	}

	//Create host
	h, err := libp2p.New(context.Background(), libp2p.Defaults)
	m.host = h
	if err != nil {
		return &Mesh{}, nil
	}

	return &m, nil
}

//Start mesh network
func (m *Mesh) Start(cb func(err error)) {

	cb(nil)

	//Hang till closed
	<-m.close
}

//Stop the mesh network
func (m *Mesh) Close() error {

	//Stop mesh network
	if err := m.host.Close(); err != nil {
		m.close <- struct{}{}
		return err
	}

	m.close <- struct{}{}
	return nil
}

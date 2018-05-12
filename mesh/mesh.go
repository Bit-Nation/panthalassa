package mesh

import (
	"context"

	pangeaProtocol "github.com/Bit-Nation/panthalassa_new/protocols/pangea"
	libp2p "github.com/libp2p/go-libp2p"
	host "github.com/libp2p/go-libp2p-host"
	net "github.com/libp2p/go-libp2p-net"
)

//configuration for mesh network
func meshConfig(cfg *libp2p.Config) error {
	return libp2p.Defaults(cfg)
}

type Mesh struct {
	host      host.Host
	close     chan struct{}
	protocols []Protocol
}

//Create a new instance of the mesh network
func NewMesh() (*Mesh, error) {

	m := Mesh{
		close: make(chan struct{}, 1),
		protocols: []Protocol{
			&pangeaProtocol.Protocol{},
		},
	}

	//Create host
	h, err := libp2p.New(context.Background(), libp2p.Defaults)
	m.host = h
	if err != nil {
		return nil, err
	}

	//Register protocols
	for _, p := range m.protocols {
		p.Register(h)
	}

	//Register network listener
	h.Network().Notify(&NotifyBundle{
		OpenedStreamF: func(_ net.Network, stream net.Stream) {

		},
	})

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

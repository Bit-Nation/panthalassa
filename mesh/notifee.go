package mesh

//This was taken from the github.com/libp2p/go-libp2p-net repo as the current GX version doesn't contain this file

import (
	net "github.com/libp2p/go-libp2p-net"
	ma "github.com/multiformats/go-multiaddr"
)

// NotifyBundle implements Notifiee by calling any of the functions set on it,
// and nop'ing if they are unset. This is the easy way to register for
// notifications.
type NotifyBundle struct {
	ListenF      func(net.Network, ma.Multiaddr)
	ListenCloseF func(net.Network, ma.Multiaddr)

	ConnectedF    func(net.Network, net.Conn)
	DisconnectedF func(net.Network, net.Conn)

	OpenedStreamF func(net.Network, net.Stream)
	ClosedStreamF func(net.Network, net.Stream)
}

func (nb *NotifyBundle) Listen(n net.Network, a ma.Multiaddr) {
	if nb.ListenF != nil {
		nb.ListenF(n, a)
	}
}

func (nb *NotifyBundle) ListenClose(n net.Network, a ma.Multiaddr) {
	if nb.ListenCloseF != nil {
		nb.ListenCloseF(n, a)
	}
}

func (nb *NotifyBundle) Connected(n net.Network, c net.Conn) {
	if nb.ConnectedF != nil {
		nb.ConnectedF(n, c)
	}
}

func (nb *NotifyBundle) Disconnected(n net.Network, c net.Conn) {
	if nb.DisconnectedF != nil {
		nb.DisconnectedF(n, c)
	}
}

func (nb *NotifyBundle) OpenedStream(n net.Network, s net.Stream) {
	if nb.OpenedStreamF != nil {
		nb.OpenedStreamF(n, s)
	}
}

func (nb *NotifyBundle) ClosedStream(n net.Network, s net.Stream) {
	if nb.ClosedStreamF != nil {
		nb.ClosedStreamF(n, s)
	}
}

package mesh

import (
	host "github.com/libp2p/go-libp2p-host"
)

type Protocol interface {
	Register(h host.Host)
	UnRegister(h host.Host)
}

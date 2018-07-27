package p2p

import (
	"context"
	"log"

	formatter "github.com/ipfs/go-log/writer"
	ps "github.com/libp2p/go-libp2p-peerstore"
	goLog "github.com/whyrusleeping/go-logging"
)

const logProtocol = "/pangea/logger/1.0.0"

func (n *Network) ConnectLogger(pInfo ps.PeerInfo) error {

	// connect to host first
	err := n.Host.Connect(context.Background(), pInfo)
	if err != nil {
		return err
	}

	// dial to protocol
	str, err := n.Host.NewStream(context.Background(), pInfo.ID, logProtocol)
	if err != nil {
		return err
	}

	// set formatter
	goLog.SetFormatter(&formatter.PoliteJSONFormatter{})

	// set log backend
	goLog.SetBackend(goLog.NewLogBackend(str, "", log.Ltime|log.Llongfile))

	return nil

}

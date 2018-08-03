package p2p

import (
	"context"

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
	formatter, err := goLog.NewStringFormatter("%{time} - %{shortfile} (%{level}): %{module} %{message}")
	if err != nil {
		return err
	}
	goLog.SetFormatter(formatter)

	// set log backend
	b := goLog.NewLogBackend(str, "", 0)
	b.Color = true
	goLog.SetBackend(b)

	return nil

}

package registry

import (
	"bufio"
	"fmt"

	dapp "github.com/Bit-Nation/panthalassa/dapp"
	net "github.com/libp2p/go-libp2p-net"
	mc "github.com/multiformats/go-multicodec/json"
)

// this stream handler is used for development purpose
// when we receive a DApp we will send it to the client
// the client will then decide what to do with it.
func (r *Registry) devStreamHandler(str net.Stream) {

	// exit if peer is not white listed for DApp development
	if !r.state.HasDAppDevPeer(str.Conn().RemotePeer()) {
		logger.Warning(fmt.Sprintf("peer: %s wasn't whitelisted for dapp development", str.Conn().RemotePeer().Pretty()))
		if err := str.Reset(); err != nil {
			logger.Error(err)
		}
		return
	}

	go func() {

		reader := bufio.NewReader(str)
		decoder := mc.Multicodec(true).Decoder(reader)

		// decode app from stream
		var app dapp.JsonRepresentation
		if err := decoder.Decode(&app); err != nil {
			logger.Error(err)
		}

		// add stream to registry so that we can
		// associate it with the DApp
		r.addDAppDevStream(app.SignaturePublicKey, str)

		//@todo send dapp to device

	}()

}

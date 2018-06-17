package registry

import (
	"bufio"
	"encoding/json"

	dapp "github.com/Bit-Nation/panthalassa/dapp"
	pb "github.com/Bit-Nation/panthalassa/dapp/registry/pb"
	net "github.com/libp2p/go-libp2p-net"
	protoMc "github.com/multiformats/go-multicodec/protobuf"
)

// this stream handler is used for development purpose
// when we receive a DApp we will send it to the client
// the client will then decide what to do with it.
func (r *Registry) devStreamHandler(str net.Stream) {

	go func() {

		reader := bufio.NewReader(str)
		decoder := protoMc.Multicodec(nil).Decoder(reader)

		for {

			// decode app from stream
			msg := pb.Message{}
			if err := decoder.Decode(&msg); err != nil {
				logger.Error(err)
				continue
			}

			if msg.Type != pb.Message_DApp {
				logger.Error("i can only handle DApps")
				continue
			}

			var app dapp.JsonRepresentation
			if err := json.Unmarshal(msg.DApp, &app); err != nil {
				logger.Error(err)
				continue
			}

			// add stream to registry so that we can
			// associate it with the DApp
			r.addDAppDevStream(app.SignaturePublicKey, str)

			valid, err := app.VerifySignature()
			if err != nil {
				logger.Error(err)
				continue
			}
			if !valid {
				logger.Error("Received invalid signature for DApp: ", app.Name)
				continue
			}

			// push received DApp upstream
			if err := r.client.HandleReceivedDApp(app); err != nil {
				logger.Error(err)
			}

		}

	}()

}

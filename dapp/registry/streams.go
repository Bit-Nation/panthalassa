package registry

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"io"

	dapp "github.com/Bit-Nation/panthalassa/dapp"
	net "github.com/libp2p/go-libp2p-net"
)

// this stream handler is used for development purpose
// when we receive a DApp we will send it to the client
// the client will then decide what to do with it.
func (r *Registry) devStreamHandler(str net.Stream) {

	go func() {

		reader := bufio.NewReader(str)

		for {

			// read DApp from stream
			jsonDAppBytes, err := reader.ReadBytes(0x0A)
			if err != nil {
				logger.Error(err)
				if err == io.EOF {
					str.Close()
				} else {
					str.Reset()
				}
				break
			}

			// decode base64 json
			rawJsonDApp, err := base64.StdEncoding.DecodeString(string(jsonDAppBytes))
			if err != nil {
				logger.Error(err)
				continue
			}

			// unmarshal DApp data
			rawDAppData := dapp.RawData{}
			if err := json.Unmarshal(rawJsonDApp, &rawDAppData); err != nil {
				logger.Error(err)
				continue
			}

			// parse json to DApp Data
			dAppData, err := dapp.ParseJsonToData(rawDAppData)
			if err != nil {
				logger.Error(err)
				continue
			}

			// add stream to registry so that we can
			// associate it with the DApp
			r.addDevStreamChan <- addDevStreamChanStr{
				signingKey: dAppData.UsedSigningKey,
				stream:     str,
			}

			valid, err := dAppData.VerifySignature()
			if err != nil {
				logger.Error(err)
				continue
			}
			if !valid {
				logger.Error("Received invalid signature for DApp: ", dAppData.Name)
				continue
			}

			// persist received app
			if err := r.dAppDB.SaveDApp(dAppData); err != nil {
				logger.Error(err)
			}

		}

	}()

}

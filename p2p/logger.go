package p2p

import (
	"context"
	"encoding/json"
	"io"
	"log"

	ps "github.com/libp2p/go-libp2p-peerstore"
	goLog "github.com/whyrusleeping/go-logging"
)

const logProtocol = "/pangea/logger/1.0.0"

// to this from the ipfs logger since GX caused strange re writes in the import path
type politeJSONFormatter struct{}

// Format encodes a logging.Record in JSON and writes it to Writer.
func (f *politeJSONFormatter) Format(calldepth int, r *goLog.Record, w io.Writer) error {
	entry := make(map[string]interface{})
	entry["id"] = r.Id
	entry["level"] = r.Level
	entry["time"] = r.Time
	entry["module"] = r.Module
	entry["message"] = r.Message()
	err := json.NewEncoder(w).Encode(entry)
	if err != nil {
		return err
	}

	w.Write([]byte{'\n'})
	return nil
}

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
	goLog.SetFormatter(&politeJSONFormatter{})

	// set log backend
	goLog.SetBackend(goLog.NewLogBackend(str, "", log.Ltime|log.Llongfile))

	return nil

}

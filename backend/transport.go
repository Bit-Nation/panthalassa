package backend

import bpb "github.com/Bit-Nation/protobuffers"

type Transport interface {
	// will be called by the backend to send a message
	Send(msg *bpb.BackendMessage) error
	// will return the next message from the transport
	NextMessage() (*bpb.BackendMessage, error)
	// close the transport
	Close() error
}

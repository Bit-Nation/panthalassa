package backend

import (
	bpb "github.com/Bit-Nation/protobuffers"
)

type testTransport struct {
	send func(msg *bpb.BackendMessage) error
	// this will be the callback that should be called on a message
	// from the transport
	nextMessage func() (*bpb.BackendMessage, error)
}

func (t *testTransport) Send(msg *bpb.BackendMessage) error {
	return t.send(msg)
}

func (t *testTransport) NextMessage() (*bpb.BackendMessage, error) {
	return t.nextMessage()
}

func (t *testTransport) Close() error {
	return nil
}

func (t *testTransport) Start() error {
	return nil
}

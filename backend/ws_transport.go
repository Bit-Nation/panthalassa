package backend

import (
	bpb "github.com/Bit-Nation/protobuffers"
)

type WSTransport struct {
}

func (t *WSTransport) Send(msg *bpb.BackendMessage) error {
	return nil
}

func (t *WSTransport) NextMessage() (*bpb.BackendMessage, error) {
	return nil, nil
}

func (t *WSTransport) Close() error {
	return nil
}

func (t *WSTransport) Start() error {
	return nil
}

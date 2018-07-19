package backend

import bpb "github.com/Bit-Nation/protobuffers"

type Transport interface {
	Send(msg *bpb.BackendMessage) error
	OnMessage(func(msg *bpb.BackendMessage) error)
	Close() error
}

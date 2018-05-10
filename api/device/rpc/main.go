package rpc

type JsonRPCCall interface {
	Type() string
	Data() (string, error)
	Valid() error
}

type JsonRPCResponse interface {
	Valid()
	Close()
}

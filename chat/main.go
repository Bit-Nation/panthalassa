package chat

import (
	host "github.com/libp2p/go-libp2p-host"
	ma "github.com/multiformats/go-multiaddr"
)

type Command struct {
	Type string `json:"type"`
	Data string `json:"data"`
}

type Chat struct {
	stream *Stream
}

const ProtocolID = "chat/1.0.0"

func New(h host.Host, chatNode ma.Multiaddr) (*Chat, error) {

	s, err := NewStream(h, chatNode)
	if err != nil {
		return nil, err
	}

	return &Chat{
		stream: s,
	}, nil

}

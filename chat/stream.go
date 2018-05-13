package chat

import (
	"bufio"
	"context"

	host "github.com/libp2p/go-libp2p-host"
	net "github.com/libp2p/go-libp2p-net"
	pstore "github.com/libp2p/go-libp2p-peerstore"
	ma "github.com/multiformats/go-multiaddr"
	json "github.com/multiformats/go-multicodec/json"
)

type Stream struct {
	Write chan<- Command
	Read  <-chan Command
	s     net.Stream
}

func (s *Stream) Close() error {
	return s.s.Close()
}

func NewStream(h host.Host, chatNode ma.Multiaddr) (*Stream, error) {

	writeChan := make(chan Command)
	readChan := make(chan Command)

	//Get peer info from multi address
	pi, err := pstore.InfoFromP2pAddr(chatNode)
	if err != nil {
		return nil, err
	}

	//Connect to chat node
	if err = h.Connect(context.Background(), *pi); err != nil {
		return nil, err
	}

	//Create stream to chat node
	s, err := h.NewStream(context.Background(), pi.ID, ProtocolID)
	if err != nil {
		return nil, err
	}

	//Write to stream
	go func() {

		write := bufio.NewWriter(s)
		//@todo check that the "false" does.
		writeEncoder := json.Multicodec(false).Encoder(write)

		for {
			select {
			case m := <-writeChan:
				if err = writeEncoder.Encode(m); err != nil {
					//@todo find better way then throwing
					panic(err)
				}
				write.Flush()

			}
		}
	}()

	//Read from stream
	go func() {

		read := bufio.NewReader(s)
		readDecoder := json.Multicodec(false).Decoder(read)

		for {
			var c Command
			if err := readDecoder.Decode(&c); err != nil {
				//@todo find better way then throwing
				panic(err)
			}
			readChan <- c
		}

	}()

	return &Stream{
		Write: writeChan,
		Read:  readChan,
		s:     s,
	}, nil

}

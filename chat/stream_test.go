package chat

import (
	"context"
	"fmt"
	"testing"

	"bufio"
	libp2p "github.com/libp2p/go-libp2p"
	net "github.com/libp2p/go-libp2p-net"
	ma "github.com/multiformats/go-multiaddr"
	json "github.com/multiformats/go-multicodec/json"
	require "github.com/stretchr/testify/require"
)

func TestStream(t *testing.T) {

	//Create chat node
	chatNode, err := libp2p.New(context.Background(), libp2p.Defaults)
	require.Nil(t, err)
	validMa, err := ma.NewMultiaddr(fmt.Sprintf("%s/ipfs/%s", chatNode.Addrs()[0], chatNode.ID().Pretty()))
	require.Nil(t, err)

	//Register chat protocol on chat node
	chatNode.SetStreamHandler(ProtocolID, func(stream net.Stream) {

		//Read / Write
		reader := bufio.NewReader(stream)
		writer := bufio.NewWriter(stream)

		//Encoder
		enc := json.Multicodec(false).Encoder(writer)
		dec := json.Multicodec(false).Decoder(reader)

		//Decoded command
		var c Command
		err := dec.Decode(&c)
		require.Nil(t, err)

		//Check if test data is correct
		require.Equal(t, c.Type, "GREET")
		require.Equal(t, c.Data, "3")

		//Send a command back
		enc.Encode(Command{
			Type: "GREET_BACK",
			Data: "6",
		})
		writer.Flush()

	})

	//My peer
	me, err := libp2p.New(context.Background(), libp2p.Defaults)
	require.Nil(t, err)

	//New stream to chat node
	s, err := NewStream(me, validMa)
	require.Nil(t, err)

	s.Write <- Command{
		Type: "GREET",
		Data: "3",
	}

	c := <-s.Read
	require.Equal(t, "GREET_BACK", c.Type)
	require.Equal(t, "6", c.Data)

}

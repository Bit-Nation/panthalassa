package ethws

import (
	"testing"
	"time"

	require "github.com/stretchr/testify/require"
)

func TestEthereumWebSocket(t *testing.T) {

	ws := New(Config{
		Retry: time.Second,
		WSUrl: "wss://mainnet.infura.io/_ws",
	})

	respChan, err := ws.SendRequest(Request{
		Method:  `eth_protocolVersion`,
		JsonRPC: rpcVersion,
	})

	require.Nil(t, err)

	resp := <-respChan
	require.Nil(t, resp.Error)

	require.Equal(t, "2.0", resp.JsonRPC)
	require.Equal(t, "0x3f", resp.Result)

}

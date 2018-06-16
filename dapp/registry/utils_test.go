package registry

import (
	"sync"
	"testing"

	net "github.com/libp2p/go-libp2p-net"
	swarm "github.com/libp2p/go-libp2p-swarm"
	"github.com/stretchr/testify/require"
)

// test add / get DApp dev stream
func TestAddGetDAppDevStream(t *testing.T) {

	reg := Registry{
		lock:           sync.Mutex{},
		dAppDevStreams: map[string]net.Stream{},
	}

	s := &swarm.Stream{}

	// add stream
	reg.addDAppDevStream([]byte("my_app_id"), s)

	// get stream
	exist, str := reg.getDAppDevStream([]byte("my_app_id"))
	require.Equal(t, s, str)
	require.True(t, exist)

}

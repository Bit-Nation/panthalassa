package registry

import (
	"sync"
	"testing"

	net "github.com/libp2p/go-libp2p-net"
	"github.com/stretchr/testify/require"
)

// test add / get DApp dev stream
func TestAddGetDAppDevStream(t *testing.T) {

	reg := Registry{
		lock:           sync.Mutex{},
		dAppDevStreams: map[string]net.Stream{},
	}

	s := &stream{}

	key := []byte("my_app_id")

	// test what happens if stream doesn't exist
	exist, str := reg.getDAppDevStream(key)
	require.Nil(t, str)
	require.False(t, exist)

	// add stream
	reg.addDAppDevStream(key, s)

	// get stream
	exist, str = reg.getDAppDevStream(key)
	require.Equal(t, s, str)
	require.True(t, exist)

}

package registry

import (
	"encoding/hex"

	net "github.com/libp2p/go-libp2p-net"
)

// add a stream that relates to DApp development
// to the registry. Use this to add a stream
// in a thread safe way
func (r *Registry) addDAppDevStream(key []byte, str net.Stream) {
	r.lock.Lock()
	defer r.lock.Unlock()
	r.dAppDevStreams[hex.EncodeToString(key)] = str
}

// get a stream that relates to DApp development
// from the registry. Use this for thread safety
func (r *Registry) getDAppDevStream(key []byte) (bool, net.Stream) {
	r.lock.Lock()
	defer r.lock.Unlock()
	exist, stream := r.dAppDevStreams[hex.EncodeToString(key)]
	return stream, exist
}

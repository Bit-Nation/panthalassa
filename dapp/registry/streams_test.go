package registry

import (
	"errors"
	"sync"
	"testing"
	"time"

	state "github.com/Bit-Nation/panthalassa/state"
	crypto "github.com/libp2p/go-libp2p-crypto"
	net "github.com/libp2p/go-libp2p-net"
	peer "github.com/libp2p/go-libp2p-peer"
	ma "github.com/multiformats/go-multiaddr"
	require "github.com/stretchr/testify/require"
)

// test stream implementation
type stream struct {
	net.Stream
	data                           []byte
	failRead, failWrite, failClose bool
	reset                          bool
	conn                           net.Conn
}

func (s *stream) Reset() error {
	s.reset = true
	return nil
}

func (s *stream) Close() error {
	return errors.New("Close is not implemented")
}

func (s *stream) SetDeadline(t time.Time) error {
	return errors.New("SetDeadline is not implemented")
}

func (s *stream) SetReadDeadline(t time.Time) error {
	return errors.New("SetReadDeadline is not implemented")
}

func (s *stream) SetWriteDeadline(t time.Time) error {
	return errors.New("SetWriteDeadline is not implemented")
}

func (s *stream) Write(b []byte) (int, error) {
	return 0, errors.New("write is not implemented")
}

func (s *stream) Read(b []byte) (int, error) {
	return 0, errors.New("read is not implemented")
}

func (s *stream) Conn() net.Conn {
	return s.conn
}

type conn struct {
	remotePeerId peer.ID
}

func (c *conn) Close() error {
	return errors.New("Close not implemented")
}

func (c *conn) NewStream() (net.Stream, error) {
	return nil, errors.New("NewStream not implemented")
}

func (c *conn) GetStreams() []net.Stream {
	panic("GetStreams not implemented")
}

// LocalMultiaddr is the Multiaddr on this side
func (c *conn) LocalMultiaddr() ma.Multiaddr {
	panic("LocalMultiaddr - not implemented")
}

// LocalPeer is the Peer on our side of the connection
func (c *conn) LocalPeer() peer.ID {
	panic("LocalPeer - not implemented")
}

// LocalPrivateKey is the private key of the peer on our side.
func (c *conn) LocalPrivateKey() crypto.PrivKey {
	panic("LocalPrivateKey - not implemented")
}

// RemoteMultiaddr is the Multiaddr on the remote side
func (c *conn) RemoteMultiaddr() ma.Multiaddr {
	panic("RemoteMultiaddr - not implemented")
}

// RemotePeer is the Peer on the remote side
func (c *conn) RemotePeer() peer.ID {
	return c.remotePeerId
}

// RemotePublicKey is the private key of the peer on our side.
func (c *conn) RemotePublicKey() crypto.PubKey {
	panic("not implemented")
	return nil
}

func TestDevStreamHandler(t *testing.T) {

	// app state
	s := state.New()

	// create DApp registry
	r := Registry{
		state:          s,
		lock:           sync.Mutex{},
		dAppDevStreams: map[string]net.Stream{},
	}

	// create fake connection
	c := conn{
		remotePeerId: "i_am_the_remote_peer_id",
	}

	// create fake stream
	testStream := stream{
		conn: &c,
	}

	// make sure the default value of reset is false
	require.False(t, testStream.reset)

	// now try to connect it should fail
	// since the peer is not whitelisted
	r.devStreamHandler(&testStream)

	// make sure reset is true since we reset
	// the stream when a peer is not allowed
	// to connect on that protocol
	require.True(t, testStream.reset)

}

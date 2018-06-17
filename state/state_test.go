package state

import (
	"testing"

	require "github.com/stretchr/testify/require"
)

func TestState_AddDAppDevPeer(t *testing.T) {

	s := New()
	s.AddDAppDevPeer("fake-peer-id")

	// The value is "2w4H9nYaM9EYbazsd" since
	// peer id is a struct that implements "String()"
	// with a encoding function
	require.Equal(t, "2w4H9nYaM9EYbazsd", s.dappDevPeers[0])

}

func TestState_HasDAppDevPeer(t *testing.T) {

	s := New()
	s.AddDAppDevPeer("fake-peer-id")

	require.True(t, s.HasDAppDevPeer("fake-peer-id"))

}

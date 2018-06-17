package state

import (
	"github.com/libp2p/go-libp2p-peer"
	"sync"
)

type State struct {
	lock         sync.Mutex
	dappDevPeers []string
}

func New() *State {
	return &State{
		lock:         sync.Mutex{},
		dappDevPeers: []string{},
	}
}

// add a peer that is whitelisted
// for DApp Development
func (s *State) AddDAppDevPeer(pi peer.ID) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.dappDevPeers = append(s.dappDevPeers, pi.Pretty())
}

// remove a peer from the DApp Development list
func (s *State) HasDAppDevPeer(pi peer.ID) bool {
	s.lock.Lock()
	defer s.lock.Unlock()
	for _, peerId := range s.dappDevPeers {
		if peerId == pi.Pretty() {
			return true
		}
	}
	return false
}

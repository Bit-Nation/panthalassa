package device_api

import (
	"errors"
	"fmt"
	"math/rand"
	"sync"
)

type State struct {
	requests map[uint32]chan Response
	m        *sync.Mutex
}

func newState() *State {

	return &State{
		requests: make(map[uint32]chan Response),
		m:        &sync.Mutex{},
	}

}

func (s *State) Add(respChan chan Response) uint32 {

	s.m.Lock()
	var key uint32
	//@todo we should have a backup break for the for loop
	for {
		key = rand.Uint32()
		if _, exist := s.requests[key]; !exist {
			break
		}
	}
	s.requests[key] = respChan
	s.m.Unlock()

	return key

}

//Return's the channel an removes it from the state map
func (s *State) Cut(index uint32) (chan Response, error) {

	s.m.Lock()
	respChan, exist := s.requests[index]
	if !exist {
		return nil, errors.New(fmt.Sprintf("a request channel for id (%d) does not exist", index))
	}
	delete(s.requests, index)
	s.m.Unlock()

	return respChan, nil

}

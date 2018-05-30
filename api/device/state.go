package device_api

import (
	"errors"
	"fmt"
	"sync"

	uuid "github.com/satori/go.uuid"
)

type State struct {
	requests map[string]chan Response
	m        *sync.Mutex
}

func newState() *State {

	return &State{
		requests: make(map[string]chan Response),
		m:        &sync.Mutex{},
	}

}

func (s *State) Add(respChan chan Response) (string, error) {

	s.m.Lock()
	var id string
	//@todo we should have a backup break for the for loop
	for {
		uid, err := uuid.NewV4()
		if err != nil {
			return "", err
		}
		id = uid.String()

		if _, exist := s.requests[id]; !exist {
			break
		}
	}
	logger.Debug(fmt.Sprintf("added response channel with id: %s", id))
	s.requests[id] = respChan
	s.m.Unlock()

	return id, nil

}

//Return's the channel an removes it from the state map
func (s *State) Cut(id string) (chan Response, error) {

	s.m.Lock()
	respChan, exist := s.requests[id]
	if !exist {
		return nil, errors.New(fmt.Sprintf("a request channel for id (%s) does not exist", id))
	}
	delete(s.requests, id)
	s.m.Unlock()
	logger.Debug(fmt.Sprintf("fetched response channel: for id: %s", id))
	return respChan, nil

}

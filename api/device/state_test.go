package device_api

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestState(t *testing.T) {

	s := newState()

	testChan := make(chan Response)

	//Register test channel
	id := s.Add(testChan)

	//Check if successfully registered
	s.m.Lock()
	require.Equal(t, testChan, s.requests[id])
	s.m.Unlock()

	//Cut should remove the channel from the state and return it
	registeredChan, err := s.Cut(id)
	require.Nil(t, err)
	require.Equal(t, testChan, registeredChan)

	//Cutting a already received channel should as well result in an error
	registeredChan, err = s.Cut(id)
	require.EqualError(t, err, "a request channel for id (4039455774) does not exist")

}

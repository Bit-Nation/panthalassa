package modal

import (
	"testing"
	"errors"
	"time"
	
	log "github.com/op/go-logging"
	require "github.com/stretchr/testify/require"
	otto "github.com/robertkrimen/otto"
)

type testDevice struct {
	handler func(title, layout string) error
}

func (d testDevice) ShowModal(title, layout string) error {
	return d.handler(title, layout)
}

func TestWithoutCallback(t *testing.T) {
	
	// get test logger
	logger := log.MustGetLogger("")
	
	// create VM
	vm := otto.New()
	
	// create modal module
	modalModule := New(logger, testDevice{
		handler: func(title, layout string) error {
			require.Equal(t, "my title", title)
			require.Equal(t, "{}", layout)
			return nil
		},
	})
	
	// register show modal functionality
	require.Nil(t, modalModule.Register(vm))
	
	// try to display modal
	_, err := vm.Call(`showModal`, vm, "my title", "{}")
	require.Nil(t, err)
	
}

func TestFirstParamString(t *testing.T)  {
	
	// get test logger
	logger := log.MustGetLogger("")
	
	// create VM
	vm := otto.New()
	
	// create modal module
	modalModule := New(logger, nil)
	
	// register show modal functionality
	require.Nil(t, modalModule.Register(vm))
	
	// try to display modal
	showModalError, err := vm.Call(`showModal`, vm)
	require.Nil(t, err)
	require.Equal(t, "ValidationError: expected parameter 0 to be of type string", showModalError.String())
	
}

func TestSecondParamString(t *testing.T)  {
	
	// get test logger
	logger := log.MustGetLogger("")
	
	// create VM
	vm := otto.New()
	
	// create modal module
	modalModule := New(logger, nil)
	
	// register show modal functionality
	require.Nil(t, modalModule.Register(vm))
	
	// try to display modal
	showModalError, err := vm.Call(`showModal`, vm, ``)
	require.Nil(t, err)
	require.Equal(t, "ValidationError: expected parameter 1 to be of type string", showModalError.String())
	
}

func TestErrorForwarding(t *testing.T)  {
	
	// get test logger
	logger := log.MustGetLogger("")
	
	// create VM
	vm := otto.New()
	
	// create modal module
	modalModule := New(logger, testDevice{
		handler: func(title, layout string) error {
			// return test error
			return errors.New("test error")
		},
	})
	
	// register show modal functionality
	require.Nil(t, modalModule.Register(vm))
	
	c := make(chan bool)
	
	go func() {
		select {
		case <-c:
			return
		case <- time.After(time.Second * 1):
			require.FailNow(t, "time out")
		}
	}()
	
	// try to display modal
	_, err := vm.Call(`showModal`, vm, "", "", func(err string) {
	
		if err != "test error" {
			panic("expected error")
		}

		c <- true
		
	})
	require.Nil(t, err)
	
}
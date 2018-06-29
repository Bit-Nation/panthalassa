package modal

import (
	"errors"
	"testing"
	"time"

	"crypto/rand"
	log "github.com/op/go-logging"
	otto "github.com/robertkrimen/otto"
	require "github.com/stretchr/testify/require"
	ed25519 "golang.org/x/crypto/ed25519"
)

type testDevice struct {
	handler func(title, layout string, dAppIDKey ed25519.PublicKey) error
}

func (d testDevice) ShowModal(title, layout string, dAppIDKey ed25519.PublicKey) error {
	return d.handler(title, layout, dAppIDKey)
}

func TestWithoutCallback(t *testing.T) {

	// get test logger
	logger := log.MustGetLogger("")

	// create VM
	vm := otto.New()

	pub, _, err := ed25519.GenerateKey(rand.Reader)
	require.Nil(t, err)

	// create modal module
	modalModule := New(logger, testDevice{
		handler: func(title, layout string, dAppIdKey ed25519.PublicKey) error {
			require.Equal(t, "my title", title)
			require.Equal(t, "{}", layout)
			return nil
		},
	}, pub)

	// register show modal functionality
	require.Nil(t, modalModule.Register(vm))

	// try to display modal
	_, err = vm.Call(`showModal`, vm, "my title", "{}")
	require.Nil(t, err)

}

func TestFirstParamString(t *testing.T) {

	// get test logger
	logger := log.MustGetLogger("")

	// create VM
	vm := otto.New()

	pub, _, err := ed25519.GenerateKey(rand.Reader)
	require.Nil(t, err)

	// create modal module
	modalModule := New(logger, nil, pub)

	// register show modal functionality
	require.Nil(t, modalModule.Register(vm))

	// try to display modal
	showModalError, err := vm.Call(`showModal`, vm, nil)
	require.Nil(t, err)
	require.Equal(t, "ValidationError: expected parameter 0 to be of type string", showModalError.String())

}

func TestSecondParamString(t *testing.T) {

	// get test logger
	logger := log.MustGetLogger("")

	// create VM
	vm := otto.New()

	pub, _, err := ed25519.GenerateKey(rand.Reader)
	require.Nil(t, err)

	// create modal module
	modalModule := New(logger, nil, pub)

	// register show modal functionality
	require.Nil(t, modalModule.Register(vm))

	// try to display modal
	showModalError, err := vm.Call(`showModal`, vm, ``)
	require.Nil(t, err)
	require.Equal(t, "ValidationError: expected parameter 1 to be of type string", showModalError.String())

}

func TestThirdParamCallback(t *testing.T) {

	// get test logger
	logger := log.MustGetLogger("")

	// create VM
	vm := otto.New()

	pub, _, err := ed25519.GenerateKey(rand.Reader)
	require.Nil(t, err)

	// create modal module
	modalModule := New(logger, nil, pub)

	// register show modal functionality
	require.Nil(t, modalModule.Register(vm))

	// try to display modal
	showModalError, err := vm.Call(`showModal`, vm, ``, ``)
	require.Nil(t, err)
	require.Equal(t, "ValidationError: expected parameter 2 to be of type function", showModalError.String())

}

func TestErrorForwarding(t *testing.T) {

	// get test logger
	logger := log.MustGetLogger("")

	// create VM
	vm := otto.New()

	pub, _, err := ed25519.GenerateKey(rand.Reader)
	require.Nil(t, err)

	// create modal module
	modalModule := New(logger, testDevice{
		handler: func(title, layout string, key ed25519.PublicKey) error {
			// return test error
			return errors.New("test error")
		},
	}, pub)

	// register show modal functionality
	require.Nil(t, modalModule.Register(vm))

	c := make(chan bool, 1)

	// try to display modal
	_, err = vm.Call(`showModal`, vm, "", "", func(err string) {

		if err != "test error" {
			panic("expected error")
		}

		c <- true

	})

	select {
	case <-c:
	case <-time.After(time.Second * 1):
		require.FailNow(t, "time out")
	}

	require.Nil(t, err)

}

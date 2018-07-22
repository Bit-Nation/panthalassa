package modal

import (
	"testing"
	"time"

	log "github.com/op/go-logging"
	otto "github.com/robertkrimen/otto"
	uuid "github.com/satori/go.uuid"
	require "github.com/stretchr/testify/require"
	ed25519 "golang.org/x/crypto/ed25519"
)

type testDevice struct {
	handler func(uiID, layout, renderType string, dAppPubKey ed25519.PublicKey) error
}

func (d testDevice) RenderModal(uiID, layout, renderType string, dAppPubKey ed25519.PublicKey) error {
	return d.handler(uiID, layout, renderType, dAppPubKey)
}

func TestModule_CloseModal(t *testing.T) {

	// create modal module
	logger := log.MustGetLogger("")
	m := New(logger, nil, []byte(""))
	vm := otto.New()
	require.Nil(t, m.Register(vm))

	// closer
	closed := false
	closer := func() {
		closed = true
	}

	closeTest := make(chan struct{}, 1)

	// create new uuid
	_, err := vm.Call("newModalUIID", vm, closer, func(call otto.FunctionCall) otto.Value {

		// fetch callback data
		err := call.Argument(0)
		modalID := call.Argument(1)

		// error must be undefined
		require.True(t, err.IsUndefined())

		// convert returned id to uuid
		id, convertErr := uuid.FromString(modalID.String())
		require.Nil(t, convertErr)
		require.Equal(t, modalID.String(), id.String())

		// id must be registered in modal id map
		_, exist := m.modalIDs[id.String()]
		require.True(t, exist)

		// close modal
		m.CloseModal(id.String())

		// id must NOT be registered in modal id map
		_, exist = m.modalIDs[id.String()]
		require.False(t, exist)

		require.True(t, closed)

		// close test
		closeTest <- struct{}{}
		return otto.Value{}
	})
	require.Nil(t, err)

	select {
	case <-closeTest:
		require.Nil(t, err)
	case <-time.After(time.Second * 3):
		require.Fail(t, "timed out")
	}

}

func TestModule_RenderModal(t *testing.T) {

	uiID := "i_am_the_ui_id"

	calledDevice := false
	device := &testDevice{
		handler: func(receivedUiID, layout, receivedRenderType string, dAppPubKey ed25519.PublicKey) error {
			calledDevice = true
			require.Equal(t, uiID, receivedUiID)
			require.Equal(t, "{jsx: 'tree'}", layout)
			require.Equal(t, renderType, receivedRenderType)
			require.Equal(t, "id pub key", string(dAppPubKey))
			return nil
		},
	}

	// create module
	logger := log.MustGetLogger("")
	m := New(logger, device, []byte("id pub key"))
	vm := otto.New()
	require.Nil(t, m.Register(vm))

	// we just register a fake it here to just
	// make sure that we have an ID in the vm
	m.lock.Lock()
	m.modalIDs[uiID] = &otto.Value{}
	m.lock.Unlock()

	done := make(chan struct{}, 1)

	_, err := vm.Call("renderModal", vm, uiID, "{jsx: 'tree'}", func(call otto.FunctionCall) otto.Value {

		// make sure device has been called
		require.True(t, calledDevice)

		// close test
		done <- struct{}{}

		return otto.Value{}
	})
	require.Nil(t, err)

	select {
	case <-done:
	case <-time.After(time.Second * 2):
		require.Fail(t, "time out")
	}

}

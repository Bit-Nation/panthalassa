package modal

import (
	"testing"
	"time"

	log "github.com/op/go-logging"
	uuid "github.com/satori/go.uuid"
	require "github.com/stretchr/testify/require"
	ed25519 "golang.org/x/crypto/ed25519"
	duktape "gopkg.in/olebedev/go-duktape.v3"
)

type testDevice struct {
	handler func(uiID, layout string, dAppPubKey ed25519.PublicKey) error
}

func (d testDevice) RenderModal(uiID, layout string, dAppPubKey ed25519.PublicKey) error {
	return d.handler(uiID, layout, dAppPubKey)
}

func TestModule_CloseModal(t *testing.T) {

	// create modal module
	logger := log.MustGetLogger("")
	m := New(logger, nil, []byte(""))
	vm := duktape.New()
	require.Nil(t, m.Register(vm))

	// closer
	closed := false

	_, err := vm.PushGlobalGoFunction("callbackCloserCloseModal", func(context *duktape.Context) int {
		closed = true
		return 0
	})
	require.Nil(t, err)

	closeTest := make(chan struct{}, 1)

	// create new uuid
	_, err = vm.PushGlobalGoFunction("callbackTestModuleCloseModal", func(context *duktape.Context) int {
		// fetch callback data
		errBool := context.IsUndefined(0)
		modalID := context.ToString(1)

		// error must be undefined
		require.True(t, errBool)
		// convert returned id to uuid
		id, convertErr := uuid.FromString(modalID)
		require.Nil(t, convertErr)
		require.Equal(t, modalID, id.String())
		// id must be registered in modal id map
		respChan := make(chan *duktape.Context)
		m.fetchModalCloserChan <- fetchModalCloser{
			id:       id.String(),
			respChan: respChan,
		}
		require.NotNil(t, <-respChan)
		// close modal
		vm.PevalString(`callbackCloserCloseModal`)
		m.CloseModal(id.String())

		// id must NOT be registered in modal id map
		respChan = make(chan *duktape.Context)
		m.fetchModalCloserChan <- fetchModalCloser{
			id:       id.String(),
			respChan: respChan,
		}
		require.Nil(t, <-respChan)
		require.True(t, closed)

		// close test
		closeTest <- struct{}{}
		return 0
	})
	require.Nil(t, err)
	err = vm.PevalString(`newModalUIID(callbackCloserCloseModal,callbackTestModuleCloseModal)`)
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
		handler: func(receivedUIID, receivedLayout string, dAppPubKey ed25519.PublicKey) error {
			calledDevice = true
			require.Equal(t, uiID, receivedUIID)
			require.Equal(t, "{jsx: 'tree'}", receivedLayout)
			require.Equal(t, "id pub key", string(dAppPubKey))
			return nil
		},
	}

	// create module
	logger := log.MustGetLogger("")
	m := New(logger, device, []byte("id pub key"))
	vm := duktape.New()
	require.Nil(t, m.Register(vm))

	// we just register a fake it here to just
	// make sure that we have an ID in the vm
	m.addModalIDChan <- addModalID{
		id:     uiID,
		closer: &duktape.Context{},
	}

	done := make(chan struct{}, 1)
	_, err := vm.PushGlobalGoFunction("callbackTestModuleCloseModal", func(context *duktape.Context) int {
		// make sure device has been called
		require.True(t, calledDevice)

		// close test
		done <- struct{}{}

		return 0
	})
	require.Nil(t, err)
	err = vm.PevalString(`renderModal("` + uiID + `","{jsx: 'tree'}",callbackTestModuleCloseModal)`)
	require.Nil(t, err)

	select {
	case <-done:
	case <-time.After(time.Second * 2):
		require.Fail(t, "time out")
	}

}

func TestModal_RenderWithoutID(t *testing.T) {

	// create module
	logger := log.MustGetLogger("")
	m := New(logger, nil, []byte("id pub key"))
	vm := duktape.New()
	require.Nil(t, m.Register(vm))

	done := make(chan struct{}, 1)
	_, err := vm.PushGlobalGoFunction("callbackTestModalRenderWithoutID", func(context *duktape.Context) int {
		err := context.ToString(0)
		require.Equal(t, "MissingModalID: modal UI ID: 'id_do_not_exist' does not exist", err)
		done <- struct{}{}
		return 0
	})
	require.Nil(t, err)
	err = vm.PevalString(`renderModal("id_do_not_exist", "", callbackTestModalRenderWithoutID)`)
	require.Nil(t, err)

	select {
	case <-done:
	case <-time.After(time.Second * 3):
		require.Fail(t, "timed out")
	}

}

func TestModal_RequestLimitation(t *testing.T) {

	// create module
	logger := log.MustGetLogger("")
	m := New(logger, nil, []byte("id pub key"))
	vm := duktape.New()
	require.Nil(t, m.Register(vm))

	_, err := vm.PushGlobalGoFunction("callbackCloserRequestLimitation", func(context *duktape.Context) int {
		return 0
	})
	require.Nil(t, err)

	// make sure that current amount of registered ids is 0
	require.Equal(t, uint(0), m.modalIDsReqLim.Current())

	// closer
	done := make(chan struct{}, 1)
	_, err = vm.PushGlobalGoFunction("callbackTestModalRequestLimitation", func(context *duktape.Context) int {
		// newModalUIID must register a new id
		require.Equal(t, uint(1), m.modalIDsReqLim.Current())

		// close modal with UI ID
		id := context.ToString(1)
		vm.PevalString(`callbackCloserRequestLimitation`)
		m.CloseModal(id)

		// wait a bit to sync go routines
		time.Sleep(time.Millisecond * 100)

		// make sure id was removed from current
		require.Equal(t, uint(0), m.modalIDsReqLim.Current())

		done <- struct{}{}

		return 0
	})
	require.Nil(t, err)
	err = vm.PevalString(`newModalUIID(callbackCloserRequestLimitation,callbackTestModalRequestLimitation)`)
	require.Nil(t, err)
	select {
	case <-done:
	case <-time.After(time.Second * 3):
	}

}

package modal

import (
	"errors"
	"fmt"
	"time"

	reqLim "github.com/Bit-Nation/panthalassa/dapp/request_limitation"
	validator "github.com/Bit-Nation/panthalassa/dapp/validator"
	log "github.com/ipfs/go-log"
	logger "github.com/op/go-logging"
	uuid "github.com/satori/go.uuid"
	ed25519 "golang.org/x/crypto/ed25519"
	duktape "gopkg.in/olebedev/go-duktape.v3"
)

type Device interface {
	RenderModal(uiID, layout string, dAppPubKey ed25519.PublicKey) error
}

var sysLog = log.Logger("modal")

type addModalID struct {
	id     string
	closer *duktape.Context
}

type fetchModalCloser struct {
	id       string
	respChan chan *duktape.Context
}

type Module struct {
	device               Device
	logger               *logger.Logger
	dAppIDKey            ed25519.PublicKey
	modalIDsReqLim       *reqLim.CountThrottling
	addModalIDChan       chan addModalID
	fetchModalCloserChan chan fetchModalCloser
	deleteModalID        chan string
}

const renderType = "modal"

// create new Modal Module
func New(l *logger.Logger, device Device, dAppIDKey ed25519.PublicKey) *Module {
	m := &Module{
		device:               device,
		logger:               l,
		dAppIDKey:            dAppIDKey,
		modalIDsReqLim:       reqLim.NewCountThrottling(6, time.Minute, 20, errors.New("can't put more show modal requests in queue")),
		addModalIDChan:       make(chan addModalID),
		fetchModalCloserChan: make(chan fetchModalCloser),
		deleteModalID:        make(chan string),
	}

	go func() {

		modals := map[string]*duktape.Context{}

		// exit if channels are closed
		if m.addModalIDChan == nil || m.fetchModalCloserChan == nil || m.deleteModalID == nil {
			return
		}

		for {
			select {
			// add
			case add := <-m.addModalIDChan:
				modals[add.id] = add.closer
			// fetch
			case fetch := <-m.fetchModalCloserChan:
				closer, exist := modals[fetch.id]
				if !exist {
					fetch.respChan <- closer
					continue
				}
				fetch.respChan <- closer
			// delete
			case delID := <-m.deleteModalID:
				delete(modals, delID)
				m.modalIDsReqLim.Decrease()
			}
		}

	}()

	return m
}

func (m *Module) Close() error {
	close(m.addModalIDChan)
	close(m.fetchModalCloserChan)
	close(m.deleteModalID)
	return nil
}

// renderModal provides a way to display a modal
// the first parameter should be the modal title
// the second parameter should be the layout to render
// and the third parameter is an optional callback that
// will called with an optional error
func (m *Module) Register(vm *duktape.Context) error {

	fmt.Println("[DApp] Register Start")

	fmt.Println("[DApp] renderModal Start")
	_, err := vm.PushGlobalGoFunction("renderModal", func(context *duktape.Context) int {
		sysLog.Debug("render modal")

		// validate function call
		v := validator.New()
		// ui id
		v.Set(0, &validator.TypeString)
		// layout
		v.Set(1, &validator.TypeString)
		// callback
		v.Set(2, &validator.TypeFunction)
		if err := v.Validate(context); err != nil {
			m.logger.Error(err.Error())
			return 1
		}

		// modal ui id
		uiID := context.SafeToString(0)

		// make sure ui id is registered
		mIdChan := make(chan *duktape.Context)
		m.fetchModalCloserChan <- fetchModalCloser{
			id:       uiID,
			respChan: mIdChan,
		}
		modalCloser := <-mIdChan
		// a modal closer can only exist with relation to a modal ID.
		// So if the modal closer exist the modal id exist as well
		if modalCloser == nil {
			errMsg := fmt.Sprintf("MissingModalID: modal UI ID: '%s' does not exist", uiID)
			context.PushString(errMsg)
			context.PushString(uiID)
			context.Call(2)
			return 1
		}

		// get layout
		layout := context.SafeToString(1)

		sysLog.Debugf("going to render: %s with UI id %s", layout, uiID)

		// execute show modal action in
		// context of request limitation
		// @TODO figure out why concurrently using "go func" instead of "func" makes the vm/context lose the stack
		func() {
			// request to show modal
			if err := m.device.RenderModal(uiID, layout, m.dAppIDKey); err != nil {
				context.PushString(`Error: failed to render modal ` + err.Error())
				context.Call(1)
				return
			}

			context.PushUndefined()
			context.Call(1)

		}()

		return 0

	})
	fmt.Println("[DApp] renderModal End")

	if err != nil {
		return err
	}

	fmt.Println("[DApp] newModalUIID Start")
	_, err = vm.PushGlobalGoFunction("newModalUIID", func(context *duktape.Context) int {
		// validate function call
		v := validator.New()
		// close handler
		v.Set(0, &validator.TypeFunction)
		// callback
		v.Set(1, &validator.TypeFunction)
		if err := v.Validate(context); err != nil {
			m.logger.Error(err.Error())
			return 1
		}

		// create new id
		id, err := uuid.NewV4()
		if err != nil {
			context.PushString("Error " + err.Error())
			context.Call(1)
			return 1
		}

		// increase request limitation counter
		throttlingFunc := func(dec chan struct{}) {
			// add ui id & closer to stack
			m.addModalIDChan <- addModalID{
				id:     id.String(),
				closer: context,
			}

			// call callback
			context.PopN(0)
			context.PushUndefined()
			context.PushString(id.String())
			context.Call(2)

			//@TODO find a way to determine if the callback errored out
			//if err != nil {
			// in the case of an error we would like to remove the id
			// and decrease our request limitation
			//	m.deleteModalID <- id.String()
			//	m.logger.Error(err.Error())
			//}
		}
		dec := make(chan struct{}, 1)
		throttlingFunc(dec)
		//m.modalIDsReqLim.Exec(throttlingFunc)
		return 0

	})
	fmt.Println("[DApp] newModalUIID/Register End")
	return err
}

// close the modal
func (m *Module) CloseModal(uiID string) {

	// fetch closer
	respChan := make(chan *duktape.Context)
	m.fetchModalCloserChan <- fetchModalCloser{
		id:       uiID,
		respChan: respChan,
	}
	closer := <-respChan
	if closer == nil {
		return
	}

	// close modal
	closer.Call(0)

	// finish close (decrease counter & delete from id's")
	m.deleteModalID <- uiID

}

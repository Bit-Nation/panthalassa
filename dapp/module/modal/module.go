package modal

import (
	"errors"
	"fmt"
	"time"

	reqLim "github.com/Bit-Nation/panthalassa/dapp/request_limitation"
	validator "github.com/Bit-Nation/panthalassa/dapp/validator"
	log "github.com/ipfs/go-log"
	logger "github.com/op/go-logging"
	otto "github.com/robertkrimen/otto"
	uuid "github.com/satori/go.uuid"
	ed25519 "golang.org/x/crypto/ed25519"
)

type Device interface {
	RenderModal(uiID, layout string, dAppPubKey ed25519.PublicKey) error
}

var sysLog = log.Logger("modal")

type addModalID struct {
	id     string
	closer *otto.Value
}

type fetchModalCloser struct {
	id       string
	respChan chan *otto.Value
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

		modals := map[string]*otto.Value{}

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
func (m *Module) Register(vm *otto.Otto) error {

	err := vm.Set("renderModal", func(call otto.FunctionCall) otto.Value {

		sysLog.Debug("render modal")

		// validate function call
		v := validator.New()
		// ui id
		v.Set(0, &validator.TypeString)
		// layout
		v.Set(1, &validator.TypeString)
		// callback
		v.Set(2, &validator.TypeFunction)
		if err := v.Validate(vm, call); err != nil {
			m.logger.Error(err.String())
			return otto.Value{}
		}

		// modal ui id
		uiID := call.Argument(0).String()

		// callback
		cb := call.Argument(2)

		// make sure ui id is registered
		mIdChan := make(chan *otto.Value)
		m.fetchModalCloserChan <- fetchModalCloser{
			id:       uiID,
			respChan: mIdChan,
		}
		modalCloser := <-mIdChan
		// a modal closer can only exist with relation to a modal ID.
		// So if the modal closer exist the modal id exist as well
		if modalCloser == nil {
			errMsg := fmt.Sprintf("modal UI ID: '%s' does not exist", uiID)
			if _, err := cb.Call(cb, vm.MakeCustomError("MissingModalID", errMsg), uiID); err != nil {
				m.logger.Error(errMsg)
			}
			return otto.Value{}
		}

		// get layout
		layout := call.Argument(1).String()

		// execute show modal action in
		// context of request limitation
		go func() {
			// request to show modal
			if err := m.device.RenderModal(uiID, layout, m.dAppIDKey); err != nil {
				cb.Call(cb, vm.MakeCustomError("Error", "failed to render modal"))
				return
			}

			cb.Call(cb)
		}()

		return otto.Value{}

	})

	if err != nil {
		return err
	}

	return vm.Set("newModalUIID", func(call otto.FunctionCall) otto.Value {

		// validate function call
		v := validator.New()
		// close handler
		v.Set(0, &validator.TypeFunction)
		// callback
		v.Set(1, &validator.TypeFunction)
		if err := v.Validate(vm, call); err != nil {
			m.logger.Error(err.String())
			return otto.Value{}
		}

		cb := call.Argument(1)

		// create new id
		id, err := uuid.NewV4()
		if err != nil {
			_, err = cb.Call(cb, vm.MakeCustomError("Error", err.Error()))
			if err != nil {
				m.logger.Error(err.Error())
			}
			return otto.Value{}
		}

		// increase request limitation counter
		m.modalIDsReqLim.Exec(func(dec chan struct{}) {

			// add ui id & closer to stack
			closer := call.Argument(0)
			m.addModalIDChan <- addModalID{
				id:     id.String(),
				closer: &closer,
			}

			// call callback
			_, err = cb.Call(cb, nil, id.String())
			if err != nil {
				// in the case of an error we would like to remove the id
				// and decrease our request limitation
				m.deleteModalID <- id.String()
				m.logger.Error(err.Error())
			}
		})

		return otto.Value{}

	})

}

// close the modal
func (m *Module) CloseModal(uiID string) {

	// fetch closer
	respChan := make(chan *otto.Value)
	m.fetchModalCloserChan <- fetchModalCloser{
		id:       uiID,
		respChan: respChan,
	}
	closer := <-respChan
	if closer == nil {
		return
	}

	// close modal
	if _, err := closer.Call(*closer); err != nil {
		m.logger.Error(err.Error())
	}

	// finish close (decrease counter & delete from id's")
	m.deleteModalID <- uiID

}

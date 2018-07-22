package modal

import (
	"errors"
	"fmt"
	"sync"
	"time"

	reqLim "github.com/Bit-Nation/panthalassa/dapp/request_limitation"
	validator "github.com/Bit-Nation/panthalassa/dapp/validator"
	log "github.com/op/go-logging"
	otto "github.com/robertkrimen/otto"
	uuid "github.com/satori/go.uuid"
	ed25519 "golang.org/x/crypto/ed25519"
)

type Device interface {
	RenderModal(uiID, layout string, dAppPubKey ed25519.PublicKey) error
}

type Module struct {
	device         Device
	logger         *log.Logger
	dAppIDKey      ed25519.PublicKey
	modalIDs       map[string]*otto.Value
	modalIDsReqLim *reqLim.CountThrottling
	lock           sync.Mutex
}

const renderType = "modal"

// create new Modal Module
func New(l *log.Logger, device Device, dAppIDKey ed25519.PublicKey) *Module {
	return &Module{
		device:         device,
		logger:         l,
		dAppIDKey:      dAppIDKey,
		modalIDs:       map[string]*otto.Value{},
		modalIDsReqLim: reqLim.NewCountThrottling(6, time.Minute, 20, errors.New("can't put more show modal requests in queue")),
		lock:           sync.Mutex{},
	}
}

func (m *Module) Register(vm *otto.Otto) error {

	err := vm.Set("renderModal", func(call otto.FunctionCall) otto.Value {

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
		m.lock.Lock()
		defer m.lock.Unlock()
		if _, exist := m.modalIDs[uiID]; !exist {
			errMsg := fmt.Sprintf("modal UI ID: '%s' does not exist", uiID)
			_, err := cb.Call(cb, vm.MakeCustomError("MissingModalID", errMsg), uiID)
			if err != nil {
				m.logger.Error(errMsg)
			}
			return otto.Value{}
		}

		// make sure layout is valid
		// @todo validate layout
		layout := call.Argument(1).String()

		// execute show modal action in
		// context of request limitation
		go func() {
			// request to show modal
			err := m.device.RenderModal(
				uiID,
				layout,
				m.dAppIDKey,
			)

			if err != nil {
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

		m.lock.Lock()
		defer m.lock.Unlock()

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
		m.modalIDsReqLim.Exec(func() {
			// add ui id to stack
			closer := call.Argument(0)
			m.modalIDs[id.String()] = &closer

			// call callback
			_, err = cb.Call(cb, nil, id.String())
			if err != nil {
				// in the case of an error we would like to remove the id
				// and decrease our request limitation
				delete(m.modalIDs, id.String())
				m.modalIDsReqLim.Decrease()

				m.logger.Error(err.Error())
			}
		})

		return otto.Value{}

	})

}

// close the modal
func (m *Module) CloseModal(uiID string) {

	m.lock.Lock()
	defer m.lock.Unlock()

	closeModal := m.modalIDs[uiID]
	if _, err := closeModal.Call(*closeModal); err != nil {
		m.logger.Error(err.Error())
	}
	m.modalIDsReqLim.Decrease()
	delete(m.modalIDs, uiID)

}

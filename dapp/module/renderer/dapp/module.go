package dapp

import (
	"errors"

	validator "github.com/Bit-Nation/panthalassa/dapp/validator"
	log "github.com/ipfs/go-log"
	logger "github.com/op/go-logging"
	otto "github.com/robertkrimen/otto"
)

type Module struct {
	logger             *logger.Logger
	vm                 *otto.Otto
	setOpenHandlerChan chan *otto.Value
	getOpenHandlerChan chan chan *otto.Value
	addCBChan          chan *chan error
	rmCBChan           chan *chan error
	// returns a cb chan from the stack
	nextCBChan chan chan *chan error
}

var sysLog = log.Logger("renderer - dapp")

func (m *Module) Register(vm *otto.Otto) error {
	m.vm = vm
	// setOpenHandler must be called with a callback
	// the callback that is passed to `setOpenHandler`
	// will be called with an "data" object and a callback
	// the callback should be called (with an optional error)
	// in order to return from the function
	return vm.Set("setOpenHandler", func(call otto.FunctionCall) otto.Value {

		//sysLog.Debug("set open handler")

		// validate function call
		v := validator.New()
		v.Set(0, &validator.TypeFunction)
		if err := v.Validate(vm, call); err != nil {
			m.logger.Error(err.String())
			return *err
		}

		// set renderer
		fn := call.Argument(0)
		m.setOpenHandlerChan <- &fn

		return otto.Value{}
	})
}

// payload can be an arbitrary set of key value pairs (as a json string)
func (m *Module) OpenDApp(payload string) error {

	// fetch handler
	handlerChan := make(chan *otto.Value)
	m.getOpenHandlerChan <- handlerChan
	handler := <-handlerChan

	// make sure an renderer has been set
	if handler == nil {
		return errors.New("failed to open DApp - no open handler set")
	}

	// convert context to otto js object
	dataObj, err := m.vm.Object("(" + payload + ")")
	if err != nil {
		return err
	}

	cbDone := make(chan error)

	// add cb chan to state
	m.addCBChan <- &cbDone

	go func() {

		// call the renderer
		// we pass in data object and a callback
		_, err = handler.Call(*handler, dataObj, func(call otto.FunctionCall) otto.Value {

			// remove cb chan from state
			defer func() {
				m.rmCBChan <- &cbDone
			}()

			// fetch params from the callback call
			err := call.Argument(0)

			// if there is an error, set it in the response
			if !err.IsUndefined() {
				cbDone <- errors.New(err.String())
				return otto.Value{}
			}

			cbDone <- nil

			return otto.Value{}

		})

		if err != nil {
			m.logger.Error(err.Error())
		}

	}()

	return <-cbDone
}

func (m *Module) Close() error {
	close(m.setOpenHandlerChan)
	close(m.getOpenHandlerChan)
	// close all open callback
	for {
		// fetch next response channel
		respChan := make(chan *chan error)
		m.nextCBChan <- respChan
		cbChan := <-respChan

		// exit if there is no next channel
		if cbChan == nil {
			break
		}

		// send error in response
		*cbChan <- errors.New("closed the application")
		m.rmCBChan <- cbChan
	}
	close(m.addCBChan)
	close(m.rmCBChan)
	close(m.nextCBChan)
	return nil
}

func New(l *logger.Logger) *Module {

	m := &Module{
		logger:             l,
		setOpenHandlerChan: make(chan *otto.Value),
		getOpenHandlerChan: make(chan chan *otto.Value),
		addCBChan:          make(chan *chan error),
		rmCBChan:           make(chan *chan error),
		nextCBChan:         make(chan chan *chan error),
	}

	go func() {

		openHandler := new(otto.Value)
		cbChans := map[*chan error]bool{}

		for {

			// exit if channels got closed
			if m.setOpenHandlerChan == nil || m.getOpenHandlerChan == nil || m.addCBChan == nil || m.rmCBChan == nil || m.nextCBChan == nil {
				return
			}

			select {
			case h := <-m.setOpenHandlerChan:
				//@todo somehow there is a bug related to otto were it leads to
				//@todo that setOpenHandler is called twice but wil nil for the handler
				//@todo once we get rid of otto we can remove this check
				if h == nil {
					continue
				}
				openHandler = h
			case respChan := <-m.getOpenHandlerChan:
				//@todo somehow there is a bug related to otto were it leads to
				//@todo that setOpenHandler is called twice but wil nil for the handler
				//@todo once we get rid of otto we can remove this check
				if respChan == nil {
					continue
				}
				respChan <- openHandler
			case cb := <-m.addCBChan:
				cbChans[cb] = true
			case cb := <-m.rmCBChan:
				delete(cbChans, cb)
			case respChan := <-m.nextCBChan:
				if len(cbChans) == 0 {
					respChan <- nil
					continue
				}
				// my a bit ugly, but couldn't think of how to fetch the next callback
				// in another way from a map
				for cbChan, _ := range cbChans {
					respChan <- cbChan
					break
				}
			}

		}

	}()

	return m
}

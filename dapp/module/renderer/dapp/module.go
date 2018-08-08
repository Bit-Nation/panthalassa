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

		sysLog.Debug("set open handler")

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

	c := make(chan error, 1)

	go func() {

		// call the renderer
		// we pass in data object and a callback
		_, err = handler.Call(*handler, dataObj, func(call otto.FunctionCall) otto.Value {

			// fetch params from the callback call
			err := call.Argument(0)

			// if there is an error, set it in the response
			if !err.IsUndefined() {
				c <- errors.New(err.String())
				return otto.Value{}
			}

			c <- nil

			return otto.Value{}

		})

		if err != nil {
			m.logger.Error(err.Error())
		}

	}()

	return <-c
}

func (m *Module) Close() error {
	close(m.setOpenHandlerChan)
	close(m.getOpenHandlerChan)
	return nil
}

func New(l *logger.Logger) *Module {

	m := &Module{
		logger:             l,
		setOpenHandlerChan: make(chan *otto.Value),
		getOpenHandlerChan: make(chan chan *otto.Value),
	}

	go func() {

		openHandler := new(otto.Value)

		for {

			// exit if channels got closed
			if m.setOpenHandlerChan == nil || m.getOpenHandlerChan == nil {
				return
			}

			select {
			case h := <-m.setOpenHandlerChan:
				openHandler = h
			case respChan := <-m.getOpenHandlerChan:
				respChan <- openHandler
			}

		}

	}()

	return m
}

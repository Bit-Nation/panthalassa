package dapp

import (
	"errors"

	validator "github.com/Bit-Nation/panthalassa/dapp/validator"
	log "github.com/ipfs/go-log"
	logger "github.com/op/go-logging"
	duktape "gopkg.in/olebedev/go-duktape.v3"
)

type Module struct {
	logger             *logger.Logger
	vm                 *duktape.Context
	setOpenHandlerChan chan *duktape.Context
	getOpenHandlerChan chan chan *duktape.Context
	addCBChan          chan *chan error
	rmCBChan           chan *chan error
	// returns a cb chan from the stack
	nextCBChan chan chan *chan error
	closer     chan struct{}
}

var sysLog = log.Logger("renderer - dapp")

// setOpenHandler must be called with a callback
// the callback that is passed to `setOpenHandler`
// will be called with an "data" object and a callback
// the callback should be called (with an optional error)
// in order to return from the function
func (m *Module) Register(vm *duktape.Context) error {
	m.vm = vm
	_, err := vm.PushGlobalGoFunction("setOpenHandler", func(context *duktape.Context) int {

		sysLog.Debug("set open handler")

		// validate function call
		v := validator.New()
		v.Set(0, &validator.TypeFunction)
		if err := v.Validate(context); err != nil {
			m.logger.Error(err.Error())
			return 1
		}

		// set renderer
		m.setOpenHandlerChan <- context

		return 0
	})
	return err
}

// payload can be an arbitrary set of key value pairs (as a json string)
func (m *Module) OpenDApp(payload string) error {

	// fetch handler
	handlerChan := make(chan *duktape.Context)
	m.getOpenHandlerChan <- handlerChan
	handler := <-handlerChan

	// make sure an renderer has been set
	if handler == nil {
		return errors.New("failed to open DApp - no open handler set")
	}

	// convert context to otto js object
	dataObj := "(" + payload + ")"

	cbDone := make(chan error)

	// add cb chan to state
	m.addCBChan <- &cbDone

	go func() {

		// call the renderer
		// we pass in data object and a callback
		_, err := handler.PushGlobalGoFunction("callbackOpenDApp", func(context *duktape.Context) int {

			// remove cb chan from state
			defer func() {
				m.rmCBChan <- &cbDone
			}()

			// if there is an error, set it in the response
			if !context.IsUndefined(0) && !context.IsNull(0) {
				callbackerr := context.ToString(0)
				cbDone <- errors.New(callbackerr)
				return 1
			}

			cbDone <- nil

			return 0

		})

		if err != nil {
			m.logger.Error(err.Error())
		}
		err = handler.PevalString(dataObj)
		if err != nil {
			m.logger.Error(err.Error())
		}
		err = handler.PevalString(`callbackOpenDApp`)
		if err != nil {
			m.logger.Error(err.Error())
		}

		handler.Call(2)

	}()

	return <-cbDone
}

func (m *Module) Close() error {
	m.closer <- struct{}{}
	return nil
}

func New(l *logger.Logger) *Module {

	m := &Module{
		logger:             l,
		setOpenHandlerChan: make(chan *duktape.Context),
		getOpenHandlerChan: make(chan chan *duktape.Context),
		addCBChan:          make(chan *chan error),
		rmCBChan:           make(chan *chan error),
		nextCBChan:         make(chan chan *chan error),
		closer:             make(chan struct{}),
	}

	go func() {

		openHandler := new(duktape.Context)
		cbChans := map[*chan error]bool{}

		for {

			select {
			case <-m.closer:
				for cbChan, _ := range cbChans {
					*cbChan <- errors.New("closed the application")
				}
				return
			case h := <-m.setOpenHandlerChan:
				openHandler = h
			case respChan := <-m.getOpenHandlerChan:
				respChan <- openHandler
			case cb := <-m.addCBChan:
				cbChans[cb] = true
			case cb := <-m.rmCBChan:
				delete(cbChans, cb)
			}

		}

	}()

	return m
}

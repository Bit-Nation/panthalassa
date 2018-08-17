package dapp

import (
	"errors"

	validator "github.com/Bit-Nation/panthalassa/dapp/validator"
	log "github.com/ipfs/go-log"
	logger "github.com/op/go-logging"
	otto "github.com/robertkrimen/otto"
)

type Module struct {
	logger          *logger.Logger
	vm              *otto.Otto
	setRendererChan chan *otto.Value
	getRendererChan chan chan *otto.Value
	addCBChan       chan *chan resp
	rmCBChan        chan *chan resp
	closer          chan struct{}
}

var sysLog = log.Logger("renderer - message")

// register module function in the VM
// setOpenHandler must be called with a callback
// the callback that is passed to `setMessageRenderer`
// should accept two parameters:
// 1. The "data" will hold the data (object) passed into the open call
//    (will e.g. hold the message and the context)
// 2. The "callback" should be can be called with two parameters:
// 		1. an error
// 		2. the rendered layout
func (m *Module) Register(vm *otto.Otto) error {
	m.vm = vm
	return vm.Set("setMessageRenderer", func(call otto.FunctionCall) otto.Value {

		sysLog.Debug("set message renderer")

		// validate function call
		v := validator.New()
		v.Set(0, &validator.TypeFunction)
		if err := v.Validate(vm, call); err != nil {
			m.logger.Error(err.String())
			return *err
		}

		// set renderer
		fn := call.Argument(0)
		m.setRendererChan <- &fn

		return otto.Value{}
	})
}

type resp struct {
	layout string
	error  error
}

// payload can be an arbitrary set of key value pairs
// should contain the "message" and the "context" tho
func (m *Module) RenderMessage(payload string) (string, error) {

	// fetch renderer
	rendererChan := make(chan *otto.Value)
	m.getRendererChan <- rendererChan
	renderer := <-rendererChan

	// make sure an renderer has been set
	if renderer == nil || !renderer.IsDefined() {
		return "", errors.New("failed to render message - no open handler set")
	}

	// convert context to otto js object
	payloadObj, err := m.vm.Object("(" + payload + ")")
	if err != nil {
		return "", err
	}

	cbChan := make(chan resp, 1)
	m.addCBChan <- &cbChan

	go func() {

		// call the message renderer
		// @todo what happens if we call the callback twice?
		_, err = renderer.Call(*renderer, payloadObj, func(call otto.FunctionCall) otto.Value {

			// delete cb chan from stack when we are done
			defer func() {
				m.rmCBChan <- &cbChan
			}()

			// fetch params from the callback call
			err := call.Argument(0)
			layout := call.Argument(1)

			r := resp{}

			// if there is an error, set it in the response
			if !err.IsUndefined() && !err.IsNull() {
				r.error = errors.New(err.String())
			}

			// set the layout in the response
			if layout.IsObject() {
				r.layout, r.error = layout.ToString()
			}

			cbChan <- r

			return otto.Value{}

		})

		if err != nil {
			m.logger.Error(err.Error())
		}

	}()

	r := <-cbChan
	return r.layout, r.error
}

func (m *Module) Close() error {
	m.closer <- struct{}{}
	return nil
}

func New(l *logger.Logger) *Module {
	m := &Module{
		logger:          l,
		setRendererChan: make(chan *otto.Value),
		getRendererChan: make(chan chan *otto.Value),
		addCBChan:       make(chan *chan resp),
		rmCBChan:        make(chan *chan resp),
		closer:          make(chan struct{}),
	}

	go func() {

		renderer := new(otto.Value)
		cbChans := map[*chan resp]bool{}

		for {

			select {
			case <-m.closer:
				for cbChan, _ := range cbChans {
					*cbChan <- resp{
						error: errors.New("closed the application"),
					}
				}
				return
			case r := <-m.setRendererChan:
				renderer = r
			case respChan := <-m.getRendererChan:
				respChan <- renderer
			case addCB := <-m.addCBChan:
				cbChans[addCB] = true
			case rmCB := <-m.rmCBChan:
				delete(cbChans, rmCB)
			}

		}

	}()

	return m
}

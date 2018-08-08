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

// payload can be an arbitrary set of key value pairs
// should contain the "message" and the "context" tho
func (m *Module) RenderMessage(payload string) (string, error) {

	// fetch renderer
	rendererChan := make(chan *otto.Value)
	m.getRendererChan <- rendererChan
	renderer := <-rendererChan

	// make sure an renderer has been set
	if renderer == nil {
		return "", errors.New("failed to render message - no open handler set")
	}

	// convert context to otto js object
	payloadObj, err := m.vm.Object("(" + payload + ")")
	if err != nil {
		return "", err
	}

	type resp struct {
		layout string
		error  error
	}

	c := make(chan resp, 1)

	go func() {

		// call the message renderer
		// @todo what happens if we call the callback twice?
		_, err = renderer.Call(*renderer, payloadObj, func(call otto.FunctionCall) otto.Value {

			// fetch params from the callback call
			err := call.Argument(0)
			layout := call.Argument(1)

			r := resp{}

			// if there is an error, set it in the response
			if !err.IsUndefined() {
				r.error = errors.New(err.String())
			}

			// set the layout in the response
			if layout.IsString() {
				r.layout = layout.String()
			}

			c <- r

			return otto.Value{}

		})

		if err != nil {
			m.logger.Error(err.Error())
		}

	}()

	r := <-c

	return r.layout, r.error
}

func (m *Module) Close() error {
	close(m.getRendererChan)
	close(m.setRendererChan)
	return nil
}

func New(l *logger.Logger) *Module {
	m := &Module{
		logger:          l,
		setRendererChan: make(chan *otto.Value),
		getRendererChan: make(chan chan *otto.Value),
	}

	go func() {

		renderer := new(otto.Value)

		for {

			// exit if channels got closed
			if m.setRendererChan == nil || m.getRendererChan == nil {
				return
			}

			select {
			case r := <-m.setRendererChan:
				renderer = r
			case respChan := <-m.getRendererChan:
				respChan <- renderer
			}

		}

	}()

	return m
}

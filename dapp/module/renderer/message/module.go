package dapp

import (
	"errors"

	validator "github.com/Bit-Nation/panthalassa/dapp/validator"
	log "github.com/ipfs/go-log"
	logger "github.com/op/go-logging"
	duktape "gopkg.in/olebedev/go-duktape.v3"
)

type Module struct {
	logger          *logger.Logger
	vm              *duktape.Context
	setRendererChan chan *duktape.Context
	getRendererChan chan chan *duktape.Context
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
func (m *Module) Register(vm *duktape.Context) error {
	m.vm = vm
	_, err := vm.PushGlobalGoFunction("setMessageRenderer", func(context *duktape.Context) int {
		sysLog.Debug("set message renderer")

		// validate function call
		v := validator.New()
		v.Set(0, &validator.TypeFunction)
		if err := v.Validate(context); err != nil {
			m.logger.Error(err.Error())
			return 1
		}

		// set renderer
		m.setRendererChan <- context

		return 0
	})
	return err
}

type resp struct {
	layout string
	error  error
}

// payload can be an arbitrary set of key value pairs
// should contain the "message" and the "context" tho
func (m *Module) RenderMessage(payload string) (string, error) {
	// fetch renderer
	rendererChan := make(chan *duktape.Context)
	m.getRendererChan <- rendererChan
	renderer := <-rendererChan

	// make sure an renderer has been set
	if renderer == nil {
		return "", errors.New("failed to render message - no open handler set")
	}

	// convert context to otto js object
	payloadObj := "(" + payload + ")"

	cbChan := make(chan resp, 1)
	m.addCBChan <- &cbChan

	go func() {

		// call the message renderer
		// @todo what happens if we call the callback twice?
		_, err := renderer.PushGlobalGoFunction("callbackRenderMessage", func(context *duktape.Context) int {

			// delete cb chan from stack when we are done
			defer func() {
				m.rmCBChan <- &cbChan
			}()

			// fetch params from the callback call
			r := resp{}

			// if there is an error, set it in the response
			if !context.IsUndefined(0) && !context.IsNull(0) {
				callbackerr := context.SafeToString(0)
				r.error = errors.New(callbackerr)
			}

			// set the layout in the response
			if context.IsString(1) {
				layout := context.SafeToString(1)
				r.layout = layout
			}

			cbChan <- r

			return 0

		})
		if err != nil {
			m.logger.Error(err.Error())
		}
		err = renderer.PevalString(payloadObj)
		if err != nil {
			m.logger.Error(err.Error())
		}
		err = renderer.PevalString(`callbackRenderMessage`)
		if err != nil {
			m.logger.Error(err.Error())
		}

		renderer.Call(2)
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
		setRendererChan: make(chan *duktape.Context),
		getRendererChan: make(chan chan *duktape.Context),
		addCBChan:       make(chan *chan resp),
		rmCBChan:        make(chan *chan resp),
		closer:          make(chan struct{}),
	}

	go func() {

		renderer := new(duktape.Context)
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

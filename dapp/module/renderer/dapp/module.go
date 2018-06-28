package dapp

import (
	"errors"
	"sync"

	validator "github.com/Bit-Nation/panthalassa/dapp/validator"
	log "github.com/op/go-logging"
	otto "github.com/robertkrimen/otto"
)

type Module struct {
	lock     sync.Mutex
	logger   *log.Logger
	renderer *otto.Value
	vm       *otto.Otto
}

func (m *Module) Name() string {
	return "RENDERER:DAPP"
}

func (m *Module) Register(vm *otto.Otto) error {
	m.vm = vm
	// setOpenHandler must be called with a callback
	// the callback that is passed to `setOpenHandler`
	// will be called with an "data" object and a callback
	// the callback should be called (with an optional error)
	// in order to return from the function
	return vm.Set("setOpenHandler", func(call otto.FunctionCall) otto.Value {

		// validate function call
		v := validator.New()
		v.Set(0, &validator.TypeFunction)
		if err := v.Validate(vm, call); err != nil {
			m.logger.Error(err.String())
			return *err
		}

		// set renderer
		m.lock.Lock()
		fn := call.Argument(0)
		m.renderer = &fn
		m.lock.Unlock()

		return otto.Value{}
	})
}

// payload can be an arbitrary set of key value pairs (as a json string)
func (m *Module) OpenDApp(payload string) error {

	// lock
	m.lock.Lock()
	defer m.lock.Unlock()

	// make sure an renderer has been set
	if m.renderer == nil {
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
		_, err = m.renderer.Call(*m.renderer, dataObj, func(call otto.FunctionCall) otto.Value {

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

func New(l *log.Logger) *Module {
	return &Module{
		lock:   sync.Mutex{},
		logger: l,
	}
}

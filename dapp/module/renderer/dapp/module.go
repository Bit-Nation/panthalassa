package dapp

import (
	"errors"
	"fmt"
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

// register module function in the VM
// setOpenHandler must be called with a callback
// the callback that is passed to `setOpenHandler`
// The "callback" should be can be called with an error if there is one
func (m *Module) Register(vm *otto.Otto) error {
	m.vm = vm
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

func (m *Module) OpenDApp(context string) error {

	// lock
	m.lock.Lock()
	defer m.lock.Unlock()

	// make sure an renderer has been set
	if m.renderer == nil {
		return errors.New("failed to open DApp - no open handler set")
	}

	// convert context to otto js object
	ctxObj, err := m.vm.Object(fmt.Sprintf(`(%s)`, context))
	if err != nil {
		return err
	}

	c := make(chan error, 1)

	// call the callback
	_, err = m.renderer.Call(*m.renderer, ctxObj, func(call otto.FunctionCall) otto.Value {

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
		return err
	}

	return <-c
}

func New(l *log.Logger) *Module {
	return &Module{
		lock:   sync.Mutex{},
		logger: l,
	}
}

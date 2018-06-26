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
	return "RENDERER:MESSAGE"
}

// register module function in the VM
// setOpenHandler must be called with a callback
// the callback that is passed to `setOpenHandler`
// should accept two parameters:
// 1. The "context" will hold the context in which the DApp is opened
// 2. The "callback" should be can be called with two parameters:
// 		1. an error
// 		2. the rendered layout
func (m *Module) Register(vm *otto.Otto) error {
	m.vm = vm
	return vm.Set("setMessageHandler", func(call otto.FunctionCall) otto.Value {

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

func (m *Module) RenderMessage(message string, context string) (string, error) {

	// lock
	m.lock.Lock()
	defer m.lock.Unlock()

	// make sure an renderer has been set
	if m.renderer == nil {
		return "", errors.New("failed to render message - no open handler set")
	}

	// convert context to otto js object
	ctxObj, err := m.vm.Object(fmt.Sprintf(`(%s)`, context))
	if err != nil {
		return "", err
	}
	msgObj, err := m.vm.Object(fmt.Sprintf(`(%s)`, message))
	if err != nil {
		return "", err
	}

	type resp struct {
		layout string
		error  error
	}

	c := make(chan resp, 1)

	// call the callback
	_, err = m.renderer.Call(*m.renderer, msgObj, ctxObj, func(call otto.FunctionCall) otto.Value {

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
		return "", err
	}

	r := <-c

	return r.layout, r.error
}

func New(l *log.Logger) *Module {
	return &Module{
		lock:   sync.Mutex{},
		logger: l,
	}
}

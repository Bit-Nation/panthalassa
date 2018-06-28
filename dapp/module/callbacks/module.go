package callbacks

import (
	"errors"
	"fmt"
	"sync"

	validator "github.com/Bit-Nation/panthalassa/dapp/validator"
	logger "github.com/op/go-logging"
	otto "github.com/robertkrimen/otto"
)

// with this module it's possible to register functions
// from inside of the vm and call them by there id

func New(l *logger.Logger) *Module {
	return &Module{
		lock:      sync.Mutex{},
		functions: map[uint]*otto.Value{},
		logger:    l,
	}
}

type Module struct {
	logger    *logger.Logger
	lock      sync.Mutex
	functions map[uint]*otto.Value
	counter   uint
	vm        *otto.Otto
}

func (m *Module) Name() string {
	return "CALLBACKS"
}

func (m *Module) CallFunction(id uint, args string) error {

	// lock
	m.lock.Lock()
	defer m.lock.Unlock()

	// check if function is registered
	fn := *m.functions[id]
	if !fn.IsFunction() {
		return errors.New(fmt.Sprintf("function with id: %d does not exist", id))
	}

	// parse params
	objArgs, err := m.vm.Object(fmt.Sprintf(`(%s)`, args))
	if err != nil {
		return err
	}

	done := make(chan error, 1)

	alreadyCalled := false

	fn.Call(fn, objArgs, func(call otto.FunctionCall) otto.Value {

		// check if callback has already been called
		if alreadyCalled {
			m.logger.Error("Already called callback")
			return m.vm.MakeCustomError("Callback", "Already called callback")
		}
		alreadyCalled = true

		// check parameters
		err := call.Argument(0)

		if !err.IsUndefined() {
			done <- errors.New(err.String())
			return otto.Value{}
		}

		done <- nil
		return otto.Value{}

	})

	return <-done

}

func (m *Module) Register(vm *otto.Otto) error {
	m.vm = vm
	return vm.Set("registerFunction", func(call otto.FunctionCall) otto.Value {
		// validate function call
		v := validator.New()
		v.Set(0, &validator.TypeFunction)
		if err := v.Validate(vm, call); err != nil {
			vm.Run(fmt.Sprintf(`throw new Error("registerFunction needs a callback as it's first param' %s")`, err.String()))
			return *err
		}

		// lock
		m.lock.Lock()
		defer m.lock.Unlock()

		// add function
		m.counter++
		fn := call.Argument(0)
		if _, exist := m.functions[m.counter]; exist {
			vm.Run(fmt.Sprintf(`throw new Error("Failed to register function. Id (%d) already in use")`, m.counter))
			return otto.Value{}
		}

		m.functions[m.counter] = &fn

		functionId, err := otto.ToValue(m.counter)
		if err != nil {
			m.logger.Error(err.Error())
		}

		return functionId

	})
}

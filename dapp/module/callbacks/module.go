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

// this will call the given function (identified by the id)
// with the given payload as an object and a callback
// e.g. myRegisteredFunction(payloadObj, cb)
// the callback must be called in order to "return" from the function
func (m *Module) CallFunction(id uint, payload string) error {

	// lock
	m.lock.Lock()
	defer m.lock.Unlock()

	// check if function is registered
	fn := *m.functions[id]
	if !fn.IsFunction() {
		return errors.New(fmt.Sprintf("function with id: %d does not exist", id))
	}

	// parse params
	objArgs, err := m.vm.Object(fmt.Sprintf(`(%s)`, payload))
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

// registerFunction will take a function as it's first and only parameter
// if the parameter is not a function it will throw an error
// a ID (uint) is returned that represents the id of the registered function
// a registered function will be called with an object containing information
// and a callback that should be called (with an optional error) in order to
// "return"
func (m *Module) Register(vm *otto.Otto) error {
	m.vm = vm
	return vm.Set("registerFunction", func(call otto.FunctionCall) otto.Value {
		// validate function call
		v := validator.New()
		v.Set(0, &validator.TypeFunction)
		if err := v.Validate(vm, call); err != nil {
			vm.Run(fmt.Sprintf(`throw new Error("registerFunction needs a callback as it's first param %s")`, err.String()))
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

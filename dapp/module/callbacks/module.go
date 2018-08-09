package callbacks

import (
	"errors"
	"fmt"

	reqLim "github.com/Bit-Nation/panthalassa/dapp/request_limitation"
	validator "github.com/Bit-Nation/panthalassa/dapp/validator"
	log "github.com/ipfs/go-log"
	logger "github.com/op/go-logging"
	otto "github.com/robertkrimen/otto"
)

var debugger = log.Logger("callbacks")

// with this module it's possible to register functions
// from inside of the vm and call them by there id

func New(l *logger.Logger) *Module {

	m := &Module{
		logger:             l,
		reqLim:             reqLim.NewCount(10000, errors.New("can't register more functions")),
		addFunctionChan:    make(chan addFunction),
		fetchFunctionChan:  make(chan fetchFunction),
		deleteFunctionChan: make(chan uint),
		addCBChan:          make(chan *chan error),
		rmCBChan:           make(chan *chan error),
		nextCBChan:         make(chan chan *chan error),
	}

	// state machine
	go func() {

		var count uint

		functions := map[uint]*otto.Value{}
		cbChans := map[*chan error]bool{}

		for {

			// exit
			if m.deleteFunctionChan == nil || m.addFunctionChan == nil || m.fetchFunctionChan == nil {
				return
			}

			select {
			// add function
			case addFn := <-m.addFunctionChan:
				if addFn.respChan == nil || addFn.fn == nil {
					continue
				}
				// exit if function exist
				if _, exist := functions[count+1]; exist {
					addFn.respChan <- struct {
						id    uint
						error error
					}{id: 0, error: fmt.Errorf("function for id: %d already registered", count+1)}
					continue
				}

				// inc request limitation
				err := m.reqLim.Increase()
				if err != nil {
					addFn.respChan <- struct {
						id    uint
						error error
					}{id: 0, error: err}
					continue
				}

				// inc counter
				count++

				// register function
				functions[count] = addFn.fn
				addFn.respChan <- struct {
					id    uint
					error error
				}{id: count, error: nil}

			// remove function
			case remFn := <-m.deleteFunctionChan:
				fmt.Println(remFn)
				_, exist := functions[remFn]
				if !exist {
					continue
				}
				m.reqLim.Decrease()
				delete(functions, remFn)
			// fetch function
			case fetchFn := <-m.fetchFunctionChan:
				// exit if resp chan
				if fetchFn.respChan == nil {
					continue
				}
				fn, exist := functions[fetchFn.id]
				if !exist {
					fetchFn.respChan <- nil
					continue
				}
				fetchFn.respChan <- fn
			// add callback channel
			case addCB := <-m.addCBChan:
				cbChans[addCB] = true
			// remove callback channel
			case rmCB := <-m.rmCBChan:
				delete(cbChans, rmCB)
			// next callback channel
			case nextCBChan := <-m.nextCBChan:
				if len(cbChans) == 0 {
					nextCBChan <- nil
					continue
				}
				for cbChan, _ := range cbChans {
					nextCBChan <- cbChan
					break
				}
			}

		}

	}()

	return m
}

type addFunction struct {
	fn       *otto.Value
	respChan chan struct {
		id    uint
		error error
	}
}

type fetchFunction struct {
	id       uint
	respChan chan *otto.Value
}

type Module struct {
	logger             *logger.Logger
	vm                 *otto.Otto
	reqLim             *reqLim.Count
	addFunctionChan    chan addFunction
	fetchFunctionChan  chan fetchFunction
	deleteFunctionChan chan uint
	addCBChan          chan *chan error
	rmCBChan           chan *chan error
	// returns a cb chan from the stack
	nextCBChan chan chan *chan error
}

func (m *Module) Close() error {
	close(m.addFunctionChan)
	close(m.fetchFunctionChan)
	close(m.deleteFunctionChan)
	return nil
}

// this will call the given function (identified by the id)
// with the given payload as an object and a callback
// e.g. myRegisteredFunction(payloadObj, cb)
// the callback must be called in order to "return" from the function
func (m *Module) CallFunction(id uint, payload string) error {

	debugger.Debug(fmt.Errorf("call function with id: %d and payload: %s", id, payload))

	respChan := make(chan *otto.Value)
	m.fetchFunctionChan <- fetchFunction{
		id:       id,
		respChan: respChan,
	}
	fn := <-respChan

	// check if function is registered
	if !fn.IsFunction() {
		return errors.New(fmt.Sprintf("function with id: %d does not exist", id))
	}

	// parse params
	objArgs, err := m.vm.Object("(" + payload + ")")
	if err != nil {
		return err
	}

	done := make(chan error, 1)
	m.addCBChan <- &done

	alreadyCalled := false

	fn.Call(*fn, objArgs, func(call otto.FunctionCall) otto.Value {

		defer func() {
			m.rmCBChan <- &done
		}()

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
	err := vm.Set("registerFunction", func(call otto.FunctionCall) otto.Value {
		// validate function call
		v := validator.New()
		v.Set(0, &validator.TypeFunction)
		if err := v.Validate(vm, call); err != nil {
			vm.Run(fmt.Sprintf(`throw new Error("registerFunction needs a callback as it's first param %s")`, err.String()))
			return *err
		}

		// callback
		fn := call.Argument(0)

		// add function to stack
		idRespChan := make(chan struct {
			id    uint
			error error
		}, 1)
		m.addFunctionChan <- addFunction{
			respChan: idRespChan,
			fn:       &fn,
		}
		// response
		addResp := <-idRespChan

		// exit vm on error
		if addResp.error != nil {
			vm.Run(fmt.Sprintf(`throw new Error(%s)`, addResp.error.Error()))
			return otto.Value{}
		}

		// convert function id to otto value
		functionId, err := otto.ToValue(addResp.id)
		if err != nil {
			m.logger.Error(err.Error())
		}

		return functionId

	})
	if err != nil {
		return err
	}

	return vm.Set("unRegisterFunction", func(call otto.FunctionCall) otto.Value {

		// validate function call
		v := validator.New()
		v.Set(0, &validator.TypeNumber)
		if err := v.Validate(vm, call); err != nil {
			return *err
		}

		// function id
		funcID := call.Argument(0)
		id, err := funcID.ToInteger()
		if err != nil {
			m.logger.Error(err.Error())
			return otto.Value{}
		}
		fmt.Println("delete function")
		// delete function from channel
		m.deleteFunctionChan <- uint(id)

		return otto.Value{}
	})

}

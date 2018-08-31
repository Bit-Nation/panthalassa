package callbacks

import (
	"errors"
	"fmt"

	reqLim "github.com/Bit-Nation/panthalassa/dapp/request_limitation"
	validator "github.com/Bit-Nation/panthalassa/dapp/validator"
	log "github.com/ipfs/go-log"
	logger "github.com/op/go-logging"
	duktape "gopkg.in/olebedev/go-duktape.v3"
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
		closer:             make(chan struct{}),
	}

	// state machine
	go func() {

		var count uint

		functions := map[uint]*duktape.Context{}
		cbChans := map[*chan error]bool{}

		for {

			select {
			// exit from go routine
			case <-m.closer:
				// close all callback's
				for cbChan, _ := range cbChans {
					*cbChan <- errors.New("closed application")
				}
				return
			// add function
			case addFn := <-m.addFunctionChan:
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
				_, exist := functions[remFn]
				if !exist {
					continue
				}
				m.reqLim.Decrease()
				delete(functions, remFn)
			// fetch function
			case fetchFn := <-m.fetchFunctionChan:
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
			}

		}

	}()

	return m
}

type addFunction struct {
	fn       *duktape.Context
	respChan chan struct {
		id    uint
		error error
	}
}

type fetchFunction struct {
	id       uint
	respChan chan *duktape.Context
}

type Module struct {
	logger             *logger.Logger
	vm                 *duktape.Context
	reqLim             *reqLim.Count
	addFunctionChan    chan addFunction
	fetchFunctionChan  chan fetchFunction
	deleteFunctionChan chan uint
	addCBChan          chan *chan error
	rmCBChan           chan *chan error
	// returns a cb chan from the stack
	nextCBChan chan chan *chan error
	closer     chan struct{}
}

func (m *Module) Close() error {
	m.closer <- struct{}{}
	return nil
}

// this will call the given function (identified by the id)
// with the given payload as an object and a callback
// e.g. myRegisteredFunction(payloadObj, cb)
// the callback must be called in order to "return" from the function
func (m *Module) CallFunction(id uint, payload string) error {
	debugger.Debug(fmt.Errorf("call function with id: %d and payload: %s", id, payload))

	respChan := make(chan *duktape.Context)
	m.fetchFunctionChan <- fetchFunction{
		id:       id,
		respChan: respChan,
	}
	vm := <-respChan

	if vm == nil || vm.GetType(0).IsNone() {
		return errors.New(fmt.Sprintf("function with id: %d does not exist", id))
	}

	// check if function is registered
	if !vm.IsFunction(0) {
		return errors.New(fmt.Sprintf("function with id: %d does not exist", id))
	}

	// parse params
	objArgs := "(" + payload + ")"

	done := make(chan error, 1)
	m.addCBChan <- &done

	alreadyCalled := false

	_, err := vm.PushGlobalGoFunction("callbackCallFunction", func(context *duktape.Context) int {

		defer func() {
			m.rmCBChan <- &done
		}()

		// check if callback has already been called
		if alreadyCalled {
			m.logger.Error("Already called callback")
			if context.IsFunction(1) {
				context.PushString("Callback: Already called callback")
				context.Call(1)
			}
			return 1
		}
		alreadyCalled = true

		if !context.IsUndefined(0) {
			firstParameter := context.ToString(0)
			if context.IsFunction(1) {
				context.PushString(firstParameter)
				context.Call(1)
			}
			done <- errors.New(firstParameter)
			return 1
		}

		done <- nil
		return 0

	})
	if err != nil {
		m.logger.Error(err.Error())
	}
	err = vm.PevalString(objArgs)
	if err != nil {
		m.logger.Error(err.Error())
	}
	err = vm.PevalString(`callbackCallFunction`)
	if err != nil {
		m.logger.Error(err.Error())
	}
	vm.Call(2)
	return <-done

}

// registerFunction will take a function as it's first and only parameter
// if the parameter is not a function it will throw an error
// a ID (uint) is returned that represents the id of the registered function
// a registered function will be called with an object containing information
// and a callback that should be called (with an optional error) in order to
// "return"
func (m *Module) Register(vm *duktape.Context) error {
	m.vm = vm
	_, err := vm.PushGlobalGoFunction("registerFunction", func(context *duktape.Context) int {
		// validate function call
		v := validator.New()
		v.Set(0, &validator.TypeFunction)
		if err := v.Validate(context); err != nil {
			m.logger.Error(fmt.Sprintf(`registerFunction needs a callback as it's first param %s`, err.Error()))
			return 1
		}

		// add function to stack
		idRespChan := make(chan struct {
			id    uint
			error error
		}, 1)
		m.addFunctionChan <- addFunction{
			respChan: idRespChan,
			fn:       context,
		}
		// response
		addResp := <-idRespChan

		// exit vm on error
		if addResp.error != nil {
			m.logger.Error(addResp.error.Error())
			return 1
		}

		// convert function id to int
		functionId := int(addResp.id)

		return functionId

	})
	if err != nil {
		return err
	}

	_, err = vm.PushGlobalGoFunction("unRegisterFunction", func(context *duktape.Context) int {

		// validate function call
		v := validator.New()
		v.Set(0, &validator.TypeNumber)
		if err := v.Validate(context); err != nil {
			return 1
		}

		// function id
		funcID := context.ToInt(0)

		// delete function from channel
		m.deleteFunctionChan <- uint(funcID)

		return 0
	})
	return err
}

package uuidv4

import (
	otto "github.com/robertkrimen/otto"
	uuid "github.com/satori/go.uuid"
)

var newUuid = uuid.NewV4

// register an module on the given
// vm that allows to generate uuid's
// of version 4. it expects a
// callback as it's only argument
type UUIDV4 struct{}

func (r *UUIDV4) Name() string {
	return "UUIDV4"
}

func (r *UUIDV4) Register(vm *otto.Otto) error {

	return vm.Set("uuidV4", func(call otto.FunctionCall) otto.Value {

		// make sure callback is a function
		cb := call.Argument(0)
		if !cb.IsFunction() {
			_, err := cb.Call(call.This, nil, "expected first argument to be a function")
			if err != nil {
				// @todo process error
			}
			return otto.Value{}
		}

		// create uuid
		id, err := newUuid()
		if err != nil {
			_, err = cb.Call(call.This, nil, err.Error())
			return otto.Value{}
		}

		// call callback with uuid
		_, err = cb.Call(call.This, id.String(), nil)
		if err != nil {
			// @todo process error
		}
		return otto.Value{}

	})

}

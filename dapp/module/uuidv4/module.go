package uuidv4

import (
	validator "github.com/Bit-Nation/panthalassa/dapp/validator"
	log "github.com/ipfs/go-log"
	logger "github.com/op/go-logging"
	otto "github.com/robertkrimen/otto"
	uuid "github.com/satori/go.uuid"
)

var newUuid = uuid.NewV4

var sysLogger = log.Logger("uuidv4")

// register an module on the given
// vm that allows to generate uuid's
// of version 4. it expects a
// callback as it's only argument
type UUIDV4 struct {
	logger *logger.Logger
}

func New(l *logger.Logger) *UUIDV4 {
	return &UUIDV4{logger: l}
}

func (r *UUIDV4) Register(vm *otto.Otto) error {

	return vm.Set("uuidV4", func(call otto.FunctionCall) otto.Value {

		sysLogger.Debug("generate uuidv4")

		// validate function call
		v := validator.New()
		v.Set(0, &validator.TypeFunction)
		if err := v.Validate(vm, call); err != nil {
			vm.Run(`throw new Error("uuidV4 expects an callback as it's parameter")`)
			return otto.Value{}
		}

		cb := call.Argument(0)

		// create uuid
		id, err := newUuid()

		// call callback with error
		if err != nil {
			_, err = cb.Call(call.This, err.Error())
			if err != nil {
				r.logger.Error(err.Error())
			}
			return otto.Value{}
		}

		// call callback with uuid
		_, err = cb.Call(call.This, nil, id.String())
		if err != nil {
			r.logger.Error(err.Error())
		}

		return otto.Value{}

	})

}

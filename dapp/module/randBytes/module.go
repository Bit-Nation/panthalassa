package randBytes

import (
	"crypto/rand"

	validator "github.com/Bit-Nation/panthalassa/dapp/validator"
	log "github.com/ipfs/go-log"
	logger "github.com/op/go-logging"
	otto "github.com/robertkrimen/otto"
)

var randSource = rand.Reader

var sysLog = log.Logger("rand bytes")

func New(l *logger.Logger) *Module {
	return &Module{
		logger: l,
	}
}

type Module struct {
	logger *logger.Logger
}

func (m *Module) Register(vm *otto.Otto) error {
	return vm.Set("randomBytes", func(call otto.FunctionCall) otto.Value {

		sysLog.Debug("generate random bytes")

		// validate call
		v := validator.New()
		v.Set(0, &validator.TypeNumber)
		v.Set(1, &validator.TypeFunction)
		if err := v.Validate(vm, call); err != nil {
			m.logger.Error(err.String())
		}

		cb := call.Argument(1)

		// convert to integer
		amount, err := call.Argument(0).ToInteger()
		if err != nil {
			_, err := cb.Call(cb, err.Error())
			if err != nil {
				m.logger.Error(err.Error())
			}
			return otto.Value{}
		}

		// read random bytes
		destination := make([]byte, amount)
		_, err = randSource.Read(destination)
		if err != nil {
			_, err := cb.Call(cb, err.Error())
			if err != nil {
				m.logger.Error(err.Error())
			}
			return otto.Value{}
		}

		// call callback
		_, err = cb.Call(cb, nil, destination)
		if err != nil {
			_, err := cb.Call(cb, err.Error())
			if err != nil {
				m.logger.Error(err.Error())
			}
			return otto.Value{}
		}

		return otto.Value{}

	})
}

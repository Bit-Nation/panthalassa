package sendEthTx

import (
	validator "github.com/Bit-Nation/panthalassa/dapp/validator"
	log "github.com/op/go-logging"
	otto "github.com/robertkrimen/otto"
)

type SendEthereumTransaction interface {
	// will return a JSON serialized transaction
	// and an error if there is one
	Send(value, to, data string) (string, error)
}

func New(ethApi SendEthereumTransaction, l *log.Logger) *Module {
	return &Module{
		logger: l,
		ethApi: ethApi,
	}
}

type Module struct {
	logger *log.Logger
	ethApi SendEthereumTransaction
}

func (m *Module) Name() string {
	return "SEND_ETH_TX"
}

func (m *Module) Register(vm *otto.Otto) error {

	return vm.Set("sendETHTransaction", func(call otto.FunctionCall) otto.Value {

		// validate function call
		v := validator.New()
		v.Set(0, &validator.TypeObject)
		v.Set(1, &validator.TypeFunction)
		if err := v.Validate(vm, call); err != nil {
			m.logger.Error(err.String())
			return *err
		}

		cb := call.Argument(1)

		// validate transaction object
		obj := call.Argument(0).Object()
		objVali := validator.NewObjValidator()
		objVali.Set("value", validator.ObjTypeString, true)
		objVali.Set("to", validator.ObjTypeAddress, true)
		objVali.Set("data", validator.ObjTypeString, true)
		if err := objVali.Validate(vm, *obj); err != nil {
			cb.Call(cb, err.String())
			return otto.Value{}
		}

		// make the request async
		go func() {

			to, err := obj.Get("to")
			if err != nil {
				cb.Call(cb, err.Error())
				return
			}

			value, err := obj.Get("value")
			if err != nil {
				cb.Call(cb, err.Error())
				return
			}

			data, err := obj.Get("data")
			if err != nil {
				cb.Call(cb, err.Error())
				return
			}

			// try to sign a transaction
			tx, err := m.ethApi.Send(
				value.String(),
				to.String(),
				data.String(),
			)

			// exit on error
			if err != nil {
				cb.Call(cb, err.Error())
				return
			}

			// call callback with transaction
			cb.Call(cb, nil, tx)

		}()

		return otto.Value{}

	})

}

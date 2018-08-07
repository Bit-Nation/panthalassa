package sendEthTx

import (
	"errors"
	"time"

	reqLim "github.com/Bit-Nation/panthalassa/dapp/request_limitation"
	validator "github.com/Bit-Nation/panthalassa/dapp/validator"
	log "github.com/ipfs/go-log"
	logger "github.com/op/go-logging"
	otto "github.com/robertkrimen/otto"
)

var sysLogger = log.Logger("send eth tx")

type SendEthereumTransaction interface {
	// will return a JSON serialized transaction
	// and an error if there is one
	SendEthereumTransaction(value, to, data string) (string, error)
}

func New(ethApi SendEthereumTransaction, l *logger.Logger) *Module {
	return &Module{
		logger:     l,
		ethApi:     ethApi,
		throttling: reqLim.NewThrottling(6, time.Minute, 50, errors.New("can't add more signing requests to stack")),
	}
}

type Module struct {
	logger     *logger.Logger
	ethApi     SendEthereumTransaction
	throttling *reqLim.Throttling
}

func (m *Module) Register(vm *otto.Otto) error {

	// send an ethereum transaction
	// musst be called with an object that holds value, to and data
	return vm.Set("sendETHTransaction", func(call otto.FunctionCall) otto.Value {

		sysLogger.Debug("send eth transaction")

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

		// execute in the context of the throttling
		// request limitation
		m.throttling.Exec(func() {

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
			tx, err := m.ethApi.SendEthereumTransaction(
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

		})

		return otto.Value{}

	})

}

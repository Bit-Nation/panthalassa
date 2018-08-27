package sendEthTx

import (
	"errors"
	"time"

	reqLim "github.com/Bit-Nation/panthalassa/dapp/request_limitation"
	validator "github.com/Bit-Nation/panthalassa/dapp/validator"
	log "github.com/ipfs/go-log"
	logger "github.com/op/go-logging"
	duktape "gopkg.in/olebedev/go-duktape.v3"
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

func (m *Module) Close() error {
	return nil
}

func (m *Module) Register(vm *duktape.Context) error {

	// send an ethereum transaction
	// musst be called with an object that holds value, to and data
	_, err := vm.PushGlobalGoFunction("sendETHTransaction", func(context *duktape.Context) int {
		var itemsToPopBeforeCallback int
		sysLogger.Debug("send eth transaction")

		// validate function call
		v := validator.New()
		v.Set(0, &validator.TypeObject)
		v.Set(1, &validator.TypeFunction)
		// utils to handle an occurred error
		handleError := func(errMsg string) int {
			if context.IsFunction(1) {
				context.PopN(itemsToPopBeforeCallback)
				context.PushString(errMsg)
				context.Call(1)
				return 0
			}
			m.logger.Error(errMsg)
			return 1
		}

		if err := v.Validate(context); err != nil {
			m.logger.Error(err.Error())
			return 1
		}

		// validate transaction object
		objVali := validator.NewObjValidator()
		objVali.Set("value", validator.ObjTypeString, true)
		objVali.Set("to", validator.ObjTypeAddress, true)
		objVali.Set("data", validator.ObjTypeString, true)
		if err := objVali.Validate(vm, 0); err != nil {
			handleError(err.Error())
			return 1
		}
		// execute in the context of the throttling
		// request limitation
		throttlingFunc := func() {

			if !context.GetPropString(0, "to") {
				err := errors.New(`key "to" doesn't exist`)
				handleError(err.Error())
				return
			}
			to := context.ToString(-1)
			itemsToPopBeforeCallback++
			itemsToPopBeforeCallback++

			if !context.GetPropString(0, "value") {
				err := errors.New(`key "value" doesn't exist`)
				handleError(err.Error())
				return
			}
			value := context.ToString(-1)
			itemsToPopBeforeCallback++
			itemsToPopBeforeCallback++

			if !context.GetPropString(0, "data") {
				err := errors.New(`key "data" doesn't exist`)
				handleError(err.Error())
				return
			}
			data := context.ToString(-1)
			itemsToPopBeforeCallback++
			itemsToPopBeforeCallback++

			// try to sign a transaction
			tx, err := m.ethApi.SendEthereumTransaction(
				value,
				to,
				data,
			)

			// exit on error
			if err != nil {
				handleError(err.Error())
			}

			// call callback with transaction
			//@TODO FIND OUT WHY DO WE NEED TO CLEAR THE STACK LIKE THIS?
			context.PopN(itemsToPopBeforeCallback)
			context.PushUndefined()
			context.PushString(tx)
			context.Call(2)
			return

		}
		throttlingFunc()
		//@TODO Find a way to fix throttling
		//m.throttling.Exec(throttlingFunc)

		return 0

	})
	return err
}

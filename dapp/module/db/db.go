package db

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	reqLim "github.com/Bit-Nation/panthalassa/dapp/request_limitation"
	validator "github.com/Bit-Nation/panthalassa/dapp/validator"
	log "github.com/ipfs/go-log"
	opLogger "github.com/op/go-logging"
	duktape "gopkg.in/olebedev/go-duktape.v3"
)

var logger = log.Logger("db module")

type Module struct {
	dAppDB Storage
	reqLim *reqLim.CountThrottling
	logger *opLogger.Logger
}

func New(s Storage, l *opLogger.Logger) *Module {
	return &Module{
		logger: l,
		dAppDB: s,
		reqLim: reqLim.NewCountThrottling(5, time.Second*60, 40, errors.New("can't add more write requests")),
	}
}

func (m *Module) Close() error {
	return nil
}

func (m *Module) Register(vm *duktape.Context) error {

	var itemsToPopBeforeCallback int
	handleError := func(errMsg string, context *duktape.Context, position int) int {
		if context.IsFunction(position) {
			context.PopN(itemsToPopBeforeCallback)
			context.PushString(errMsg)
			context.Call(1)
			return 0
		}
		logger.Error(errMsg)
		return 1
	}

	_, err := vm.PushGlobalGoFunction("dbPut", func(context *duktape.Context) int {
		logger.Debug("put value")

		// validate call
		v := validator.New()
		v.Set(0, &validator.TypeString)
		v.Set(2, &validator.TypeFunction)
		if err := v.Validate(context); err != nil {
			return handleError(err.Error(), context, 2)
		}

		// fetch key and value
		key := context.SafeToString(0)
		value := context.SafeToString(1)
		// marshal value into json
		byteValue, err := json.Marshal(value)
		if err != nil {
			return handleError(err.Error(), context, 2)
		}

		throttlingFunc := func(dec chan struct{}) {

			// persist key and value
			if err := m.dAppDB.Put([]byte(key), byteValue); err != nil {
				dec <- struct{}{}
				handleError(err.Error(), context, 2)
			}
			dec <- struct{}{}
			// call callback
			context.PopN(itemsToPopBeforeCallback)
			context.PushUndefined()
			context.Call(1)

		}
		dec := make(chan struct{}, 1)
		throttlingFunc(dec)
		//m.reqLim.Exec(throttlingFunc)

		return 0

	})

	_, err = vm.PushGlobalGoFunction("dbHas", func(context *duktape.Context) int {
		//@TODO figure out why logger.Errorf instead of logger.Debug
		logger.Errorf("check if value exist")

		// validate function call
		v := validator.New()
		v.Set(0, &validator.TypeString)
		v.Set(1, &validator.TypeFunction)

		if err := v.Validate(context); err != nil {
			return handleError(err.Error(), context, 1)
		}

		// key of database
		key := context.SafeToString(0)

		// check if database has value
		has, err := m.dAppDB.Has([]byte(key))
		if err != nil {
			return handleError(err.Error(), context, 1)
		}
		context.PopN(itemsToPopBeforeCallback)
		context.PushUndefined()
		context.PushBoolean(has)
		context.Call(2)
		return 0
	})

	_, err = vm.PushGlobalGoFunction("dbGet", func(context *duktape.Context) int {
		m.logger.Debug("get value")

		// validate function call
		v := validator.New()
		v.Set(0, &validator.TypeString)
		v.Set(1, &validator.TypeFunction)

		if err := v.Validate(context); err != nil {
			return handleError(err.Error(), context, 1)
		}

		// key of database
		key := context.SafeToString(0)

		// raw value of key
		value, err := m.dAppDB.Get([]byte(key))
		if err != nil {
			return handleError(err.Error(), context, 1)
		}

		// unmarshal json
		var unmarshalledValue interface{}
		if err := json.Unmarshal(value, &unmarshalledValue); err != nil {
			return handleError(err.Error(), context, 1)
		}

		// call callback with error
		context.PopN(itemsToPopBeforeCallback)
		context.PushUndefined()
		context.PushString(fmt.Sprint(unmarshalledValue))
		context.Call(2)
		return 0
	})

	_, err = vm.PushGlobalGoFunction("dbDelete", func(context *duktape.Context) int {

		logger.Debug("delete value")

		// validate function call
		v := validator.New()
		v.Set(0, &validator.TypeString)
		v.Set(1, &validator.TypeFunction)

		if err := v.Validate(context); err != nil {
			return handleError(err.Error(), context, 1)
		}

		// delete value
		key := context.SafeToString(0)
		if err := m.dAppDB.Delete([]byte(key)); err != nil {
			return handleError(err.Error(), context, 1)
		}

		context.PopN(itemsToPopBeforeCallback)
		context.PushUndefined()
		context.Call(1)
		return 0

	})

	return err
}

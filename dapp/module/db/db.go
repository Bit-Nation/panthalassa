package db

import (
	"encoding/json"
	"errors"
	"time"

	reqLim "github.com/Bit-Nation/panthalassa/dapp/request_limitation"
	validator "github.com/Bit-Nation/panthalassa/dapp/validator"
	log "github.com/ipfs/go-log"
	otto "github.com/robertkrimen/otto"
)

var logger = log.Logger("db module")

type Module struct {
	dAppDB Storage
	reqLim *reqLim.CountThrottling
}

func New(s Storage) *Module {
	return &Module{
		dAppDB: s,
		reqLim: reqLim.NewCountThrottling(5, time.Second*60, 40, errors.New("can't add more write requests")),
	}
}

func (m *Module) Close() error {
	return nil
}

func (m *Module) Register(vm *otto.Otto) error {

	handleError := func(errMsg string, cb otto.Value) otto.Value {
		if cb.IsFunction() {
			cb.Call(cb, errMsg)
			return otto.Value{}
		}
		logger.Error(errMsg)
		return otto.Value{}
	}

	return vm.Set("db", map[string]interface{}{
		"put": func(call otto.FunctionCall) otto.Value {

			logger.Debug("put value")

			// validate call
			v := validator.New()
			v.Set(0, &validator.TypeString)
			v.Set(2, &validator.TypeFunction)
			cb := call.Argument(2)
			if err := v.Validate(vm, call); err != nil {
				return handleError(err.String(), cb)
			}

			// fetch key and value
			key := call.Argument(0)
			value, err := call.Argument(1).Export()
			if err != nil {
				return handleError(err.Error(), cb)
			}

			// marshal value into json
			byteValue, err := json.Marshal(value)
			if err != nil {
				return handleError(err.Error(), cb)
			}

			m.reqLim.Exec(func(dec chan struct{}) {

				// persist key and value
				if err := m.dAppDB.Put([]byte(key.String()), byteValue); err != nil {
					dec <- struct{}{}
					handleError(err.Error(), cb)
				}
				dec <- struct{}{}
				// call callback
				_, err = cb.Call(cb)
				if err != nil {
					logger.Error(err.Error())
				}

			})

			return otto.Value{}

		},
		"has": func(call otto.FunctionCall) otto.Value {

			logger.Errorf("check if value exist")

			// validate function call
			v := validator.New()
			v.Set(0, &validator.TypeString)
			v.Set(1, &validator.TypeFunction)
			cb := call.Argument(1)
			if err := v.Validate(vm, call); err != nil {
				return handleError(err.String(), cb)
			}

			// key of database
			key := call.Argument(0).String()

			// check if database has value
			has, err := m.dAppDB.Has([]byte(key))
			if err != nil {
				return handleError(err.Error(), cb)
			}
			_, err = cb.Call(cb, nil, has)
			if err != nil {
				logger.Error(err.Error())
			}
			return otto.Value{}

		},
		"get": func(call otto.FunctionCall) otto.Value {

			logger.Debug("get value")

			// validate function call
			v := validator.New()
			v.Set(0, &validator.TypeString)
			v.Set(1, &validator.TypeFunction)
			cb := call.Argument(1)
			if err := v.Validate(vm, call); err != nil {
				return handleError(err.String(), cb)
			}

			// key of database
			key := call.Argument(0).String()

			// raw value of key
			value, err := m.dAppDB.Get([]byte(key))
			if err != nil {
				return handleError(err.Error(), cb)
			}

			// unmarshal json
			var unmarshalledValue interface{}
			if err := json.Unmarshal(value, &unmarshalledValue); err != nil {
				return handleError(err.Error(), cb)
			}

			// call callback with error
			_, err = cb.Call(cb, nil, unmarshalledValue)
			if err != nil {
				logger.Error(err.Error())
			}
			return otto.Value{}

		},
		"delete": func(call otto.FunctionCall) otto.Value {

			logger.Debug("delete value")

			// validate function call
			v := validator.New()
			v.Set(0, &validator.TypeString)
			v.Set(1, &validator.TypeFunction)
			cb := call.Argument(1)
			if err := v.Validate(vm, call); err != nil {
				return handleError(err.String(), cb)
			}

			// delete value
			key := call.Argument(0)
			if err := m.dAppDB.Delete([]byte(key.String())); err != nil {
				return handleError(err.Error(), cb)
			}

			if _, err := cb.Call(cb); err != nil {
				logger.Error(err.Error())
			}

			return otto.Value{}

		},
	})
}

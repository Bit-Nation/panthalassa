package ethWebSocket

import (
	"encoding/json"

	validator "github.com/Bit-Nation/panthalassa/dapp/validator"
	ethws "github.com/Bit-Nation/panthalassa/ethws"
	logger "github.com/op/go-logging"
	otto "github.com/robertkrimen/otto"
)

type EthWS struct {
	logger *logger.Logger
	ethWS  *ethws.EthereumWS
}

func New(logger *logger.Logger, ethWS *ethws.EthereumWS) *EthWS {
	return &EthWS{
		logger: logger,
		ethWS:  ethWS,
	}
}

func (ws *EthWS) Name() string {
	return "ETHEREUM_WEB_SOCKET"
}

func (ws *EthWS) Register(vm *otto.Otto) error {

	return vm.Set("ethereumRequest", func(call otto.FunctionCall) otto.Value {

		// validate function call
		v := validator.New()
		v.Set(0, &validator.TypeString)
		v.Set(1, &validator.TypeFunction)
		if err := v.Validate(vm, call); err != nil {
			ws.logger.Error(err.String())
			return *err
		}

		go func() {

			cb := call.Argument(1)

			// unmarshal json rpc request
			var r ethws.Request
			if err := json.Unmarshal([]byte(call.Argument(0).String()), &r); err != nil {
				_, err := cb.Call(cb, nil, err.Error())
				if err != nil {
					ws.logger.Error(err.Error())
				}
				return
			}

			// send json rpc request
			respCha, err := ws.ethWS.SendRequest(r)
			if err != nil {
				_, err := cb.Call(cb, nil, err.Error())
				if err != nil {
					ws.logger.Error(err.Error())
				}
				return
			}

			// unmarshal response
			rawResponse, err := json.Marshal(<-respCha)
			if err != nil {
				_, err := cb.Call(cb, nil, err.Error())
				ws.logger.Error(err.Error())
				return
			}

			if _, err := cb.Call(cb, string(rawResponse)); err != nil {
				ws.logger.Error(err.Error())
			}

		}()

		return otto.Value{}

	})

}

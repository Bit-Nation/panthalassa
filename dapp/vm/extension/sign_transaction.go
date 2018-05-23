package extension

import (
	"encoding/json"

	deviceApi "github.com/Bit-Nation/panthalassa/api/device"
	ethCommon "github.com/ethereum/go-ethereum/common"
	otto "github.com/robertkrimen/otto"
)

type SignTransactionRpcCall struct {
	Reason  string `json:"reason"`
	TxTo    string `json:"to"`
	TxValue string `json:"value"`
	TxData  string `json:"data"`
}

func (c SignTransactionRpcCall) Type() string {
	return "VM:TRANSACTION:SIGN"
}

func (c *SignTransactionRpcCall) Data() (string, error) {

	raw, err := json.Marshal(c)
	return string(raw), err

}

func (c *SignTransactionRpcCall) Valid() error {
	return nil
}

// a sign call requires
// @todo there is need for validating the parameter length. E.g the data should not be longer than x char's etc
func SignTransaction(api *deviceApi.Api) OttoFunction {

	return func(call otto.FunctionCall) otto.Value {

		toSign := call.Argument(0).Object()
		callback := call.Argument(1)

		if toSign == nil {
			v, _ := otto.ToValue("missing transaction data")
			return v
		}

		if !callback.IsDefined() {
			v, _ := otto.ToValue("missing callback")
			return v
		}

		// reason for this transaction
		reason, err := toSign.Get("reason")
		if err != nil {
			callback.Call(callback, err)
		}
		if !reason.IsString() {
			callback.Call(callback, "reason has to be a string")
			return otto.Value{}
		}

		// "to" of transaction
		to, err := toSign.Get("to")
		if err != nil {
			callback.Call(callback, err)
			return otto.Value{}
		}
		if !to.IsString() {
			callback.Call(callback, "to has to be a valid ethereum address")
			return otto.Value{}
		}
		if !ethCommon.IsHexAddress(to.String()) {
			callback.Call(callback, "to has to be a valid ethereum address")
			return otto.Value{}
		}

		// "value" of transaction
		value, err := toSign.Get("value")
		if err != nil {
			callback.Call(callback, err)
			return otto.Value{}
		}
		if !value.IsString() {
			callback.Call(callback, "value must be a string")
			return otto.Value{}
		}

		// "data" of transaction
		data, err := toSign.Get("data")
		if err != nil {
			callback.Call(callback, err)
			return otto.Value{}
		}
		if !data.IsString() {
			callback.Call(callback, "data must be a string")
			return otto.Value{}
		}

		respChan, err := api.Send(&SignTransactionRpcCall{
			Reason:  reason.String(),
			TxTo:    to.String(),
			TxValue: value.String(),
			TxData:  value.String(),
		})

		if err != nil {
			callback.Call(callback, "failed ot send transaction to device")
			return otto.Value{}
		}

		resp := <-respChan

		if resp.Error != nil {
			logger.Error("sign transaction request failed with error: ", resp.Error)
			resp.Close(nil)
			callback.Call(callback, "received error for sign request. please connect a debugger to see the problem")
			return otto.Value{}
		}

		r := struct {
			SignedTransaction string
		}{}

		if err := json.Unmarshal([]byte(resp.Payload), &r); err != nil {
			logger.Error("Failed to unmarshal response for sign transaction request: ", err)
			resp.Close(err)
			callback.Call(callback, "failed to unmarshal sign transaction request")
			return otto.Value{}
		}

		_, err = callback.Call(callback, nil, r.SignedTransaction)
		resp.Close(err)

		return otto.Value{}

	}

}

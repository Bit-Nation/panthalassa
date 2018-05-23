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

}

func (c *SignTransactionRpcCall) Data() (string, error) {

	raw, err := json.Marshal(c)
	return string(raw), err

}

func (c *SignTransactionRpcCall) Valid() error {

}

type SignTransaction struct {
	Api *deviceApi.Api
}

// a sign call requires
// @todo there is need for validating the parameter length. E.g the data should not be longer than x char's etc
func (s *SignTransaction) SignTransaction(call otto.FunctionCall) {

	toSign := call.Argument(0).Object()
	callback := call.Argument(1)

	// reason for this transaction
	reason, err := toSign.Get("reason")
	if err != nil {
		callback.Call(callback, err)
	}
	if !reason.IsString() {
		callback.Call(callback, "reason has to be a string")
		return
	}

	// "to" of transaction
	to, err := toSign.Get("to")
	if err != nil {
		callback.Call(callback, err)
		return
	}
	if !to.IsString() {
		callback.Call(callback, "to has to be a valid ethereum address")
		return
	}
	if !ethCommon.IsHexAddress(to.String()) {
		callback.Call(callback, "to has to be a valid ethereum address")
		return
	}

	// "value" of transaction
	value, err := toSign.Get("value")
	if err != nil {
		callback.Call(callback, err)
		return
	}
	if !value.IsString() {
		callback.Call(callback, "value must be a string")
		return
	}

	// "data" of transaction
	data, err := toSign.Get("data")
	if err != nil {
		callback.Call(callback, err)
		return
	}
	if !data.IsString() {
		callback.Call(callback, "data must be a string")
		return
	}

	respChan, err := s.Api.Send(&SignTransactionRpcCall{
		Reason:  reason.String(),
		TxTo:    to.String(),
		TxValue: value.String(),
		TxData:  value.String(),
	})

	if err != nil {
		callback.Call(callback, "failed ot send transaction to device")
		return
	}

	resp := <-respChan

}

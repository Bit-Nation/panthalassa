package ethWebSocket

import (
	"encoding/json"
	"testing"
	"time"

	ethws "github.com/Bit-Nation/panthalassa/ethws"
	log "github.com/op/go-logging"
	otto "github.com/robertkrimen/otto"
	require "github.com/stretchr/testify/require"
)

func TestWS(t *testing.T) {

	c := make(chan bool)

	logger := log.MustGetLogger("")

	ethWS := ethws.New(ethws.Config{
		Retry: time.Second,
		WSUrl: "wss://mainnet.infura.io/_ws",
	})

	ethWsModule := New(logger, ethWS)

	vm := otto.New()

	require.Nil(t, ethWsModule.Register(vm))
	_, err := vm.Call(`ethereumRequest`, vm, `{"jsonrpc":"2.0","method":"eth_protocolVersion","params":[]}`, func(call otto.FunctionCall) otto.Value {

		if !call.Argument(0).IsUndefined() {
			panic("first argument is the error which should be undefined")
		}

		response := call.Argument(1).String()

		var resp ethws.Response
		if err := json.Unmarshal([]byte(response), &resp); err != nil {
			panic(err)
		}

		if resp.Result != "0x3f" {
			panic("got unexpected result.")
		}

		c <- true

		return otto.Value{}

	})
	require.Nil(t, err)

	select {
	case <-c:
		return
	case <-time.After(time.Second * 5):
		require.FailNow(t, "timed out")
	}

}

func TestWSFunctionSignatureValidation(t *testing.T) {

	logger := log.MustGetLogger("")

	ethWS := ethws.New(ethws.Config{
		Retry: time.Second,
		WSUrl: "wss://mainnet.infura.io/_ws",
	})

	ethWsModule := New(logger, ethWS)

	vm := otto.New()

	ethWsModule.Register(vm)

	vmErr, err := vm.Call(`ethereumRequest`, vm, nil)
	require.Nil(t, err)
	require.Equal(t, "ValidationError: expected parameter 0 to be of type string", vmErr.String())

	vmErr, err = vm.Call(`ethereumRequest`, vm, "", 9)
	require.Nil(t, err)
	require.Equal(t, "ValidationError: expected parameter 1 to be of type function", vmErr.String())

}

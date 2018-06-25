package ethWebSocket

import (
	"encoding/json"
	"io/ioutil"
	"testing"
	"time"

	ethws "github.com/Bit-Nation/panthalassa/ethws"
	log "github.com/op/go-logging"
	otto "github.com/robertkrimen/otto"
	require "github.com/stretchr/testify/require"
)

func TestWS(t *testing.T) {

	c := make(chan bool)

	ethWS := ethws.New(ethws.Config{
		Retry: time.Second,
		WSUrl: "wss://mainnet.infura.io/_ws",
	})

	m := EthWS{
		ethWS: ethWS,
	}

	vm := otto.New()

	require.Nil(t, m.Register(vm))
	_, err := vm.Call(`ethereumRequest`, vm, `{"jsonrpc":"2.0","method":"eth_protocolVersion","params":[]}`, func(response string) {

		var resp ethws.Response
		if err := json.Unmarshal([]byte(response), &resp); err != nil {
			panic(err)
		}

		if resp.Result != "0x3f" {
			panic("got unexpected result.")
		}

		c <- true

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
	backend := log.NewLogBackend(ioutil.Discard, "", 0)
	logger.SetBackend(log.AddModuleLevel(backend))

	m := EthWS{
		logger: logger,
	}

	vm := otto.New()

	m.Register(vm)

	vmErr, err := vm.Call(`ethereumRequest`, vm, nil)
	require.Nil(t, err)
	require.Equal(t, "ValidationError: expected parameter 0 to be of type string", vmErr.String())

	vmErr, err = vm.Call(`ethereumRequest`, vm, "", 9)
	require.Nil(t, err)
	require.Equal(t, "ValidationError: expected parameter 1 to be of type function", vmErr.String())

}

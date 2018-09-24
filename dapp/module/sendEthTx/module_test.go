package sendEthTx

import (
	"testing"
	"time"

	log "github.com/op/go-logging"
	require "github.com/stretchr/testify/require"
	duktape "gopkg.in/olebedev/go-duktape.v3"
)

type testApi struct {
	send func(value, to, data string) (string, error)
}

func (a *testApi) SendEthereumTransaction(value, to, data string) (string, error) {
	return a.send(value, to, data)
}

func TestSuccessCall(t *testing.T) {

	logger := log.MustGetLogger("")

	vm := duktape.New()

	m := New(&testApi{
		send: func(value, to, data string) (string, error) {
			if value != "10000" {
				panic("expected value to be 10000")
			}
			if to != "0xee60a19d0850b51b8598ca2ceb9acae3f452943d" {
				panic("expected address to be 0xee60a19d0850b51b8598ca2ceb9acae3f452943d")
			}
			if data != "0xf3..." {
				panic("expected data to be 0xf3...")
			}
			return `{}`, nil
		},
	}, logger)

	require.Nil(t, m.Register(vm))

	txReq := `({
		"value": "10000",
		"to": "0xee60a19d0850b51b8598ca2ceb9acae3f452943d",
		"data": "0xf3..."
	})`

	wait := make(chan bool)

	_, err := vm.PushGlobalGoFunction("callbackSendETHTransaction", func(context *duktape.Context) int {

		if !context.IsUndefined(0) {
			panic("expected error to be undefined")
		}

		transaction := (context.SafeToString(-1))
		if transaction != "{}" {
			// it's ok to assert {}
			// this is just a mock
			panic("expected transaction to be {}")
		}

		wait <- true
		return 0
	})
	require.Nil(t, err)

	go vm.PevalString(`sendETHTransaction(` + txReq + `,callbackSendETHTransaction)`)
	require.Nil(t, err)

	select {
	case <-wait:
	case <-time.After(time.Second * 5):
		require.FailNow(t, "timed out")
	}

}

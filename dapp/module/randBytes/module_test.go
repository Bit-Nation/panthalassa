package randBytes

import (
	"bytes"
	"testing"

	log "github.com/op/go-logging"
	otto "github.com/robertkrimen/otto"
	require "github.com/stretchr/testify/require"
)

func TestCreateRandomBytes(t *testing.T) {

	randSource = bytes.NewReader([]byte{1, 4, 6})

	mod := New(log.MustGetLogger(""))

	vm := otto.New()

	require.Nil(t, mod.Register(vm))

	_, err := vm.Call("randomBytes", vm, 3, func(call otto.FunctionCall) otto.Value {

		if !call.Argument(0).IsUndefined() {
			panic("expected error to be undefined")
		}

		generatedBytes, err := call.Argument(1).ToString()
		if err != nil {
			panic(err)
		}

		if generatedBytes != "1, 4, 6" {
			panic(err)
		}

		return otto.Value{}

	})
	require.Nil(t, err)

}

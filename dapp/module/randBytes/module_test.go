package randBytes

import (
	"bytes"
	"testing"

	log "github.com/op/go-logging"
	require "github.com/stretchr/testify/require"
	duktape "gopkg.in/olebedev/go-duktape.v3"
)

func TestCreateRandomBytes(t *testing.T) {

	randSource = bytes.NewReader([]byte{1, 4, 6})

	mod := New(log.MustGetLogger(""))

	vm := duktape.New()

	require.Nil(t, mod.Register(vm))

	vm.PushGlobalGoFunction("callbackRandomBytes", func(context *duktape.Context) int {
		if !context.IsUndefined(0) {
			panic("expected error to be undefined")
		}

		generatedBytes := (context.ToString(-1))
		if generatedBytes != "1, 4, 6" {
			panic("it's not the same bytes")
		}

		return 0
	})
	err := vm.PevalString(`randomBytes(3,callbackRandomBytes)`)
	require.Nil(t, err)

}

package uuidv4

import (
	"errors"
	"testing"
	"time"

	log "github.com/op/go-logging"
	otto "github.com/robertkrimen/otto"
	uuid "github.com/satori/go.uuid"
	require "github.com/stretchr/testify/require"
)

func TestUUIDV4ModelSuccess(t *testing.T) {

	m := UUIDV4{}

	// mock the new uuid function
	newUuid = func() (uuid.UUID, error) {
		return uuid.FromString("9b781c39-2bd3-41c6-a246-150a9f4421a3")
	}

	vm := otto.New()
	vm.Set("test", func(call otto.FunctionCall) otto.Value {
		require.True(t, call.Argument(0).IsUndefined())
		require.Equal(t, "9b781c39-2bd3-41c6-a246-150a9f4421a3", call.Argument(1).String())
		return otto.Value{}
	})

	require.Nil(t, m.Register(vm))

	vm.Run(`uuidV4(test)`)

}

func TestUUIDV4ModelError(t *testing.T) {

	m := New(log.MustGetLogger(""))

	// mock the new uuid function
	newUuid = func() (uuid.UUID, error) {
		return uuid.UUID{}, errors.New("I am a nice error message")
	}

	vm := otto.New()

	require.Nil(t, m.Register(vm))

	wait := make(chan bool, 1)

	vm.Call("uuidV4", vm, func(call otto.FunctionCall) otto.Value {

		if call.Argument(0).String() != "I am a nice error message" {
			panic("expected error message: I am a nice error message")
		}

		if !call.Argument(1).IsUndefined() {
			panic("id should be undefined")
		}

		wait <- true

		return otto.Value{}

	})

	select {
	case <-wait:
	case <-time.After(time.Second):
		require.FailNow(t, "timed out")
	}

}

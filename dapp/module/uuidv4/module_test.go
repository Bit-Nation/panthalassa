package uuidv4

import (
	"errors"
	"testing"

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
		require.Equal(t, "9b781c39-2bd3-41c6-a246-150a9f4421a3", call.Argument(0).String())
		require.True(t, call.Argument(1).IsUndefined())
		return otto.Value{}
	})

	require.Nil(t, m.Register(vm))

	vm.Run(`uuidV4(test)`)

}

func TestUUIDV4ModelError(t *testing.T) {

	m := UUIDV4{}
	// mock the new uuid function
	newUuid = func() (uuid.UUID, error) {
		return uuid.UUID{}, errors.New("I am a nice error message")
	}

	vm := otto.New()
	vm.Set("test", func(call otto.FunctionCall) otto.Value {
		require.True(t, call.Argument(0).IsUndefined())
		require.Equal(t, "I am a nice error message", call.Argument(1).String())
		return otto.Value{}
	})

	require.Nil(t, m.Register(vm))

	_, err := vm.Run(`uuidV4(test)`)
	require.Nil(t, err)

}

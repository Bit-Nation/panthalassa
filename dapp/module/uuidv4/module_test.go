package uuidv4

import (
	"errors"
	"testing"
	"time"

	log "github.com/op/go-logging"
	uuid "github.com/satori/go.uuid"
	require "github.com/stretchr/testify/require"
	duktape "gopkg.in/olebedev/go-duktape.v3"
)

func TestUUIDV4ModelSuccess(t *testing.T) {

	m := UUIDV4{}

	// mock the new uuid function
	newUuid = func() (uuid.UUID, error) {
		return uuid.FromString("9b781c39-2bd3-41c6-a246-150a9f4421a3")
	}

	vm := duktape.New()

	require.Nil(t, m.Register(vm))

	_, err := vm.PushGlobalGoFunction("test", func(context *duktape.Context) int {
		require.True(t, context.IsUndefined(0))
		require.Equal(t, "9b781c39-2bd3-41c6-a246-150a9f4421a3", context.SafeToString(1))
		return 0
	})
	require.Nil(t, err)
	require.Nil(t, m.Register(vm))

	vm.PevalString(`uuidV4(test)`)

}

func TestUUIDV4ModelError(t *testing.T) {

	m := New(log.MustGetLogger(""))

	// mock the new uuid function
	newUuid = func() (uuid.UUID, error) {
		return uuid.UUID{}, errors.New("I am a nice error message")
	}

	vm := duktape.New()

	require.Nil(t, m.Register(vm))

	wait := make(chan bool, 1)

	_, err := vm.PushGlobalGoFunction("callbackUuidV4", func(context *duktape.Context) int {
		if context.SafeToString(0) != "I am a nice error message" {
			panic("expected error message: I am a nice error message")
		}

		if !context.IsUndefined(1) {
			panic("id should be undefined")
		}

		wait <- true

		return 0
	})
	require.Nil(t, err)

	vm.PevalString("uuidV4(callbackUuidV4)")
	select {
	case <-wait:
	case <-time.After(time.Second):
		require.FailNow(t, "timed out")
	}

}

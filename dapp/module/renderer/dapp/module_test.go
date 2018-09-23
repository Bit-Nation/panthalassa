package dapp

import (
	"testing"

	log "github.com/op/go-logging"
	require "github.com/stretchr/testify/require"
	duktape "gopkg.in/olebedev/go-duktape.v3"
)

func TestModule_OpenDAppError(t *testing.T) {

	l := log.MustGetLogger("")

	vm := duktape.New()

	m := New(l)
	require.Nil(t, m.Register(vm))

	_, err := vm.PushGlobalGoFunction("callbackTestModuleOpenDAppError", func(context *duktape.Context) int {

		if !context.IsObject(0) {
			panic("callbackTestFuncCallSuccess : 0 is not an object")
		}
		if !context.GetPropString(0, "key") {
			panic("callbackTestFuncCallSuccess : key missing")
		}

		v := context.SafeToString(-1)

		if v != "value" {
			panic("Expected value of key to be: value")
		}

		context.Pop()
		context.PushString("I am an error")
		context.Call(1)

		return 0
	})
	require.Nil(t, err)
	err = vm.PevalString(`setOpenHandler(callbackTestModuleOpenDAppError)`)
	require.Nil(t, err)
	vm.PevalString(`callbackTestModuleOpenDAppError`)
	err = m.OpenDApp(`{key: "value"}`)
	require.EqualError(t, err, "I am an error")

}

func TestModule_OpenDAppSuccess(t *testing.T) {

	l := log.MustGetLogger("")

	vm := duktape.New()

	m := New(l)
	require.Nil(t, m.Register(vm))

	_, err := vm.PushGlobalGoFunction("callbackTestModuleOpenDAppSuccess", func(context *duktape.Context) int {

		if !context.IsObject(0) {
			panic("callbackTestFuncCallSuccess : 0 is not an object")
		}
		if !context.GetPropString(0, "key") {
			panic("callbackTestFuncCallSuccess : key missing")
		}

		v := context.SafeToString(-1)

		if v != "value" {
			panic("Expected value of key to be: value")
		}

		context.Pop()
		context.PushUndefined()
		context.Call(1)

		return 0
	})
	require.Nil(t, err)
	err = vm.PevalString(`setOpenHandler(callbackTestModuleOpenDAppSuccess)`)
	require.Nil(t, err)
	err = vm.PevalString(`callbackTestModuleOpenDAppSuccess`)
	require.Nil(t, err)
	err = m.OpenDApp(`{key: "value"}`)
	require.Nil(t, err)

}

func TestModule_Close(t *testing.T) {

	// setup
	vm := duktape.New()
	m := New(nil)
	require.Nil(t, m.Register(vm))

	// set fake open handler
	_, err := vm.PushGlobalGoFunction("callbackTestModuleClose", func(context *duktape.Context) int {
		m.Close()
		return 0
	})
	require.Nil(t, err)
	err = vm.PevalString(`setOpenHandler(callbackTestModuleClose)`)
	require.Nil(t, err)
	vm.PevalString(`callbackTestModuleClose`)
	require.EqualError(t, m.OpenDApp("{}"), "closed the application")

}

package dapp

import (
	"testing"

	log "github.com/op/go-logging"
	require "github.com/stretchr/testify/require"
	duktape "gopkg.in/olebedev/go-duktape.v3"
)

func TestModule_RenderMessageError(t *testing.T) {

	l := log.MustGetLogger("")

	vm := duktape.New()

	m := New(l)
	require.Nil(t, m.Register(vm))
	_, err := vm.PushGlobalGoFunction("callbackTestModuleRenderMessageError", func(context *duktape.Context) int {
		context.PushString("I am an error")
		context.Call(1)
		return 0
	})
	require.Nil(t, err)
	err = vm.PevalString("setMessageRenderer(callbackTestModuleRenderMessageError)")
	require.Nil(t, err)

	vm.PevalString(`callbackTestModuleRenderMessageError`)
	_, err = m.RenderMessage(`{}`)
	require.EqualError(t, err, "I am an error")

}

func TestModule_RenderMessageSuccess(t *testing.T) {

	l := log.MustGetLogger("")

	vm := duktape.New()

	m := New(l)
	require.Nil(t, m.Register(vm))
	_, err := vm.PushGlobalGoFunction("callbackTestModuleRenderMessageSuccess", func(context *duktape.Context) int {

		if !context.IsObject(0) {
			panic("callbackTestFuncCallSuccess : 0 is not an object")
		}
		if !context.GetPropString(0, "message") {
			panic("callbackTestFuncCallSuccess : message missing")
		}

		if !context.IsObject(-1) {
			panic("Expected message to be an object")
		}

		context.Pop()
		context.PushUndefined()
		context.PushString(`{}`)
		context.Call(2)

		return 0

	})
	require.Nil(t, err)
	err = vm.PevalString("setMessageRenderer(callbackTestModuleRenderMessageSuccess)")
	require.Nil(t, err)
	vm.PevalString(`callbackTestModuleRenderMessageSuccess`)
	layout, err := m.RenderMessage(`{message: {}, context: {}}`)
	require.Nil(t, err)
	require.Equal(t, "{}", layout)

}

func TestModule_Close(t *testing.T) {

	vm := duktape.New()
	m := New(log.MustGetLogger(""))
	require.Nil(t, m.Register(vm))

	_, err := vm.PushGlobalGoFunction("callbackTestModuleClose", func(context *duktape.Context) int {
		m.Close()
		return 0
	})
	require.Nil(t, err)
	err = vm.PevalString(`setMessageRenderer(callbackTestModuleClose)`)
	require.Nil(t, err)
	vm.PevalString(`callbackTestModuleClose`)
	_, err = m.RenderMessage("{}")
	require.EqualError(t, err, "closed the application")

}

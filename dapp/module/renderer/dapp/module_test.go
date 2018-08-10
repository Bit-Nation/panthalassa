package dapp

import (
	"testing"

	log "github.com/op/go-logging"
	otto "github.com/robertkrimen/otto"
	require "github.com/stretchr/testify/require"
)

func TestModule_OpenDAppError(t *testing.T) {

	l := log.MustGetLogger("")

	vm := otto.New()

	m := New(l)
	require.Nil(t, m.Register(vm))

	vm.Call("setOpenHandler", vm, func(context otto.Value, cb otto.Value) otto.Value {

		v, err := context.Object().Get("key")
		if err != nil {
			panic(err)
		}

		if v.String() != "value" {
			panic("Expected value of key to be: value")
		}

		cb.Call(cb, "I am an error")

		return otto.Value{}

	})

	err := m.OpenDApp(`{key: "value"}`)
	require.EqualError(t, err, "I am an error")

}

func TestModule_OpenDAppSuccess(t *testing.T) {

	l := log.MustGetLogger("")

	vm := otto.New()

	m := New(l)
	require.Nil(t, m.Register(vm))

	vm.Call("setOpenHandler", vm, func(context otto.Value, cb otto.Value) otto.Value {

		v, err := context.Object().Get("key")
		if err != nil {
			panic(err)
		}

		if v.String() != "value" {
			panic("Expected value of key to be: value")
		}

		cb.Call(cb, nil)

		return otto.Value{}

	})

	err := m.OpenDApp(`{key: "value"}`)
	require.Nil(t, err)

}

func TestModule_Close(t *testing.T) {

	// setup
	vm := otto.New()
	m := New(nil)
	require.Nil(t, m.Register(vm))

	// set fake open handler
	_, err := vm.Call("setOpenHandler", vm, func(call otto.FunctionCall) otto.Value {
		m.Close()
		return otto.Value{}
	})
	require.Nil(t, err)

	require.EqualError(t, m.OpenDApp("{}"), "closed the application")

}

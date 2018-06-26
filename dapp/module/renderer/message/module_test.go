package dapp

import (
	"testing"

	log "github.com/op/go-logging"
	otto "github.com/robertkrimen/otto"
	require "github.com/stretchr/testify/require"
)

func TestModule_RenderMessageError(t *testing.T) {

	l := log.MustGetLogger("")

	vm := otto.New()

	m := New(l)
	require.Nil(t, m.Register(vm))

	vm.Call("setMessageHandler", vm, func(msg otto.Value, context otto.Value, cb otto.Value) otto.Value {

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

	_, err := m.RenderMessage(`{sender: "0x0"}`, `{key: "value"}`)
	require.EqualError(t, err, "I am an error")

}

func TestModule_RenderMessageSuccess(t *testing.T) {

	l := log.MustGetLogger("")

	vm := otto.New()

	m := New(l)
	require.Nil(t, m.Register(vm))

	vm.Call("setMessageHandler", vm, func(msg otto.Value, context otto.Value, cb otto.Value) otto.Value {

		v, err := context.Object().Get("key")
		if err != nil {
			panic(err)
		}

		if v.String() != "value" {
			panic("Expected value of key to be: value")
		}

		cb.Call(cb, nil, "{}")

		return otto.Value{}

	})

	layout, err := m.RenderMessage(`{from: "0x0"}`, `{key: "value"}`)
	require.Nil(t, err)
	require.Equal(t, "{}", layout)

}

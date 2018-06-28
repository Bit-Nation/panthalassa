package callbacks

import (
	"testing"

	"fmt"
	otto "github.com/robertkrimen/otto"
	require "github.com/stretchr/testify/require"
)

func TestFuncRegistration(t *testing.T) {

	m := New(nil)

	vm := otto.New()

	require.Nil(t, m.Register(vm))

	_, err := vm.Call(`registerFunction`, vm, func(call otto.FunctionCall) otto.Value {
		return otto.Value{}
	})
	require.Nil(t, err)

	fn, exist := m.functions[1]
	require.True(t, exist)
	require.True(t, fn.IsFunction())

}

func TestFuncCallSuccess(t *testing.T) {

	m := New(nil)

	vm := otto.New()

	require.Nil(t, m.Register(vm))

	_, err := vm.Call(`registerFunction`, vm, func(call otto.FunctionCall) otto.Value {

		valueFromObj, err := call.Argument(0).Object().Get("key")
		if err != nil {
			panic(err)
		}

		if valueFromObj.String() != "value" {
			panic("expected value of key to be: value")
		}

		cb := call.Argument(1)

		if !cb.IsFunction() {
			panic("expected second argument to be a callback")
		}

		_, err = cb.Call(cb)
		if err != nil {
			m.logger.Error(err.Error())
		}

		return otto.Value{}
	})

	require.Nil(t, err)
	fn, exist := m.functions[1]
	require.True(t, exist)
	require.True(t, fn.IsFunction())

	require.Nil(t, m.CallFunction(1, `{key: "value"}`))

}

func TestFuncCallError(t *testing.T) {

	m := New(nil)

	vm := otto.New()

	require.Nil(t, m.Register(vm))

	_, err := vm.Call(`registerFunction`, vm, func(call otto.FunctionCall) otto.Value {

		valueFromObj, err := call.Argument(0).Object().Get("key")
		if err != nil {
			panic(err)
		}

		if valueFromObj.String() != "value" {
			panic("expected value of key to be: value")
		}

		cb := call.Argument(1)

		if !cb.IsFunction() {
			panic("expected second argument to be a callback")
		}

		_, err = cb.Call(cb, "I am an error")

		if err != nil {
			m.logger.Error(err.Error())
		}

		return otto.Value{}
	})

	require.Nil(t, err)
	fn, exist := m.functions[1]
	require.True(t, exist)
	require.True(t, fn.IsFunction())

	require.Equal(t, "I am an error", m.CallFunction(1, `{key: "value"}`).Error())

}

func TestFuncCallBackTwice(t *testing.T) {

	m := New(nil)

	vm := otto.New()

	require.Nil(t, m.Register(vm))

	_, err := vm.Call(`registerFunction`, vm, func(call otto.FunctionCall) otto.Value {

		valueFromObj, err := call.Argument(0).Object().Get("key")
		if err != nil {
			panic(err)
		}

		if valueFromObj.String() != "value" {
			panic("expected value of key to be: value")
		}

		cb := call.Argument(1)

		if !cb.IsFunction() {
			panic("expected second argument to be a callback")
		}

		_, err = cb.Call(cb, "I am an error")
		_, err = cb.Call(cb, "I am an error")

		fmt.Println(err)

		return otto.Value{}
	})

	require.Nil(t, err)
	fn, exist := m.functions[1]
	require.True(t, exist)
	require.True(t, fn.IsFunction())

	require.Equal(t, "I am an error", m.CallFunction(1, `{key: "value"}`).Error())

}

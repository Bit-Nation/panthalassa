package callbacks

import (
	"strconv"
	"testing"
	"time"

	log "github.com/op/go-logging"
	require "github.com/stretchr/testify/require"
	duktape "gopkg.in/olebedev/go-duktape.v3"
)

func TestFuncRegistration(t *testing.T) {

	m := New(nil)

	vm := duktape.New()

	require.Nil(t, m.Register(vm))

	require.Equal(t, uint(0), m.reqLim.Count())
	funcId, err := vm.PushGlobalGoFunction("callbackTestFuncRegistration", func(context *duktape.Context) int {
		return 0
	})
	require.Nil(t, err)
	err = vm.PevalString(`registerFunction(callbackTestFuncRegistration)`)
	require.Nil(t, err)
	require.Equal(t, uint(1), m.reqLim.Count())

	id := funcId
	require.Equal(t, int(1), id)

	// fetch function and assert
	respChan := make(chan *duktape.Context)
	m.fetchFunctionChan <- fetchFunction{
		id:       1,
		respChan: respChan,
	}
	require.NotNil(t, <-respChan)

}

func TestFuncUnRegisterSuccess(t *testing.T) {

	m := New(nil)
	vm := duktape.New()
	require.Nil(t, m.Register(vm))

	// register function
	funcId, err := vm.PushGlobalGoFunction("callbackTestFuncUnRegisterSuccess", func(context *duktape.Context) int {
		return 0
	})
	require.Nil(t, err)
	err = vm.PevalString(`registerFunction(callbackTestFuncUnRegisterSuccess)`)
	require.Nil(t, err)
	// make sure the id is the one we expect it to be
	id := funcId
	require.Equal(t, int(1), id)

	// make sure function exist
	respChan := make(chan *duktape.Context)
	m.fetchFunctionChan <- fetchFunction{
		id:       1,
		respChan: respChan,
	}
	require.NotNil(t, <-respChan)
	require.Equal(t, uint(1), m.reqLim.Count())

	// un-register function that has never been registered
	err = vm.PevalString(`unRegisterFunction(23565)`)
	require.Nil(t, err)
	// request limitation must still be one since function that have not been registered
	// can't decrease the counter
	require.Equal(t, uint(1), m.reqLim.Count())

	// successfully un-register function
	err = vm.PevalString(`unRegisterFunction(` + strconv.Itoa(id) + `)`)
	require.Nil(t, err)
	// wait for the go routine to sync up
	time.Sleep(time.Millisecond * 100)
	// should be 0 since the the function id we passed in exists
	require.Equal(t, uint(0), m.reqLim.Count())

}

func TestFuncCallSuccess(t *testing.T) {

	m := New(nil)
	vm := duktape.New()
	require.Nil(t, m.Register(vm))

	_, err := vm.PushGlobalGoFunction("callbackTestFuncCallSuccess", func(context *duktape.Context) int {

		if !context.IsObject(0) {
			panic("callbackTestFuncCallSuccess : 0 is not an object")
		}
		if !context.GetPropString(0, "key") {
			panic("callbackTestFuncCallSuccess : key missing")
		}

		valueFromObj := context.ToString(-1)
		context.Pop()
		if valueFromObj != "value" {
			panic("expected value of key to be: value")
		}

		if !context.IsFunction(1) {
			panic("expected second argument to be a callback")
		}

		context.PushUndefined()
		context.Call(1)

		return 0
	})
	require.Nil(t, err)

	err = vm.PevalString(`registerFunction(callbackTestFuncCallSuccess)`)
	require.Nil(t, err)

	// make sure function exist
	respChan := make(chan *duktape.Context)
	m.fetchFunctionChan <- fetchFunction{
		id:       1,
		respChan: respChan,
	}
	require.NotNil(t, <-respChan)

	require.Nil(t, m.CallFunction(1, `{key: "value"}`))

}

func TestFuncCallError(t *testing.T) {

	m := New(nil)

	vm := duktape.New()

	require.Nil(t, m.Register(vm))

	_, err := vm.PushGlobalGoFunction("callbackTestFuncCallError", func(context *duktape.Context) int {

		if !context.IsObject(0) {
			panic("callbackTestFuncCallSuccess : 0 is not an object")
		}
		if !context.GetPropString(0, "key") {
			panic("callbackTestFuncCallSuccess : key missing")
		}

		valueFromObj := context.ToString(-1)
		if valueFromObj != "value" {
			panic("expected value of key to be: value")
		}

		if !context.IsFunction(1) {
			panic("expected second argument to be a callback")
		}

		context.Pop()
		context.PushString("I am an error")
		context.Call(1)

		return 0
	})
	require.Nil(t, err)
	err = vm.PevalString(`registerFunction(callbackTestFuncCallError)`)
	require.Nil(t, err)

	// make sure function exist
	respChan := make(chan *duktape.Context)
	m.fetchFunctionChan <- fetchFunction{
		id:       1,
		respChan: respChan,
	}
	require.NotNil(t, <-respChan)

	require.Equal(t, "I am an error", m.CallFunction(1, `{key: "value"}`).Error())

}

func TestFuncCallBackTwice(t *testing.T) {

	m := New(log.MustGetLogger(""))

	vm := duktape.New()

	require.Nil(t, m.Register(vm))

	_, err := vm.PushGlobalGoFunction("callbackTestFuncCallBackTwice", func(context *duktape.Context) int {

		if !context.IsObject(0) {
			panic("callbackTestFuncCallSuccess : 0 is not an object")
		}
		if !context.GetPropString(0, "key") {
			panic("callbackTestFuncCallSuccess : key missing")
		}

		valueFromObj := context.ToString(-1)

		if valueFromObj != "value" {
			panic("expected value of key to be: value")
		}

		if !context.IsFunction(1) {
			panic("expected second argument to be a callback")
		}
		context.Pop()
		context.DupTop()
		context.PushUndefined()
		context.Call(1)
		// @TODO find a way to test "expected value to be undefined"
		//if !val.IsUndefined() {
		//	panic("expected value to be undefined")
		//}
		context.Pop()
		context.PushUndefined()
		context.Call(1)
		// @TODO find a way to test "Callback: Already called callback"
		//if val.String() != "Callback: Already called callback" {
		//	panic("Expected an error that tells me that I alrady called the callback")
		//}

		return 0

	})
	require.Nil(t, err)
	err = vm.PevalString(`registerFunction(callbackTestFuncCallBackTwice)`)
	require.Nil(t, err)

	// make sure function exist
	respChan := make(chan *duktape.Context)
	m.fetchFunctionChan <- fetchFunction{
		id:       1,
		respChan: respChan,
	}
	require.NotNil(t, <-respChan)

	m.CallFunction(1, `{key: "value"}`)

}

func TestModule_CallFunctionThatIsNotRegistered(t *testing.T) {

	m := New(log.MustGetLogger(""))
	vm := duktape.New()
	require.Nil(t, m.Register(vm))

	err := m.CallFunction(1, "{}")
	require.EqualError(t, err, "function with id: 1 does not exist")

}

func TestModule_Close(t *testing.T) {

	m := New(log.MustGetLogger(""))
	vm := duktape.New()
	require.Nil(t, m.Register(vm))

	_, err := vm.PushGlobalGoFunction("callbackTestModuleClose", func(context *duktape.Context) int {
		m.Close()
		return 0
	})
	require.Nil(t, err)

	// register function
	err = vm.PevalString(`registerFunction(callbackTestModuleClose)`)
	require.Nil(t, err)

	require.EqualError(t, m.CallFunction(1, "{}"), "closed application")

}

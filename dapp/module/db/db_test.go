package db

import (
	"encoding/json"
	"testing"
	"time"

	otto "github.com/robertkrimen/otto"
	require "github.com/stretchr/testify/require"
)

type inMemoryDB struct {
	storage map[string][]byte
}

func (db *inMemoryDB) Put(key, value []byte) error {
	db.storage[string(key)] = value
	return nil
}

func (db *inMemoryDB) Get(key []byte) ([]byte, error) {
	value := db.storage[string(key)]
	return value, nil
}

func (db *inMemoryDB) Has(key []byte) (bool, error) {
	_, exist := db.storage[string(key)]
	return exist, nil
}

func (db *inMemoryDB) Delete(key []byte) error {
	delete(db.storage, string(key))
	return nil
}

func TestModulePut(t *testing.T) {

	m := Module{
		dAppDB: &inMemoryDB{
			storage: map[string][]byte{},
		},
	}

	vm := otto.New()

	require.Nil(t, m.Register(vm))

	closer := make(chan struct{}, 1)
	vm.Call("db.put", vm, "key", "value", func(call otto.FunctionCall) otto.Value {
		err := call.Argument(0)
		require.False(t, err.IsDefined())

		closer <- struct{}{}

		return otto.Value{}
	})

	select {
	case <-closer:
	case <-time.After(time.Second * 2):
		require.Fail(t, "timed out")
	}

}

func TestModuleHas(t *testing.T) {

	m := Module{
		dAppDB: &inMemoryDB{
			storage: map[string][]byte{
				"key": []byte("value"),
			},
		},
	}

	vm := otto.New()

	require.Nil(t, m.Register(vm))

	closer := make(chan struct{}, 1)
	vm.Call("db.has", vm, "key", func(call otto.FunctionCall) otto.Value {

		err := call.Argument(0)
		require.False(t, err.IsDefined())

		exist, e := call.Argument(1).ToBoolean()
		require.Nil(t, e)
		require.True(t, exist)

		closer <- struct{}{}

		return otto.Value{}
	})

	select {
	case <-closer:
	case <-time.After(time.Second * 2):
		require.Fail(t, "timed out")
	}

}

func TestModuleGet(t *testing.T) {

	marshaledValue, err := json.Marshal("value")
	require.Nil(t, err)

	m := Module{
		dAppDB: &inMemoryDB{
			storage: map[string][]byte{
				"key": marshaledValue,
			},
		},
	}

	vm := otto.New()

	require.Nil(t, m.Register(vm))

	closer := make(chan struct{}, 1)
	vm.Call("db.get", vm, "key", func(call otto.FunctionCall) otto.Value {

		err := call.Argument(0)
		require.False(t, err.IsDefined())

		value := call.Argument(1).String()
		require.Equal(t, "value", value)

		closer <- struct{}{}

		return otto.Value{}
	})

	select {
	case <-closer:
	case <-time.After(time.Second * 2):
		require.Fail(t, "timed out")
	}

}

func TestModuleDelete(t *testing.T) {

	marshaledValue, err := json.Marshal("value")
	require.Nil(t, err)

	m := Module{
		dAppDB: &inMemoryDB{
			storage: map[string][]byte{
				"key": marshaledValue,
			},
		},
	}

	vm := otto.New()

	require.Nil(t, m.Register(vm))

	closer := make(chan struct{}, 1)
	vm.Call("db.delete", vm, "key", func(call otto.FunctionCall) otto.Value {

		err := call.Argument(0)
		require.False(t, err.IsDefined())

		closer <- struct{}{}

		return otto.Value{}
	})

	select {
	case <-closer:
	case <-time.After(time.Second * 2):
		require.Fail(t, "timed out")
	}

}

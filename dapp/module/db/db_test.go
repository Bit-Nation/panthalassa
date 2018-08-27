package db

import (
	"encoding/json"
	"testing"
	"time"

	log "github.com/op/go-logging"
	require "github.com/stretchr/testify/require"
	duktape "gopkg.in/olebedev/go-duktape.v3"
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

	m := New(&inMemoryDB{
		storage: map[string][]byte{},
	}, log.MustGetLogger(""))

	vm := duktape.New()

	require.Nil(t, m.Register(vm))

	closer := make(chan struct{}, 1)
	_, err := vm.PushGlobalGoFunction("callbackDbPut", func(context *duktape.Context) int {
		errBool := !context.IsUndefined(0)
		require.False(t, errBool)

		closer <- struct{}{}

		return 0
	})
	require.Nil(t, err)

	vm.PevalString(`dbPut("key","value",callbackDbPut)`)

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

	vm := duktape.New()

	require.Nil(t, m.Register(vm))

	closer := make(chan struct{}, 1)

	_, err := vm.PushGlobalGoFunction("callbackDbHas", func(context *duktape.Context) int {
		errBool := !context.IsUndefined(0)
		require.False(t, errBool)

		exist := context.IsBoolean(1)
		require.True(t, exist)

		closer <- struct{}{}

		return 0
	})
	require.Nil(t, err)

	vm.PevalString(`dbHas("key",callbackDbHas)`)

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
		logger: log.MustGetLogger(""),
	}

	vm := duktape.New()

	require.Nil(t, m.Register(vm))

	closer := make(chan struct{}, 1)
	_, err = vm.PushGlobalGoFunction("callbackDbGet", func(context *duktape.Context) int {

		errBool := !context.IsUndefined(0)
		require.False(t, errBool)

		value := context.ToString(1)
		require.Equal(t, "value", value)

		closer <- struct{}{}

		return 0
	})
	require.Nil(t, err)

	vm.PevalString(`dbGet("key",callbackDbGet)`)
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

	vm := duktape.New()

	require.Nil(t, m.Register(vm))

	closer := make(chan struct{}, 1)

	_, err = vm.PushGlobalGoFunction("callbackDbDelete", func(context *duktape.Context) int {

		errBool := !context.IsUndefined(0)
		require.False(t, errBool)

		closer <- struct{}{}

		return 0
	})

	vm.PevalString(`dbDelete("key", callbackDbDelete)`)

	select {
	case <-closer:
	case <-time.After(time.Second * 2):
		require.Fail(t, "timed out")
	}

}

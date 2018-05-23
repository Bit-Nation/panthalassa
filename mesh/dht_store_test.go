package mesh

import (
	"encoding/json"
	"fmt"
	"testing"

	deviceApi "github.com/Bit-Nation/panthalassa/api/device"
	ds "github.com/ipfs/go-datastore"
	require "github.com/stretchr/testify/require"
)

type testUpStream struct {
	commChan chan string
}

func (u *testUpStream) Send(data string) {
	u.commChan <- data
}

func requireEqual(expected interface{}, value interface{}) {

	if expected != value {
		panic(fmt.Sprintf("expected: %s got %s", expected, value))
	}

}

func TestRequireEqual(t *testing.T) {

	// should throw since a != b
	require.Panics(t, func() {
		requireEqual("a", "b")
	})

	// should pass since a == b
	require.NotPanics(t, func() {
		requireEqual("a", "a")
	})

}

func TestDHTPutSuccess(t *testing.T) {

	device := make(chan string, 1)

	u := testUpStream{
		commChan: device,
	}

	api := deviceApi.New(&u)

	go func() {

		// wait for the rpc call
		rpcCall := <-device

		// unmarshal rpc call
		var call deviceApi.ApiCall
		if err := json.Unmarshal([]byte(rpcCall), &call); err != nil {
			panic(err)
		}

		requireEqual("DHT:PUT", call.Type)

		// unmarshal submitted values
		var payload DHTPutCall
		if err := json.Unmarshal([]byte(call.Data), &payload); err != nil {
			panic(err)
		}

		requireEqual("/my_key", payload.Key)
		requireEqual("my_value", payload.Value)

		err := api.Receive(call.Id, `{}`)
		if err != nil {
			panic(err)
		}

	}()

	dht := DataStore{api}

	err := dht.Put(ds.NewKey("/my_key"), "my_value")
	require.Nil(t, err)

}

// test PUT but fail to persist and send error back
func TestDHTPutFail(t *testing.T) {

	device := make(chan string, 1)

	u := testUpStream{
		commChan: device,
	}

	api := deviceApi.New(&u)

	go func() {

		// wait for the rpc call
		rpcCall := <-device

		// unmarshal rpc call
		var call deviceApi.ApiCall
		if err := json.Unmarshal([]byte(rpcCall), &call); err != nil {
			panic(err)
		}

		requireEqual("DHT:PUT", call.Type)

		// unmarshal submitted values
		var payload DHTPutCall
		if err := json.Unmarshal([]byte(call.Data), &payload); err != nil {
			panic(err)
		}

		requireEqual("/my_key", payload.Key)
		requireEqual("my_value", payload.Value)

		err := api.Receive(call.Id, `{"error": "not enough dist space or what ever"}`)
		if err != nil {
			panic(err)
		}

	}()

	dht := DataStore{api}

	err := dht.Put(ds.NewKey("/my_key"), "my_value")
	require.EqualError(t, err, "not enough dist space or what ever")

}

func TestDHTGetSuccess(t *testing.T) {

	device := make(chan string, 1)

	u := testUpStream{
		commChan: device,
	}

	api := deviceApi.New(&u)

	go func() {

		// wait for the rpc call
		rpcCall := <-device

		// unmarshal rpc call
		var call deviceApi.ApiCall
		if err := json.Unmarshal([]byte(rpcCall), &call); err != nil {
			panic(err)
		}

		requireEqual("DHT:GET", call.Type)

		// unmarshal submitted values
		var payload DHTGetCall
		if err := json.Unmarshal([]byte(call.Data), &payload); err != nil {
			panic(err)
		}

		requireEqual("/my_key", payload.Key)

		err := api.Receive(call.Id, `{"payload":"{\"value\":\"mapped_value\"}"}`)
		if err != nil {
			panic(err)
		}

	}()

	dht := DataStore{api}

	value, err := dht.Get(ds.NewKey("my_key"))
	require.Nil(t, err)
	require.Equal(t, "mapped_value", value)

}

func TestDHTGetError(t *testing.T) {

	device := make(chan string, 1)

	u := testUpStream{
		commChan: device,
	}

	api := deviceApi.New(&u)

	go func() {

		// wait for the rpc call
		rpcCall := <-device

		// unmarshal rpc call
		var call deviceApi.ApiCall
		if err := json.Unmarshal([]byte(rpcCall), &call); err != nil {
			panic(err)
		}

		requireEqual("DHT:GET", call.Type)

		// unmarshal submitted values
		var payload DHTGetCall
		if err := json.Unmarshal([]byte(call.Data), &payload); err != nil {
			panic(err)
		}

		requireEqual("/my_key", payload.Key)

		err := api.Receive(call.Id, `{"error":"couldn't find mapped value"}`)
		if err != nil {
			panic(err)
		}

	}()

	dht := DataStore{api}

	value, err := dht.Get(ds.NewKey("my_key"))
	require.EqualError(t, err, "couldn't find mapped value")
	require.Nil(t, value)

}

func TestDHTHasSuccess(t *testing.T) {

	device := make(chan string, 1)

	u := testUpStream{
		commChan: device,
	}

	api := deviceApi.New(&u)

	go func() {

		// wait for the rpc call
		rpcCall := <-device

		// unmarshal rpc call
		var call deviceApi.ApiCall
		if err := json.Unmarshal([]byte(rpcCall), &call); err != nil {
			panic(err)
		}

		requireEqual("DHT:HAS", call.Type)

		// unmarshal submitted values
		var payload DHTGetCall
		if err := json.Unmarshal([]byte(call.Data), &payload); err != nil {
			panic(err)
		}

		requireEqual("/my_key", payload.Key)

		err := api.Receive(call.Id, `{"payload":"{\"exist\":true}"}`)
		if err != nil {
			panic(err)
		}

	}()

	dht := DataStore{api}

	exist, err := dht.Has(ds.NewKey("my_key"))
	require.Nil(t, err)
	require.True(t, exist)

}

func TestDHTHasError(t *testing.T) {

	device := make(chan string, 1)

	u := testUpStream{
		commChan: device,
	}

	api := deviceApi.New(&u)

	go func() {

		// wait for the rpc call
		rpcCall := <-device

		// unmarshal rpc call
		var call deviceApi.ApiCall
		if err := json.Unmarshal([]byte(rpcCall), &call); err != nil {
			panic(err)
		}

		requireEqual("DHT:HAS", call.Type)

		// unmarshal submitted values
		var payload DHTGetCall
		if err := json.Unmarshal([]byte(call.Data), &payload); err != nil {
			panic(err)
		}

		requireEqual("/my_key", payload.Key)

		err := api.Receive(call.Id, `{"error":"my error"}`)
		if err != nil {
			panic(err)
		}

	}()

	dht := DataStore{api}

	exist, err := dht.Has(ds.NewKey("my_key"))
	require.EqualError(t, err, "my error")
	require.False(t, exist)

}

func TestDHTDeleteSuccess(t *testing.T) {

	device := make(chan string, 1)

	u := testUpStream{
		commChan: device,
	}

	api := deviceApi.New(&u)

	go func() {

		// wait for the rpc call
		rpcCall := <-device

		// unmarshal rpc call
		var call deviceApi.ApiCall
		if err := json.Unmarshal([]byte(rpcCall), &call); err != nil {
			panic(err)
		}

		requireEqual("DHT:DELETE", call.Type)

		// unmarshal submitted values
		var payload DHTGetCall
		if err := json.Unmarshal([]byte(call.Data), &payload); err != nil {
			panic(err)
		}

		requireEqual("/my_key", payload.Key)

		err := api.Receive(call.Id, `{}`)
		if err != nil {
			panic(err)
		}

	}()

	dht := DataStore{api}

	err := dht.Delete(ds.NewKey("my_key"))
	require.Nil(t, err)
}

func TestDHTDeleteError(t *testing.T) {

	device := make(chan string, 1)

	u := testUpStream{
		commChan: device,
	}

	api := deviceApi.New(&u)

	go func() {

		// wait for the rpc call
		rpcCall := <-device

		// unmarshal rpc call
		var call deviceApi.ApiCall
		if err := json.Unmarshal([]byte(rpcCall), &call); err != nil {
			panic(err)
		}

		requireEqual("DHT:DELETE", call.Type)

		// unmarshal submitted values
		var payload DHTGetCall
		if err := json.Unmarshal([]byte(call.Data), &payload); err != nil {
			panic(err)
		}

		requireEqual("/my_key", payload.Key)

		err := api.Receive(call.Id, `{"error": "my custom error message"}`)
		if err != nil {
			panic(err)
		}

	}()

	dht := DataStore{api}

	err := dht.Delete(ds.NewKey("my_key"))
	require.EqualError(t, err, "my custom error message")
}

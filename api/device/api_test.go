package device_api

import (
	"github.com/ipfs/go-log"
	"github.com/stretchr/testify/require"
	"testing"

	"encoding/json"
)

type upStreamTest struct {
	send func(string)
}

func (u *upStreamTest) Send(data string) {
	u.send(data)
}

type testRPCCall struct {
	callType  string
	data      string
	dataError error
	valid     func(data string) error
}

func (c *testRPCCall) Type() string {
	return c.callType
}
func (c *testRPCCall) Data() (string, error) {
	return c.data, c.dataError
}
func (c *testRPCCall) Valid() error {
	return c.valid(c.data)
}

func Test(t *testing.T) {
	log.SetDebugLogging()

	//The api call we got from the send function
	var receivedApiCall apiCall

	//Create up stream test implementation
	upStream := upStreamTest{
		send: func(data string) {
			var call apiCall
			require.Nil(t, json.Unmarshal([]byte(data), &call))
			receivedApiCall = call
		},
	}

	//Create api
	api := New(&upStream)

	//Send api call with test data
	respChan, err := api.Send(&testRPCCall{
		callType:  "Test",
		data:      `{"key": "value"}`,
		dataError: nil,
		valid: func(data string) error {
			return nil
		},
	})
	require.Nil(t, err)

	//Waiting for the response form the api AND we then close it
	go func() {
		for {
			res := <-respChan
			require.Equal(t, "response", res.Closer)
			res.Closer <- nil
		}
	}()

	//This assertion will be evaluated then "res.Closer <- nil" will be "done"
	//"Receive" will only return after a value was send into the "res.Closer" channel
	require.Nil(t, api.Receive(receivedApiCall.Id, "response"))

}

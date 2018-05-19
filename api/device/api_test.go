package device_api

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	require "github.com/stretchr/testify/require"
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

func TestSuccess(t *testing.T) {

	//The api call we got from the send function
	var receivedApiCall apiCall

	// this is the sample response
	const data = `{"error":"","payload":"my_data"}`

	// for internal use
	c := make(chan struct{}, 1)

	//Create up stream test implementation
	upStream := upStreamTest{
		// this method is implemented by the client
		// and called in an async way.
		send: func(data string) {
			// we set receivedApiCall to the call in order to use it later in the test
			var call apiCall
			require.Nil(t, json.Unmarshal([]byte(data), &call))
			receivedApiCall = call
			c <- struct{}{}
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

	// here we are waiting for the upstream to set "receivedApiCall" to the received call
	<-c

	require.Nil(t, err)
	require.Equal(t, "Test", receivedApiCall.Type)
	require.Equal(t, `{"key": "value"}`, receivedApiCall.Data)

	//Waiting for the response form the api AND we then close it
	go func() {
		for {
			select {
			case res := <-respChan:
				if res.Payload != "my_data" {
					res.Closer <- errors.New(fmt.Sprintf("payload (%s) doesn't match expected payload (%s)", res.Payload, `my_data`))
					return
				}
				res.Closer <- nil
			}
		}
	}()

	//This assertion will be evaluated then "res.Closer <- nil" will be "done"
	//"Receive" will only return after a value was send into the "res.Closer" channel
	require.Nil(t, api.Receive(receivedApiCall.Id, data))

}

func TestApi_SendWithMissingRequest(t *testing.T) {

	//Create up stream test implementation
	upStream := upStreamTest{
		send: func(data string) {
			panic("I am not supposed to be called")
		},
	}

	//Create api
	api := New(&upStream)

	err := api.Receive(uint32(3), "")
	require.EqualError(t, err, "a request channel for id (3) does not exist")

}

package api

import (
	"testing"
	"time"

	"encoding/base64"
	pb "github.com/Bit-Nation/panthalassa/api/pb"
	proto "github.com/golang/protobuf/proto"
	require "github.com/stretchr/testify/require"
)

type testUpStream struct {
	sendFn func(data string)
}

func (u *testUpStream) Send(data string) {

	// code base64 request
	rawReq, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		panic(err)
	}

	u.sendFn(string(rawReq))
}

// make sure that requests are added / removed correct
func TestAPI_addAndCutRequestWorks(t *testing.T) {

	req := pb.Request{}
	req.RequestID = "hi"

	// api
	api := New(&testUpStream{}, keyManagerFactory())

	// make sure request doesn't exist
	_, exist := api.requests["hi"]
	require.False(t, exist)

	api.addRequest(&req)

	// make sure request does exist
	_, exist = api.requests["hi"]
	require.True(t, exist)

	// now cut request our of the stack and make sure it was removed
	api.cutRequest("hi")
	_, exist = api.requests["hi"]
	require.False(t, exist)

}

func TestRequestResponse(t *testing.T) {

	dataChan := make(chan string)

	var receivedRequestID string

	km := keyManagerFactory()

	// api
	api := New(&testUpStream{
		sendFn: func(data string) {
			dataChan <- data
		},
	}, km)

	go func() {
		select {
		case data := <-dataChan:
			req := &pb.Request{}
			if err := proto.Unmarshal([]byte(data), req); err != nil {
				panic(err)
			}
			receivedRequestID = req.RequestID
			out := api.Respond(req.RequestID, &pb.Response{}, nil, time.Second)
			if out != nil {
				panic("expected nil but got: " + out.Error())
			}
		}
	}()

	resp, err := api.request(&pb.Request{}, time.Second)
	resp.Closer <- nil
	require.Nil(t, err)

}

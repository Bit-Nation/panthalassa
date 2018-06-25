package api

import (
	"sync"
	"testing"

	pb "github.com/Bit-Nation/panthalassa/api/pb"
	proto "github.com/golang/protobuf/proto"
	require "github.com/stretchr/testify/require"
	"time"
)

func TestAPI_ShowModal(t *testing.T) {

	c := make(chan string)

	api := &API{
		lock:     sync.Mutex{},
		requests: map[string]chan *Response{},
		client: &UpStreamTestImpl{
			f: func(data string) {
				c <- data
			},
		},
	}

	dAppApi := DApp{
		api: api,
	}

	go func() {

		select {
		case data := <-c:

			req := pb.Request{}
			requireNil(proto.Unmarshal([]byte(data), &req))

			if req.ShowModal.Title != "Request Money" {
				panic("Expected title to be 'Request Money'")
			}

			if req.ShowModal.Layout != "{}" {
				panic("Expected layout to be '{}'")
			}

			err := api.Respond(req.RequestID, &pb.Response{}, nil, time.Second*3)
			if err != nil {
				panic("expected error to be nil")
			}

		}

	}()

	err := dAppApi.ShowModal("Request Money", "{}")
	require.Nil(t, err)

}

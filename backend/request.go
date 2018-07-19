package backend

import (
	"fmt"
	"time"

	bpb "github.com/Bit-Nation/protobuffers"
	uuid "github.com/satori/go.uuid"
)

type response struct {
	err  error
	resp *bpb.BackendMessage_Response
}

type request struct {
	Req      *bpb.BackendMessage_Request
	ReqID    string
	RespChan chan *response
}

// will request the chat backend
func (b *ServerBackend) request(req bpb.BackendMessage_Request, timeOut time.Duration) (*bpb.BackendMessage_Response, error) {

	respChan := make(chan *response)

	// request id
	id, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}

	// add request to stack
	b.outReqQueue <- &request{
		Req:      &req,
		RespChan: respChan,
		ReqID:    id.String(),
	}

	select {
	case resp := <-respChan:
		return resp.resp, resp.err
	case <-time.After(timeOut):
		// remove request
		delete(b.stack, id.String())
		return nil, fmt.Errorf("request timed out after %d", timeOut)
	}

}

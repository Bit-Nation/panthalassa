package backend

import (
	"fmt"
	"time"

	bpb "github.com/Bit-Nation/protobuffers"
	uuid "github.com/satori/go.uuid"
	"sync"
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

// stack of requests
type requestStack struct {
	stack map[string]chan *response
	lock  sync.Mutex
}

// add response channel to stack
func (s *requestStack) Add(reqID string, responseChan chan *response) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.stack[reqID] = responseChan
}

// remove request id from stack
func (s *requestStack) Remove(reqID string) {
	delete(s.stack, reqID)
}

// cut response chanel from stack
func (s *requestStack) Cut(reqID string) chan *response {
	responseChan := s.stack[reqID]
	delete(s.stack, reqID)
	return responseChan
}

// will request the chat backend
func (b *Backend) request(req bpb.BackendMessage_Request, timeOut time.Duration) (*bpb.BackendMessage_Response, error) {

	respChan := make(chan *response)

	// request id
	id, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}

	// add request to queue
	b.outReqQueue <- &request{
		Req:      &req,
		RespChan: respChan,
		ReqID:    id.String(),
	}

	select {
	case resp := <-respChan:
		return resp.resp, resp.err
	case <-time.After(timeOut):
		// remove request from stack
		b.stack.Remove(id.String())
		return nil, fmt.Errorf("request timed out after %d", timeOut)
	}

}

package backend

import (
	"sync"

	bpb "github.com/Bit-Nation/protobuffers"
	log "github.com/ipfs/go-log"
)

var logger = log.Logger("backend")

type ServerConfig struct {
	WebSocketUrl string
	BearerToken  string
}

type ServerBackend struct {
	transport Transport
	// all outgoing requests
	outReqQueue chan *request
	lock        sync.Mutex
	stack       map[string]chan *response
}

func (b *ServerBackend) Close() error {
	return b.transport.Close()
}

func NewServerBackend(trans Transport) (*ServerBackend, error) {

	b := &ServerBackend{
		outReqQueue: make(chan *request, 150),
		transport:   trans,
		lock:        sync.Mutex{},
		stack:       map[string]chan *response{},
	}

	// send outgoing requests to transport
	go func() {
		for {
			select {
			case req := <-b.outReqQueue:
				// send request
				err := b.transport.Send(&bpb.BackendMessage{
					RequestID: req.ReqID,
					Request:   req.Req,
				})
				// close response channel on error
				if err != nil {
					req.RespChan <- &response{
						err: err,
					}
					continue
				}
				// add response channel
				b.lock.Lock()
				b.stack[req.ReqID] = req.RespChan
				b.lock.Unlock()
			}
		}
	}()

	return b, nil

}

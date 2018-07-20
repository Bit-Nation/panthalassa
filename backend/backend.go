package backend

import (
	"sync"

	km "github.com/Bit-Nation/panthalassa/keyManager"
	bpb "github.com/Bit-Nation/protobuffers"
	log "github.com/ipfs/go-log"
	"time"
)

var logger = log.Logger("backend")

// IMPORTANT - the returned error will be send to the backend.
// Make sure it only return an error message that doesn't
// have private information
type RequestHandler *func(req *bpb.BackendMessage_Request) (*bpb.BackendMessage_Response, error)

type ServerConfig struct {
	WebSocketUrl string
	BearerToken  string
}

type Backend struct {
	transport Transport
	// all outgoing requests
	outReqQueue    chan *request
	lock           sync.Mutex
	stack          map[string]chan *response
	requestHandler []RequestHandler
	km             *km.KeyManager
	authenticated  bool
}

// Add request handler that will be executed
func (b *Backend) AddRequestHandler(handler RequestHandler) {
	b.lock.Lock()
	defer b.lock.Unlock()
	b.requestHandler = append(b.requestHandler, handler)
}

func (b *Backend) Start() error {
	return b.transport.Start()
}

func (b *Backend) Close() error {
	return b.transport.Close()
}

func NewServerBackend(trans Transport, km *km.KeyManager) (*ServerBackend, error) {

	b := &Backend{
		outReqQueue: make(chan *request, 150),
		transport:   trans,
		lock:        sync.Mutex{},
		stack:       map[string]chan *response{},
		km:          km,
	}

	// handle incoming message and iterate over
	// the registered message handlers
	trans.OnMessage(func(msg *bpb.BackendMessage) error {
		b.lock.Lock()
		defer b.lock.Unlock()
		for _, handler := range b.requestHandler {
			// handler
			h := *handler
			resp, err := h(msg.Request)
			// exit on error
			if err != nil {
				return b.transport.Send(&bpb.BackendMessage{
					RequestID: msg.RequestID,
					Error:     err.Error(),
				})
			}
			// if resp is nil we know that the handler didn't handle the request
			if resp == nil {
				continue
			}
			// send response
			err = b.transport.Send(&bpb.BackendMessage{
				Response:  resp,
				RequestID: msg.RequestID,
			})
			if err != nil {
				return err
			}
			// in the case this was a auth request we need to apply some special logic
			// this will only be executed when this message was a auth request
			if msg.Request != nil {
				if msg.Response.Auth != nil {
					b.authenticated = true
				}
			}
		}
		return nil
	})

	// register auth handler
	b.AddRequestHandler(b.auth())

	// send outgoing requests to transport
	go func() {
		for {

			// wait for authentication
			b.lock.Lock()
			if !b.authenticated {
				time.Sleep(time.Second * 1)
				b.lock.Unlock()
				continue
			}
			b.lock.Unlock()

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

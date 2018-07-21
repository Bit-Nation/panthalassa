package backend

import (
	"errors"
	"fmt"
	"sync"
	"time"

	km "github.com/Bit-Nation/panthalassa/keyManager"
	bpb "github.com/Bit-Nation/protobuffers"
	log "github.com/ipfs/go-log"
)

var logger = log.Logger("backend")

// IMPORTANT - the returned error will be send to the backend.
// Make sure it only return an error message that doesn't
// have private information
type RequestHandler *func(req *bpb.BackendMessage_Request) (*bpb.BackendMessage_Response, error)

type ResponseHandler *func(resp *bpb.BackendMessage_Response) (*response, error)

type ServerConfig struct {
	WebSocketUrl string
	BearerToken  string
}

type Backend struct {
	transport Transport
	// all outgoing requests
	outReqQueue     chan *request
	lock            sync.Mutex
	stack           requestStack
	requestHandler  []RequestHandler
	responseHandler []ResponseHandler
	km              *km.KeyManager
	authenticated   bool
}

// Add request handler that will be executed
func (b *Backend) AddRequestHandler(handler RequestHandler) {
	b.lock.Lock()
	defer b.lock.Unlock()
	b.requestHandler = append(b.requestHandler, handler)
}

// add response handler
func (b *Backend) AddResponseHandler(handler ResponseHandler) {
	b.lock.Lock()
	defer b.lock.Unlock()
	b.responseHandler = append(b.responseHandler, handler)
}

func (b *Backend) Start() error {
	return b.transport.Start()
}

func (b *Backend) Close() error {
	return b.transport.Close()
}

func NewServerBackend(trans Transport, km *km.KeyManager) (*Backend, error) {

	b := &Backend{
		outReqQueue: make(chan *request, 150),
		transport:   trans,
		lock:        sync.Mutex{},
		stack: requestStack{
			stack: map[string]chan *response{},
			lock:  sync.Mutex{},
		},
		km: km,
	}

	// handle incoming message and iterate over
	// the registered message handlers
	trans.OnMessage(func(msg *bpb.BackendMessage) error {
		b.lock.Lock()
		defer b.lock.Unlock()

		// handle requests
		if msg.Request != nil {
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
				if msg.Response.Auth != nil {
					b.authenticated = true
				}
			}
		}

		// handle responses
		if msg.Response != nil {

			resp := msg.Response
			reqID := msg.RequestID

			// err will be != nil in the case of no response channel
			respChan := b.stack.Cut(reqID)
			if respChan == nil {
				return fmt.Errorf("failed to fetch response channel for id: %s", msg.RequestID)
			}

			// send error from response to request channel
			if msg.Error != "" {
				respChan <- &response{
					err: errors.New(msg.Error),
				}
				return nil
			}
			for _, handler := range b.responseHandler {

				// if an handler return nil for an internalResponse
				h := *handler
				internalResponse, err := h(resp)
				if err != nil {
					return err
				}

				// if we got a response we send it to the response channel
				if internalResponse != nil {
					respChan <- internalResponse
					return nil
				}

			}

		}

		return nil
	})

	// register request handlers
	authHandler := b.auth
	b.AddRequestHandler(&authHandler)

	// register signing key response handler
	signingKeyHandler := func(resp *bpb.BackendMessage_Response) (*response, error) {
		if resp.SignedPreKey != nil {
			return &response{resp: resp}, nil
		}
		return nil, nil
	}
	b.AddResponseHandler(&signingKeyHandler)

	// register pre key bundle response handler
	preKeyBundleHandler := func(resp *bpb.BackendMessage_Response) (*response, error) {
		if resp.PreKeyBundle != nil {
			return &response{resp: resp}, nil
		}
		return nil, nil
	}
	b.AddResponseHandler(&preKeyBundleHandler)

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
				// add response channel
				b.stack.Add(req.ReqID, req.RespChan)
				// send request
				go func() {
					err := b.transport.Send(&bpb.BackendMessage{
						RequestID: req.ReqID,
						Request:   req.Req,
					})
					// close response channel on error
					if err != nil {
						req.RespChan <- &response{
							err: err,
						}
					}
				}()
			}
		}
	}()

	return b, nil

}

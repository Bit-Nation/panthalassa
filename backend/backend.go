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
type RequestHandler func(req *bpb.BackendMessage_Request) (*bpb.BackendMessage_Response, error)

type ServerConfig struct {
	WebSocketUrl string
	BearerToken  string
}

type Backend struct {
	transport Transport
	// all outgoing requests
	outReqQueue   chan *request
	stack         requestStack
	km            *km.KeyManager
	closer        chan struct{}
	addReqHandler chan RequestHandler
	reqHandlers   chan chan []RequestHandler
	authenticated chan chan bool
	authenticate  chan bool
}

// Add request handler that will be executed
func (b *Backend) AddRequestHandler(handler RequestHandler) {
	b.addReqHandler <- handler
}

func (b *Backend) Close() error {
	b.closer <- struct{}{}
	close(b.addReqHandler)
	close(b.reqHandlers)
	close(b.authenticate)
	close(b.authenticated)
	return b.transport.Close()
}

func NewBackend(trans Transport, km *km.KeyManager) (*Backend, error) {

	b := &Backend{
		transport:   trans,
		outReqQueue: make(chan *request, 150),
		stack: requestStack{
			stack: map[string]chan *response{},
			lock:  sync.Mutex{},
		},
		km:            km,
		closer:        make(chan struct{}, 1),
		addReqHandler: make(chan RequestHandler),
		reqHandlers:   make(chan chan []RequestHandler),
		authenticated: make(chan chan bool),
		authenticate:  make(chan bool),
	}

	// backend state
	go func() {

		reqHandlers := []RequestHandler{}
		authenticated := false

		for {

			// exist if channels have been closed
			if b.addReqHandler == nil || b.reqHandlers == nil || b.authenticated == nil || b.authenticate == nil {
				return
			}

			select {
			case rh := <-b.addReqHandler:
				reqHandlers = append(reqHandlers, rh)
			case respChan := <-b.reqHandlers:
				respChan <- reqHandlers
			case a := <-b.authenticate:
				authenticated = a
			case respChan := <-b.authenticated:
				respChan <- authenticated
			}
		}

	}()

	// handle received protobuf messages
	go func() {

		for {

			msg, err := trans.NextMessage()
			if err != nil {
				logger.Error(err)
				continue
			}

			// make sure we don't get a response & a request at the same time
			// we don't accept it. It's invalid!
			if msg.Request != nil && msg.Response != nil {
				logger.Error(errors.New("a message can’t have a response and a request at the same time"))
				continue
			}

			// handle requests
			if msg.Request != nil {
				requestHandled := false
				// ask the state for the request handlers
				reqHandlersChan := make(chan []RequestHandler)
				b.reqHandlers <- reqHandlersChan
				for _, handler := range <-reqHandlersChan {
					// handler
					h := handler
					resp, err := h(msg.Request)
					// exit on error
					if err != nil {
						b.transport.Send(&bpb.BackendMessage{
							RequestID: msg.RequestID,
							Error:     err.Error(),
						})
						continue
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
					requestHandled = true
					if err != nil {
						logger.Error(err)
						continue
					}

				}

				// If request was successfully handled we don't need to handle that message further
				if requestHandled {
					continue
				}
			}

			// handle responses
			if msg.Response != nil {

				resp := msg.Response

				// in the case this was a auth request we need to apply some special logic
				// this will only be executed when this message was a auth request
				if resp.Auth != nil {
					logger.Debug("[panthalassa] Recieved auth successful response")
					b.authenticate <- true
					continue
				}

				reqID := msg.RequestID

				// err will be != nil in the case of no response channel
				respChan := b.stack.Cut(reqID)
				if respChan == nil {
					logger.Error(fmt.Errorf("failed to fetch response channel for id: %s", msg.RequestID))
				}

				// send error from response to request channel
				if msg.Error != "" {
					respChan <- &response{
						err: errors.New(msg.Error),
					}
					continue
				}

				// send received response to response channel
				respChan <- &response{
					resp: resp,
				}

			}

			logger.Warning("dropping message: %s", msg)

		}

	}()

	// auth request handler
	b.AddRequestHandler(b.auth)

	// send outgoing requests to transport
	go func() {
		for {
			select {
			case <-b.closer:
				return
			case req := <-b.outReqQueue:
				authCheck := make(chan bool)
				b.authenticated <- authCheck
				// wait for authentication
				if !<-authCheck {
					time.Sleep(time.Second * 0.5)
					b.outReqQueue <- req
					close(authCheck)
					continue
				}

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

package backend

import (
	"crypto/rand"
	"errors"
	"fmt"
	"sync"
	"time"

	prekey "github.com/Bit-Nation/panthalassa/chat/prekey"
	db "github.com/Bit-Nation/panthalassa/db"
	km "github.com/Bit-Nation/panthalassa/keyManager"
	bpb "github.com/Bit-Nation/protobuffers"
	x3dh "github.com/Bit-Nation/x3dh"
	proto "github.com/gogo/protobuf/proto"
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
	outReqQueue         chan *request
	stack               requestStack
	km                  *km.KeyManager
	closer              chan struct{}
	addReqHandler       chan RequestHandler
	reqHandlers         chan chan []RequestHandler
	signedPreKeyStorage db.SignedPreKeyStorage
}

// Add request handler that will be executed
func (b *Backend) AddRequestHandler(handler RequestHandler) {
	b.addReqHandler <- handler
}

func (b *Backend) Close() error {
	b.closer <- struct{}{}
	err := b.transport.Close()
	if err != nil {
		return err
	}
	close(b.addReqHandler)
	close(b.reqHandlers)
	return nil
}

func NewBackend(trans Transport, km *km.KeyManager, signedPreKeyStorage db.SignedPreKeyStorage) (*Backend, error) {

	b := &Backend{
		transport:   trans,
		outReqQueue: make(chan *request, 150),
		stack: requestStack{
			stack: map[string]chan *response{},
			lock:  sync.Mutex{},
		},
		km:                  km,
		closer:              make(chan struct{}, 1),
		addReqHandler:       make(chan RequestHandler),
		reqHandlers:         make(chan chan []RequestHandler),
		signedPreKeyStorage: signedPreKeyStorage,
	}

	// backend state
	go func() {

		reqHandlers := []RequestHandler{}

		for {
			select {
			case <-b.closer:
				return
			case rh := <-b.addReqHandler:
				reqHandlers = append(reqHandlers, rh)
			case respChan := <-b.reqHandlers:
				respChan <- reqHandlers
			}
		}

	}()

	// Uploading of signedPeyKey
	go func() {
		// when we have a singed pre key we just want to continue
		if len(b.signedPreKeyStorage.All()) > 0 {
			return
		}

		c25519 := x3dh.NewCurve25519(rand.Reader)
		signedPreKeyPair, err := c25519.GenerateKeyPair()
		if err != nil {
			logger.Error(err)
			return
		}

		signedPreKey := prekey.PreKey{}
		signedPreKey.PrivateKey = signedPreKeyPair.PrivateKey
		signedPreKey.PublicKey = signedPreKeyPair.PublicKey
		if err := signedPreKey.Sign(*b.km); err != nil {
			logger.Error(err)
			return
		}
		protoSignedPreKey, err := signedPreKey.ToProtobuf()
		if err != nil {
			logger.Error(err)
			return
		}
		if err := signedPreKey.Sign(*b.km); err != nil {
			logger.Error(err)
			return
		}
		signedPreKeyBytes, err := proto.Marshal(&protoSignedPreKey)
		if err != nil {
			logger.Error(err)
			return
		}

		_, err = b.request(bpb.BackendMessage_Request{
			SignedPreKey: signedPreKeyBytes,
		}, time.Second*10)

		if err != nil {
			logger.Error(err)
			return
		}

		if err := b.signedPreKeyStorage.Put(signedPreKeyPair); err != nil {
			logger.Error(err)
			return
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
			} else {

				resp := msg.Response

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

				// Message handled
				continue
			}

			logger.Warning("dropping message: %s", msg)

		}

	}()

	// send outgoing requests to transport
	go func() {
		for {
			select {
			case <-b.closer:
				return
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

package api

import (
	"errors"
	"fmt"
	"sync"
	"time"

	pb "github.com/Bit-Nation/panthalassa/api/pb"
	"github.com/Bit-Nation/panthalassa/keyManager"
	proto "github.com/golang/protobuf/proto"
	log "github.com/ipfs/go-log"
	uuid "github.com/satori/go.uuid"
)

var logger = log.Logger("api")

type UpStream interface {
	Send(data string)
}

// Create new api with given client
func New(client UpStream, km *keyManager.KeyManager) *API {

	a := &API{
		lock:     sync.Mutex{},
		requests: map[string]chan *Response{},
		client:   client,
	}

	a.drKeyStoreApi = DoubleRatchetKeyStoreApi{
		api: a,
		km:  km,
	}

	a.dAppApi = DAppApi{
		api: a,
	}

	return a

}

type API struct {
	drKeyStoreApi DoubleRatchetKeyStoreApi
	dAppApi       DAppApi
	lock          sync.Mutex
	requests      map[string]chan *Response
	client        UpStream
}

// This represent an api response
// It's important to close the response as soon as possible
// If the response is not closed by sending an error / nil to the Closer
// The response will time out for the client as it doesn't receive a response
type Response struct {
	Msg    *pb.Response
	Error  error
	Closer chan error
}

// send a response for a received api request
func (a *API) Respond(id string, resp *pb.Response, passedErr error, timeOut time.Duration) error {

	req, err := a.cutRequest(id)
	if err != nil {
		return err
	}

	closer := make(chan error)

	req <- &Response{
		Msg:    resp,
		Error:  passedErr,
		Closer: closer,
	}

	select {
	case err := <-closer:
		return err
	case <-time.After(timeOut):
		return errors.New(fmt.Sprintf("Response for id: %s timed out", id))
	}

}

// add request to request stack
func (a *API) addRequest(req *pb.Request) <-chan *Response {

	respChan := make(chan *Response)

	// add request to stack
	a.lock.Lock()
	a.requests[req.RequestID] = respChan
	a.lock.Unlock()

	return respChan

}

// fetch a request from stack and deletes it
func (a *API) cutRequest(id string) (chan *Response, error) {

	a.lock.Lock()
	req, exist := a.requests[id]
	// in case it exist we want to delete it to free the map
	if exist {
		delete(a.requests, id)
	}
	a.lock.Unlock()

	if !exist {
		return nil, errors.New(fmt.Sprintf("couldn't find request for ID: %s", id))
	}

	return req, nil

}

// send a request to the client
func (a *API) request(req *pb.Request, timeOut time.Duration) (*Response, error) {

	// create request ID
	requestId, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}
	req.RequestID = requestId.String()

	// add request to stack
	reqChan := a.addRequest(req)

	// serialize and send to api
	rawData, err := proto.Marshal(req)
	if err != nil {
		return nil, err
	}
	logger.Info("going to send this: " + string(rawData) + " to upstream")
	go a.client.Send(string(rawData))

	// wait for the response
	// or time out
	select {
	case res := <-reqChan:
		// close the response here
		// since we got an error back
		// for our request
		if res.Error != nil {
			res.Closer <- nil
		}
		return nil, res.Error
	case <-time.After(timeOut):
		// remove request from stack
		_, err := a.cutRequest(requestId.String())
		if err != nil {
			logger.Error(err)
		}
		return nil, errors.New(fmt.Sprintf("request timeout for ID: %s", requestId))
	}

}

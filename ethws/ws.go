package ethws

import (
	"encoding/json"
	"fmt"
	"time"

	wsg "github.com/gorilla/websocket"
	log "github.com/ipfs/go-log"
)

const rpcVersion = "2.0"

var logger = log.Logger("ethws")

type Config struct {
	Retry time.Duration
	WSUrl string
}

type Request struct {
	ID      int64         `json:"id"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
	JsonRPC string        `json:"jsonrpc"`
}

func (r *Request) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type Response struct {
	JsonRPC  string      `json:"jsonrpc"`
	RPCError *Error      `json:"error,omitempty"`
	Result   interface{} `json:"result"`
	ID       int64       `json:"id"`
	error    error
}

func (r *Response) Error() error {
	return r.error
}

type EthereumWS struct {
	state State
	// requests that need to be send
	requestQueue chan Request
	requests     map[int64]chan<- Response
	conn         *wsg.Conn
}

type State struct {
	addRequestIfNotExist chan StateObject
	getRequest           chan StateObject
	deleteRequest        chan StateObject
	response             chan StateObject
}

type StateObject struct {
	id int64
	c  chan<- Response
}

// start state machine
func (ws *EthereumWS) State() {
	for {
		select {
		case getRequest := <-ws.state.getRequest:
			// if the request exists, respond with the id and the response channel
			if respChan, exist := ws.requests[getRequest.id]; exist {
				ws.state.response <- StateObject{getRequest.id, respChan}
				break
			}
			// if the request doesn't exist, respond with 0 and nil
			ws.state.response <- StateObject{0, nil}
		case addRequest := <-ws.state.addRequestIfNotExist:
			// if the request exists, don't add it to the map, respond with the id and the response channel
			if respChan, exist := ws.requests[addRequest.id]; exist {
				ws.state.response <- StateObject{addRequest.id, respChan}
				break
			}
			// if the request doesn't exist, add it to the map, respond with 0 and nil
			ws.requests[addRequest.id] = addRequest.c
			ws.state.response <- StateObject{0, nil}
		case deleteRequest := <-ws.state.deleteRequest:
			delete(ws.requests, deleteRequest.id)
			ws.state.response <- StateObject{0, nil}
		} // select
	} // infinite for
} // func (ws *EthereumWS) State()

// send an request to ethereum network
func (ws *EthereumWS) SendRequest(r Request) (<-chan Response, error) {

	c := make(chan Response)

	// add request to stack
	for {
		id := time.Now().UnixNano()
		ws.state.addRequestIfNotExist <- StateObject{id, c}
		stateResponse := <-ws.state.response
		// if the id didn't exist
		if stateResponse.id == 0 {
			r.ID = id
			break
		}
	}

	// send request to queue
	ws.requestQueue <- r
	return c, nil

}

// create new ethereum websocket
func New(conf Config) *EthereumWS {

	startSendWorker := make(chan bool)
	startReadWorker := make(chan bool)

	etws := &EthereumWS{
		state: State{
			addRequestIfNotExist: make(chan StateObject),
			getRequest:           make(chan StateObject),
			deleteRequest:        make(chan StateObject),
			response:             make(chan StateObject),
		},
		requestQueue: make(chan Request, 1000),
		requests:     map[int64]chan<- Response{},
	}
	go etws.State()
	// worker that sends the requests
	go func() {

		// wait for connection
		<-startSendWorker
		// send requests
		for {
			select {
			case req := <-etws.requestQueue:

				// send request
				if err := etws.conn.WriteJSON(req); err != nil {
					logger.Error(err)

					etws.state.getRequest <- StateObject{req.ID, nil}
					stateResponse := <-etws.state.response
					respChan := stateResponse.c

					respChan <- Response{error: err}
				}
			}
		}

	}()

	// worker that read response from websocket
	go func() {

		// wait to start worker
		<-startReadWorker
		for {

			// read message
			_, response, err := etws.conn.ReadMessage()
			if err != nil {
				logger.Error(err)
				continue
			}
			// unmarshal
			var resp Response
			if err := json.Unmarshal(response, &resp); err != nil {
				logger.Error(err)
				continue
			}

			// get response channel
			etws.state.getRequest <- StateObject{resp.ID, nil}
			stateResponse := <-etws.state.response
			if stateResponse.id == 0 {
				logger.Error(fmt.Sprintf("failed to get response channel for ID: %d", stateResponse.id))
				continue
			}
			etws.state.deleteRequest <- StateObject{resp.ID, nil}
			_ = <-etws.state.response
			// send response
			respChan := stateResponse.c
			respChan <- resp

		}

	}()

	// connect to ethereum node
	go func() {

		// try to connect till success
		for {
			co, _, err := wsg.DefaultDialer.Dial(conf.WSUrl, nil)
			if err == nil {
				etws.conn = co
				break
			}
			logger.Error(err)
			// wait a bit. We don't want to stress the endpoint
			time.Sleep(conf.Retry)
		}

		// signal the workers to start
		startReadWorker <- true
		startSendWorker <- true

	}()

	return etws
}

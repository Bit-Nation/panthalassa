package ethws

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"sync"
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
	ID      uint          `json:"id"`
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
	RPCError Error       `json:"error"`
	Result   interface{} `json:"result"`
	ID       uint        `json:"id"`
	Error    error
}

type EthereumWS struct {
	lock sync.Mutex
	// requests that need to be send
	requestQueue chan Request
	requests     map[uint]chan<- Response
	conn         *wsg.Conn
}

// send an request to ethereum network
func (ws *EthereumWS) SendRequest(r Request) (<-chan Response, error) {

	c := make(chan Response)

	// add request to stack
	ws.lock.Lock()
	defer ws.lock.Unlock()
	for {
		id := rand.Uint64()
		if _, exist := ws.requests[uint(id)]; !exist {
			r.ID = uint(id)
			ws.requests[uint(id)] = c
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
		lock:         sync.Mutex{},
		requestQueue: make(chan Request, 1000),
		requests:     map[uint]chan<- Response{},
	}

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

					etws.lock.Lock()
					respChan := etws.requests[req.ID]
					etws.lock.Unlock()

					respChan <- Response{Error: err}
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
			etws.lock.Lock()
			respChan, exist := etws.requests[resp.ID]
			if !exist {
				logger.Error(fmt.Sprintf("failed to get response channel for ID: %d", resp.ID))
				continue
			}
			delete(etws.requests, resp.ID)
			etws.lock.Unlock()

			// send response
			respChan <- resp

		}

	}()

	// connect to ethereum node
	go func() {

		// try to connect till success
		for {
			co, _, err := wsg.DefaultDialer.Dial(conf.WSUrl, nil)
			if err != nil {
				logger.Error(err)
				// wait a bit. We don't want to stress the endpoint
				time.Sleep(conf.Retry)
				continue
			}
			etws.conn = co
			break
		}

		// signal the workers to start
		startReadWorker <- true
		startSendWorker <- true

	}()

	return etws
}

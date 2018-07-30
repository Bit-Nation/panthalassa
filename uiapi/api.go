package stapi

// api for updating the UI only
// changes in our state will be send over this api

import (
	"encoding/json"

	log "github.com/ipfs/go-log"
)

var logger = log.Logger("state-api")

type UpStream interface {
	Send(data string)
}

type call struct {
	Name    string                 `json:"name"`
	Payload map[string]interface{} `json:"payload"`
}

type Api struct {
	us     UpStream
	closer chan struct{}
	stack  chan call
}

func (a *Api) Close() error {
	a.closer <- struct{}{}
	return nil
}

// send to api
func (a *Api) Send(name string, payload map[string]interface{}) {
	a.stack <- call{
		Name:    name,
		Payload: payload,
	}
}

func New(us UpStream) *Api {

	api := &Api{
		us:     us,
		closer: make(chan struct{}, 1),
		stack:  make(chan call, 200),
	}

	go func() {

		for {
			select {
			case <-api.closer:
				return
			case call := <-api.stack:
				rawCall, err := json.Marshal(call)
				if err != nil {
					logger.Error(err)
					continue
				}
				api.us.Send(string(rawCall))

			}
		}

	}()

	return api
}

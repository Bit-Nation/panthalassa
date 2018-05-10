package device_api

import (
	"encoding/json"
	"github.com/Bit-Nation/panthalassa/api/device/rpc"
)

type UpStream interface {
	Send(data string)
}

type apiCall struct {
	Type string `json:"type"`
	Id   uint32 `json:"id"`
	Data string `json:"data"`
}

func (c *apiCall) Marshal() ([]byte, error) {
	return json.Marshal(c)
}

type Response struct {
	Content string
	Closer  chan error
}

type Api struct {
	device UpStream
	state  *State
}

func New(deviceInterface UpStream) *Api {

	api := Api{
		state:  newState(),
		device: deviceInterface,
	}

	return &api
}

func (a *Api) Send(call rpc.JsonRPCCall) (chan Response, error) {

	//Validate call
	if err := call.Valid(); err != nil {
		return nil, err
	}

	//Get call data
	callContent, err := call.Data()
	if err != nil {
		return nil, err
	}

	//Create internal representation
	c := apiCall{
		Type: call.Type(),
		Data: callContent,
	}
	respChan := make(chan Response, 1)
	c.Id = a.state.Add(respChan)

	//Marshal the call data
	callData, err := c.Marshal()
	if err != nil {
		return nil, err
	}

	//Send json rpc call to device
	a.device.Send(string(callData))

	return respChan, nil

}

func (a *Api) Receive(id uint32, data string) error {

	//Get the response channel
	resp, err := a.state.Cut(id)
	if err != nil {
		return err
	}

	//Closer
	closer := make(chan error)

	//Send response to response channel
	resp <- Response{
		Content: data,
		Closer:  closer,
	}

	return <-closer

}

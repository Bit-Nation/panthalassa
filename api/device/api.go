package device_api

import (
	"encoding/json"
	"errors"

	"fmt"
	"github.com/Bit-Nation/panthalassa/api/device/rpc"
	log "github.com/ipfs/go-log"
	valid "gopkg.in/asaskevich/govalidator.v4"
)

var logger = log.Logger("device_api")

type UpStream interface {
	//Send data to client
	Send(data string)
}

type ApiCall struct {
	Type string `json:"type"`
	Id   string `json:"id"`
	Data string `json:"data"`
}

func (c *ApiCall) Marshal() ([]byte, error) {
	return json.Marshal(c)
}

func UnmarshalApiCall(call string) (ApiCall, error) {

	var apiCall ApiCall

	err := json.Unmarshal([]byte(call), &apiCall)

	return apiCall, err
}

type rawResponse struct {
	Error   string `json:"error",valid:"string,optional"`
	Payload string `json:"payload",valid:"string,optional"`
}

type Response struct {
	Error   error
	Payload string
	Closer  chan error
}

func (r *Response) Close(err error) {
	r.Closer <- err
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

//Send a call to the api
func (a *Api) Send(call rpc.JsonRPCCall) (<-chan Response, error) {

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
	c := ApiCall{
		Type: call.Type(),
		Data: callContent,
	}
	respChan := make(chan Response, 1)
	c.Id, err = a.state.Add(respChan)
	if err != nil {
		return nil, err
	}

	//Marshal the call data
	callData, err := c.Marshal()
	if err != nil {
		return nil, err
	}

	//Send json rpc call to device
	// @todo maybe it's worth making this sync and let the caller decide
	// @todo if they want to wait till the data has been send
	go a.device.Send(string(callData))

	return respChan, nil

}

// @todo at the moment the fetched response channel will never close in case there we return earlier with an error
func (a *Api) Receive(id string, data string) error {

	logger.Debug(fmt.Sprintf("Got response for request (%s) - with data: %s", id, data))

	// get the response channel
	resp, err := a.state.Cut(id)

	if err != nil {
		// only try to send a response if a channel exist
		if resp != nil {
			resp <- Response{Error: err}
		}
		return err
	}

	// closer
	closer := make(chan error)

	// decode raw response
	var rr rawResponse
	if err := json.Unmarshal([]byte(data), &rr); err != nil {
		resp <- Response{
			Error: err,
		}
		return err
	}

	// validate raw response
	_, err = valid.ValidateStruct(rr)
	if err != nil {
		resp <- Response{
			Error: err,
		}
		return err
	}

	// construct response
	r := Response{
		Error:   err,
		Payload: rr.Payload,
		Closer:  closer,
	}
	if rr.Error != "" {
		r.Error = errors.New(rr.Error)
	}

	// send response to response channel
	resp <- r

	logger.Debug("send response", r)

	return <-closer

}

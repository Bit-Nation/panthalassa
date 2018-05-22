package mesh

import (
	"encoding/json"
	"errors"
	"fmt"

	api "github.com/Bit-Nation/panthalassa/api/device"
	ds "github.com/ipfs/go-datastore"
	dsq "github.com/ipfs/go-datastore/query"
)

type DataStore struct {
	deviceApi *api.Api
}

func NewDataStore(api *api.Api) ds.Batching {
	return &DataStore{
		deviceApi: api,
	}
}

func (d *DataStore) Put(key ds.Key, value interface{}) error {
	logger.Info(fmt.Sprintf("put key: %s and value: %s in datastore", key.String(), value))

	call := DHTPutCall{
		Key:   key.String(),
		Value: value,
	}

	respChan, err := d.deviceApi.Send(&call)
	if err != nil {
		return err
	}

	// wait for the response from the client / device
	resp := <-respChan

	// close it here since we don't need it anymore
	resp.Close(nil)

	return resp.Error

}

func (d *DataStore) Get(key ds.Key) (value interface{}, err error) {

	logger.Info(fmt.Sprintf("fetch value for key: %s", key.String()))

	call := DHTGetCall{
		Key: key.String(),
	}

	respChan, err := d.deviceApi.Send(&call)
	if err != nil {
		return nil, err
	}

	// wait for the response from the client / device
	resp := <-respChan

	// exit if there is an error in the response
	if err != resp.Error {
		resp.Close(nil)
		return nil, resp.Error
	}

	payload := &struct {
		Value interface{} `json:"value"`
	}{}

	if err := json.Unmarshal([]byte(resp.Payload), payload); err != nil {
		// send error back to client
		resp.Close(err)
		return nil, errors.New(fmt.Sprintf("failed to unmarshal response from client for DHT:PUT query (raw response: %s)", resp.Payload))
	}

	resp.Close(err)

	return payload.Value, nil

}

func (d *DataStore) Has(key ds.Key) (exists bool, err error) {

	logger.Info(fmt.Sprintf("check if key is mapped to an value: %s exist", key.String()))

	// device api call
	call := DHTHasCall{
		Key: key.String(),
	}

	// send to device
	respChan, err := d.deviceApi.Send(&call)
	if err != nil {
		return false, err
	}

	// wait for a response
	resp := <-respChan

	// exit with error if there is one in the response
	if resp.Error != nil {
		resp.Close(nil)
		return false, resp.Error
	}

	payload := &struct {
		Exist bool `json:"exist"`
	}{}

	// validate json response
	// @todo we also need to check the response json schema it self since the struct is initialized with it's default values.
	// @todo That could be a problem in case we get an invalid response
	if err := json.Unmarshal([]byte(resp.Payload), payload); err != nil {
		// send error back to client
		resp.Close(err)
		return false, errors.New(fmt.Sprintf("failed to unmarshal response from client for DHT:HAS query (raw response: %s)", resp.Payload))
	}

	// close the response
	resp.Close(nil)

	return payload.Exist, nil
}

func (d *DataStore) Delete(key ds.Key) error {

	logger.Info(fmt.Sprintf("check if key is mapped to an value: %s exist", key.String()))

	call := &DHTDeleteCall{
		Key: key.String(),
	}

	respChan, err := d.deviceApi.Send(call)
	if err != nil {
		return err
	}

	resp := <-respChan

	resp.Close(nil)

	return resp.Error

}

func (d *DataStore) Query(q dsq.Query) (dsq.Results, error) {

	panic("not implemented")

	return nil, nil
}

func (d *DataStore) Batch() (ds.Batch, error) {
	panic("batch is not implemented")
	return nil, nil
}

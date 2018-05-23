package main

import (
	"encoding/json"
	"fmt"

	panthalassa "github.com/Bit-Nation/panthalassa"
	jsonDB "github.com/nanobox-io/golang-scribble"
)

type apiCall struct {
	Type string `json:"type"`
	ID   uint32 `json:"id"`
	Data string `json:"data"`
}

type response struct {
	Error   string `json:"error"`
	Payload string `json:"payload"`
}

func unmarshalCall(call string) (apiCall, error) {

	var c apiCall

	if err := json.Unmarshal([]byte(call), &c); err != nil {
		return apiCall{}, nil
	}

	return c, nil

}

type Store struct {
	Account Account
	DB      *jsonDB.Driver
}

func (s Store) SendResponse(id uint32, r response) {

	data, err := json.Marshal(r)
	if err != nil {
		logger.Error(err)
		return
	}

	logger.Info("answer request: ", id, " data: ", string(data))
	err = panthalassa.SendResponse(id, string(data))
	if err != nil {
		logger.Error(err)
		return
	}
	logger.Info("finished request: ", id)

}

func (s *Store) Send(data string) {

	call, err := unmarshalCall(data)
	if err != nil {
		panic(err)
	}

	switch call.Type {
	case "CONTACT:LIST":
		contacts, err := s.DB.ReadAll("contact")
		if err != nil {
			s.SendResponse(call.ID, response{
				Error: err.Error(),
			})
			return
		}
		jContacts, err := json.Marshal(contacts)
		if err != nil {
			s.SendResponse(call.ID, response{
				Error: err.Error(),
			})
			return
		}
		s.SendResponse(call.ID, response{
			Payload: string(jContacts),
		})
	default:
		s.SendResponse(call.ID, response{
			Error: fmt.Sprintf("couldn't process call with type: %s", call.Type),
		})
	}

}

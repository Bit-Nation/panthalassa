package mesh

import (
	"encoding/json"
	valid "github.com/asaskevich/govalidator"
)

type DHTPutCall struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}

func (c *DHTPutCall) Type() string {
	return "DHT:PUT"
}
func (c *DHTPutCall) Data() (string, error) {

	data, err := json.Marshal(c)
	if err != nil {
		return "", err
	}

	return string(data), err

}
func (c *DHTPutCall) Valid() error {
	return nil
}

type DHTGetCall struct {
	Key string `json:"key",valid:"required"`
}

func (c *DHTGetCall) Type() string {
	return "DHT:GET"
}
func (c *DHTGetCall) Data() (string, error) {

	data, err := json.Marshal(c)
	if err != nil {
		return "", err
	}

	return string(data), err

}
func (c *DHTGetCall) Valid() error {
	_, err := valid.ValidateStruct(c)
	return err
}

type DHTHasCall struct {
	Key string `json:"key",valid:"required"`
}

func (d *DHTHasCall) Type() string {
	return "DHT:HAS"
}
func (d *DHTHasCall) Data() (string, error) {

	raw, err := json.Marshal(d)
	if err != nil {
		return "", err
	}
	return string(raw), nil

}
func (d *DHTHasCall) Valid() error {

	_, err := valid.ValidateStruct(d)
	return err

}

type DHTDeleteCall struct {
	Key string `json:"key",valid:"required"`
}

func (d *DHTDeleteCall) Type() string {
	return "DHT:DELETE"
}
func (d *DHTDeleteCall) Data() (string, error) {

	raw, err := json.Marshal(d)
	if err != nil {
		return "", err
	}

	return string(raw), err

}
func (d *DHTDeleteCall) Valid() error {

	_, err := valid.ValidateStruct(d)
	return err

}

package client

import (
	"encoding/json"
	"errors"

	aes "github.com/Bit-Nation/panthalassa/crypto/aes"
	keyManager "github.com/Bit-Nation/panthalassa/keyManager"
	dr "github.com/tiabc/doubleratchet"
)

var invalidKeyLen = errors.New("a key must have the length 32 bytes and must be in hex format")

type DRKeyStoreGetCall struct {
	Key    string `json:"key"`
	MsgNum uint   `json:"msg_num"`
}

func (c *DRKeyStoreGetCall) Type() string {
	return "DR:KEY_STORE:GET"
}

func (c *DRKeyStoreGetCall) Data() (string, error) {
	data, err := json.Marshal(c)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (c *DRKeyStoreGetCall) Valid() error {
	return nil
}

type DRKeyStoreGetResponse struct {
	Key          dr.Key
	EncryptedKey string `json:"key"`
	Ok           bool
}

func UnmarshalDRKeyStoreGetResponse(payload string, km *keyManager.KeyManager) (DRKeyStoreGetResponse, error) {

	var r DRKeyStoreGetResponse
	if err := json.Unmarshal([]byte(payload), &r); err != nil {
		return DRKeyStoreGetResponse{}, err
	}

	if r.EncryptedKey == "" {
		r.Ok = false
		return r, nil
	}

	keyCT, err := aes.Unmarshal([]byte(r.EncryptedKey))
	if err != nil {
		return DRKeyStoreGetResponse{}, err
	}

	plainSecret, err := km.AESDecrypt(keyCT)
	if err != nil {
		return DRKeyStoreGetResponse{}, err
	}

	if len(plainSecret) != 32 {
		return DRKeyStoreGetResponse{}, invalidKeyLen
	}

	var key [32]byte
	copy(key[:], plainSecret[:32])

	r.Key = key
	r.Ok = true

	return r, nil

}

type DRKeyStorePutCall struct {
	IndexKey         string         `json:"index_key"`
	MsgNumber        uint           `json:"msg_number"`
	DoubleRatchetKey aes.CipherText `json:"msg_key"`
}

func (c *DRKeyStorePutCall) Type() string {
	return "DR:KEY_STORE:PUT"
}

func (c *DRKeyStorePutCall) Data() (string, error) {
	data, err := json.Marshal(c)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (c *DRKeyStorePutCall) Valid() error {
	return nil
}

type DRKeyStoreDeleteMK struct {
	IndexKey  string `json:"index_key"`
	MsgNumber uint   `json:"msg_num"`
}

func (c *DRKeyStoreDeleteMK) Type() string {
	return "DR:KEY_STORE:DELETE_MESSAGE_KEY"
}

func (c *DRKeyStoreDeleteMK) Data() (string, error) {
	data, err := json.Marshal(c)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (c *DRKeyStoreDeleteMK) Valid() error {
	return nil
}

type DRKeyStoreDeleteIndexKey struct {
	IndexKey string `json:"index_key"`
}

func (c *DRKeyStoreDeleteIndexKey) Type() string {
	return "DR:KEY_STORE:DELETE_INDEX_KEY"
}

func (c *DRKeyStoreDeleteIndexKey) Data() (string, error) {
	data, err := json.Marshal(c)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (c *DRKeyStoreDeleteIndexKey) Valid() error {
	return nil
}

type DRKeyStoreCountCall struct {
	IndexKey string `json:"index_key"`
}

func (c *DRKeyStoreCountCall) Type() string {
	return "DR:KEY_STORE:COUNT_MESSAGES"
}

func (c *DRKeyStoreCountCall) Data() (string, error) {
	data, err := json.Marshal(c)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (c *DRKeyStoreCountCall) Valid() error {
	return nil
}

package client

import (
	"encoding/hex"
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

type DRKeyStoreFetchAllKeys struct{}

func (c *DRKeyStoreFetchAllKeys) Type() string {
	return "DR:KEY_STORE:FETCH_ALL_KEYS"
}

func (c *DRKeyStoreFetchAllKeys) Data() (string, error) {
	return "", nil
}

func (c *DRKeyStoreFetchAllKeys) Valid() error {
	return nil
}

func UnmarshalFetchAllKeysPayload(payload string, km *keyManager.KeyManager) (map[dr.Key]map[uint]dr.Key, error) {

	// unmarshal payload
	var raw map[string]map[uint]string
	if err := json.Unmarshal([]byte(payload), &raw); err != nil {
		return map[dr.Key]map[uint]dr.Key{}, err
	}

	allKeys := map[dr.Key]map[uint]dr.Key{}

	for k, v := range raw {

		rawIndexKey, err := hex.DecodeString(k)
		if err != nil {
			return map[dr.Key]map[uint]dr.Key{}, err
		}

		if len(rawIndexKey) != 32 {
			return nil, errors.New("an index key must be exactly 32 bytes long")
		}

		msgKeys := map[uint]dr.Key{}
		for msgNum, encryptedMsgKey := range v {

			// parse cipher text
			ct, err := aes.Unmarshal([]byte(encryptedMsgKey))
			if err != nil {
				return map[dr.Key]map[uint]dr.Key{}, err
			}

			// decrypt message key
			secret, err := km.AESDecrypt(ct)

			if len(rawIndexKey) != 32 {
				return nil, errors.New("an message key must be exactly 32 bytes long")
			}

			var msgKey dr.Key
			copy(msgKey[:], secret)

			msgKeys[msgNum] = msgKey

		}

		var indexKey dr.Key
		copy(indexKey[:], rawIndexKey)

		allKeys[indexKey] = msgKeys

	}

	return allKeys, nil

}

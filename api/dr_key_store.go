package api

import (
	"errors"
	"time"

	pb "github.com/Bit-Nation/panthalassa/api/pb"
	aes "github.com/Bit-Nation/panthalassa/crypto/aes"
	keyManager "github.com/Bit-Nation/panthalassa/keyManager"
	dr "github.com/tiabc/doubleratchet"
)

// This functions just proxy the call from the API main
// object down to the Key store api
func (a *API) Get(k dr.Key, msgNum uint) (mk dr.Key, ok bool) {
	return a.drKeyStoreApi.Get(k, msgNum)
}
func (a *API) Put(k dr.Key, msgNum uint, mk dr.Key) {
	a.drKeyStoreApi.Put(k, msgNum, mk)
}
func (a *API) DeleteMk(k dr.Key, msgNum uint) {
	a.drKeyStoreApi.DeleteMk(k, msgNum)
}
func (a *API) DeletePk(k dr.Key) {
	a.drKeyStoreApi.DeletePk(k)
}
func (a *API) Count(k dr.Key) uint {
	return a.drKeyStoreApi.Count(k)
}
func (a *API) All() map[dr.Key]map[uint]dr.Key {
	return a.drKeyStoreApi.All()
}

// Double ratched key store api
type DoubleRatchetKeyStoreApi struct {
	api *API
	km  *keyManager.KeyManager
}

// get a key by it's key and msg number
func (s *DoubleRatchetKeyStoreApi) Get(k dr.Key, msgNum uint) (mk dr.Key, ok bool) {

	req := pb.Request{
		DRKeyStoreGet: &pb.Request_DRKeyStoreGet{
			DrKey:         k[:],
			MessageNumber: uint64(msgNum),
		},
	}

	resp, err := s.api.request(&req, time.Second*8)
	if err != nil {
		logger.Error(err)
		return dr.Key{}, false
	}

	ct, err := aes.Unmarshal(resp.Msg.DRKeyStoreGet.MessageKey)
	if err != nil {
		logger.Error(err)
		resp.Closer <- err
		return dr.Key{}, false
	}

	messageKey, err := s.km.AESDecrypt(ct)
	if err != nil {
		logger.Error(err)
		resp.Closer <- err
		return dr.Key{}, false
	}

	if len(messageKey) != 32 {
		e := errors.New("a decrypted message key must have exactly 32 bytes")
		logger.Error(e)
		resp.Closer <- e
		return dr.Key{}, false
	}

	resp.Closer <- nil

	var msgKey dr.Key
	copy(msgKey[:], messageKey)

	return msgKey, true

}

// save message key (double ratchet key)
func (s *DoubleRatchetKeyStoreApi) Put(k dr.Key, msgNum uint, mk dr.Key) {

	ct, err := s.km.AESEncrypt(mk[:])
	if err != nil {
		logger.Error(err)
		return
	}

	rawCt, err := ct.Marshal()
	if err != nil {
		logger.Error(err)
		return
	}

	resp, err := s.api.request(&pb.Request{
		DRKeyStorePut: &pb.Request_DRKeyStorePut{
			Key:           k[:],
			MessageNumber: uint64(msgNum),
			MessageKey:    rawCt,
		},
	}, time.Second*8)

	if err != nil {
		logger.Error(err)
		return
	}

	resp.Closer <- nil

}

func (s *DoubleRatchetKeyStoreApi) DeleteMk(k dr.Key, msgNum uint) {

	resp, err := s.api.request(&pb.Request{
		DRKeyStoreDeleteMK: &pb.Request_DRKeyStoreDeleteMK{
			Key:    k[:],
			MsgNum: uint64(msgNum),
		},
	}, time.Second*8)

	if err != nil {
		logger.Error(err)
		return
	}

	resp.Closer <- nil

}

func (s *DoubleRatchetKeyStoreApi) DeletePk(k dr.Key) {

	resp, err := s.api.request(&pb.Request{
		DRKeyStoreDeleteKeys: &pb.Request_DRKeyStoreDeleteKeys{
			Key: k[:],
		},
	}, time.Second*8)

	if err != nil {
		logger.Error(err)
		return
	}

	resp.Closer <- nil

}

func (s *DoubleRatchetKeyStoreApi) Count(k dr.Key) uint {

	resp, err := s.api.request(&pb.Request{
		DRKeyStoreCount: &pb.Request_DRKeyStoreCount{
			Key: k[:],
		},
	}, time.Second*8)

	if err != nil {
		logger.Error(err)
		return 0
	}

	resp.Closer <- nil

	return uint(resp.Msg.DRKeyStoreCount.Count)

}

// @todo the all method is way to heavy. Long term we need to have another solution
func (s *DoubleRatchetKeyStoreApi) All() map[dr.Key]map[uint]dr.Key {

	resp, err := s.api.request(&pb.Request{
		DRKeyStoreAll: &pb.Request_DRKeyStoreAll{},
	}, time.Second*8)

	if err != nil {
		logger.Error(err)
		return map[dr.Key]map[uint]dr.Key{}
	}

	var keys = map[dr.Key]map[uint]dr.Key{}

	for _, k := range resp.Msg.DRKeyStoreAll.All {

		// exit if key len is incorrect
		if len(k.Key) != 32 {
			e := errors.New("got invalid key in All() (expected key len == 32 bytes)")
			logger.Error(e)
			resp.Closer <- e
			return map[dr.Key]map[uint]dr.Key{}
		}

		indexKey := dr.Key{}
		copy(indexKey[:], k.Key)

		messages := map[uint]dr.Key{}

		for msgNum, key := range k.MessageKeys {

			ct, err := aes.Unmarshal(key)
			if err != nil {
				resp.Closer <- err
				logger.Error(err)
				return map[dr.Key]map[uint]dr.Key{}
			}

			rawMsgKey, err := s.km.AESDecrypt(ct)
			if err != nil {
				resp.Closer <- err
				logger.Error(err)
				return map[dr.Key]map[uint]dr.Key{}
			}

			if len(rawMsgKey) != 32 {
				e := errors.New("got invalid key in All() (expected key len == 32 bytes)")
				logger.Error(e)
				resp.Closer <- e
				return map[dr.Key]map[uint]dr.Key{}
			}

			msgKey := dr.Key{}
			copy(msgKey[:], rawMsgKey)

			messages[uint(msgNum)] = msgKey
		}

		keys[indexKey] = messages

	}

	resp.Closer <- nil

	return keys

}

func NewDRKeyStoreApi(api *API, km *keyManager.KeyManager) *DoubleRatchetKeyStoreApi {
	return &DoubleRatchetKeyStoreApi{
		api: api,
		km:  km,
	}
}

package api

import (
	"errors"
	"time"
	
	keyManager "github.com/Bit-Nation/panthalassa/keyManager"
	dr "github.com/tiabc/doubleratchet"
	pb "github.com/Bit-Nation/panthalassa/api/pb"
	aes "github.com/Bit-Nation/panthalassa/crypto/aes"
)

type DoubleRatchetKeyStoreApi struct {
	api *API
	km  *keyManager.KeyManager
}

// get a key by it's key and msg number
func (s *DoubleRatchetKeyStoreApi) Get(k dr.Key, msgNum uint) (mk dr.Key, ok bool) {
	
	req := pb.Request{
		DRKeyStoreGet: &pb.Request_DRKeyStoreGet{
			DrKey: k[:],
			MessageNumber: uint64(msgNum),
		},
	}

	resp, err := s.api.request(&req, time.Second * 8)
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



}

func (s *DoubleRatchetKeyStoreApi) DeleteMk(k dr.Key, msgNum uint) {



}

func (s *DoubleRatchetKeyStoreApi) DeletePk(k dr.Key) {


}

func (s *DoubleRatchetKeyStoreApi) Count(k dr.Key) uint {



}

func (s *DoubleRatchetKeyStoreApi) All() map[dr.Key]map[uint]dr.Key {



}

func NewDRKeyStoreApi(api *API, km *keyManager.KeyManager) *DoubleRatchetKeyStoreApi {
	return &DoubleRatchetKeyStoreApi{
		api: api,
		km:  km,
	}
}

package api

import (
	"fmt"
	"sync"
	"testing"
	"time"

	pb "github.com/Bit-Nation/panthalassa/api/pb"
	"github.com/Bit-Nation/panthalassa/crypto/aes"
	keyManager "github.com/Bit-Nation/panthalassa/keyManager"
	keyStore "github.com/Bit-Nation/panthalassa/keyStore"
	mnemonic "github.com/Bit-Nation/panthalassa/mnemonic"
	proto "github.com/golang/protobuf/proto"
	require "github.com/stretchr/testify/require"
	dr "github.com/tiabc/doubleratchet"
)

func mustBeEqual(expected interface{}, got interface{}) {

	if expected != got {
		panic(fmt.Sprintf("Expected: %v to equal: %v", expected, got))
	}

}

func requireNil(value interface{}) {
	if value != nil {
		panic(fmt.Sprintf("Expected value: %s to be nil", value))
	}
}

func keyManagerFactory() *keyManager.KeyManager {

	mne, err := mnemonic.FromString("panda eyebrow bullet gorilla call smoke muffin taste mesh discover soft ostrich alcohol speed nation flash devote level hobby quick inner drive ghost inside")
	if err != nil {
		panic(err)
	}

	ks, err := keyStore.NewFromMnemonic(mne)
	if err != nil {
		panic(err)
	}

	return keyManager.CreateFromKeyStore(ks)

}

func TestDoubleRatchetKeyStore_GetSuccess(t *testing.T) {

	km := keyManagerFactory()

	k := dr.Key{}
	k[4] = 0x16

	rawMsgKey := dr.Key{}
	rawMsgKey[2] = 0x34

	ct, err := km.AESEncrypt(rawMsgKey[:])
	require.Nil(t, err)
	rawCt, err := ct.Marshal()
	require.Nil(t, err)

	c := make(chan string)

	api := API{
		lock:     sync.Mutex{},
		requests: map[string]chan *Response{},
		client: &testUpStream{
			sendFn: func(data string) {
				c <- data
			},
		},
	}

	go func() {

		select {
		case data := <-c:

			req := pb.Request{}
			requireNil(proto.Unmarshal([]byte(data), &req))

			keyRequest := req.DRKeyStoreGet
			key := dr.Key{}
			copy(key[:], keyRequest.DrKey)

			if k != key {
				panic("expected keys to match")
			}

			mustBeEqual(uint64(5), keyRequest.MessageNumber)

			err = api.Respond(req.RequestID, &pb.Response{
				DRKeyStoreGet: &pb.Response_DRKeyStoreGet{
					MessageKey: rawCt,
				},
			}, nil, time.Second)
			requireNil(err)

		}

	}()

	dra := DoubleRatchetKeyStoreApi{
		api: &api,
		km:  keyManagerFactory(),
	}

	key, exist := dra.Get(k, 5)
	require.Equal(t, rawMsgKey, key)
	require.True(t, exist)

}

func TestDoubleRatchetKeyStore_PutSuccess(t *testing.T) {

	indexKey := dr.Key{}
	indexKey[4] = 0x16

	msgKey := dr.Key{}
	msgKey[2] = 0x34

	c := make(chan string)

	api := API{
		lock:     sync.Mutex{},
		requests: map[string]chan *Response{},
		client: &testUpStream{
			sendFn: func(data string) {
				c <- data
			},
		},
	}

	go func() {

		select {
		case data := <-c:

			req := pb.Request{}
			requireNil(proto.Unmarshal([]byte(data), &req))

			keyRequest := req.DRKeyStorePut
			decodedIndexKey := dr.Key{}
			copy(decodedIndexKey[:], keyRequest.Key)

			// make sure index key is correct
			if decodedIndexKey != indexKey {
				panic("expected keys to match")
			}

			// check that message number is correct
			mustBeEqual(uint64(5), keyRequest.MessageNumber)

			km := keyManagerFactory()
			ct, err := aes.Unmarshal(keyRequest.MessageKey)
			requireNil(err)
			rawDecodedMsgKey, err := km.AESDecrypt(ct)
			decodedMsgKey := dr.Key{}
			copy(decodedMsgKey[:], rawDecodedMsgKey)

			// make sure message key correct
			if decodedMsgKey != msgKey {
				panic("expected keys to match")
			}

			err = api.Respond(req.RequestID, &pb.Response{}, nil, time.Second)
			requireNil(err)

		}

	}()

	dra := DoubleRatchetKeyStoreApi{
		api: &api,
		km:  keyManagerFactory(),
	}

	dra.Put(indexKey, 5, msgKey)

}

func TestDoubleRatchetKeyStore_DeleteMessageKeySuccess(t *testing.T) {

	c := make(chan string)

	msgKey := dr.Key{}
	msgKey[2] = 0x34

	api := API{
		lock:     sync.Mutex{},
		requests: map[string]chan *Response{},
		client: &testUpStream{
			sendFn: func(data string) {
				c <- data
			},
		},
	}

	go func() {

		select {
		case data := <-c:

			req := pb.Request{}
			requireNil(proto.Unmarshal([]byte(data), &req))
			keyRequest := req.DRKeyStoreDeleteMK

			decodedIndexKey := dr.Key{}
			copy(decodedIndexKey[:], keyRequest.Key)

			// make sure index key is correct
			if decodedIndexKey != msgKey {
				panic("expected keys to match")
			}

			// check that message number is correct
			mustBeEqual(uint64(6), keyRequest.MsgNum)

			api.Respond(req.RequestID, &pb.Response{}, nil, time.Second*2)
		}

	}()

	dra := DoubleRatchetKeyStoreApi{
		api: &api,
		km:  keyManagerFactory(),
	}

	dra.DeleteMk(msgKey, 6)

}

func TestDoubleRatchetKeyStore_DeleteIndexKeySuccess(t *testing.T) {

	c := make(chan string)

	msgKey := dr.Key{}
	msgKey[2] = 0x34

	api := API{
		lock:     sync.Mutex{},
		requests: map[string]chan *Response{},
		client: &testUpStream{
			sendFn: func(data string) {
				c <- data
			},
		},
	}

	go func() {

		select {
		case data := <-c:

			req := pb.Request{}
			requireNil(proto.Unmarshal([]byte(data), &req))
			keyRequest := req.DRKeyStoreDeleteKeys

			decodedIndexKey := dr.Key{}
			copy(decodedIndexKey[:], keyRequest.Key)

			// make sure index key is correct
			if decodedIndexKey != msgKey {
				panic("expected keys to match")
			}

			api.Respond(req.RequestID, &pb.Response{}, nil, time.Second*2)
		}

	}()

	dra := DoubleRatchetKeyStoreApi{
		api: &api,
		km:  keyManagerFactory(),
	}

	dra.DeletePk(msgKey)

}

func TestDoubleRatchetKeyStore_CountSuccess(t *testing.T) {

	c := make(chan string)

	msgKey := dr.Key{}
	msgKey[2] = 0x34

	api := API{
		lock:     sync.Mutex{},
		requests: map[string]chan *Response{},
		client: &testUpStream{
			sendFn: func(data string) {
				c <- data
			},
		},
	}

	go func() {

		select {
		case data := <-c:

			req := pb.Request{}
			requireNil(proto.Unmarshal([]byte(data), &req))
			keyRequest := req.DRKeyStoreCount

			decodedIndexKey := dr.Key{}
			copy(decodedIndexKey[:], keyRequest.Key)

			// make sure index key is correct
			if decodedIndexKey != msgKey {
				panic("expected keys to match")
			}

			api.Respond(req.RequestID, &pb.Response{
				DRKeyStoreCount: &pb.Response_DRKeyStoreCount{
					Count: uint64(4),
				},
			}, nil, time.Second*2)
		}

	}()

	dra := DoubleRatchetKeyStoreApi{
		api: &api,
		km:  keyManagerFactory(),
	}

	require.Equal(t, uint(4), dra.Count(msgKey))

}

func TestDoubleRatchetKeyStore_FetchAllKeysSuccess(t *testing.T) {

	c := make(chan string)

	indexKey := dr.Key{}
	indexKey[5] = 0x44

	msgKey := dr.Key{}
	msgKey[2] = 0x34

	encryptedMsgKey, err := keyManagerFactory().AESEncrypt(msgKey[:])
	require.Nil(t, err)
	rawEncryptedMsgKey, err := encryptedMsgKey.Marshal()
	require.Nil(t, err)

	api := API{
		lock:     sync.Mutex{},
		requests: map[string]chan *Response{},
		client: &testUpStream{
			sendFn: func(data string) {
				c <- data
			},
		},
	}

	go func() {

		select {
		case data := <-c:

			req := pb.Request{}
			requireNil(proto.Unmarshal([]byte(data), &req))

			api.Respond(req.RequestID, &pb.Response{
				DRKeyStoreAll: &pb.Response_DRKeyStoreAll{
					All: []*pb.Response_DRKeyStoreAll_Key{
						&pb.Response_DRKeyStoreAll_Key{
							Key: indexKey[:],
							MessageKeys: map[uint64][]byte{
								uint64(4): rawEncryptedMsgKey,
							},
						},
					},
				},
			}, nil, time.Second*2)
		}

	}()

	dra := DoubleRatchetKeyStoreApi{
		api: &api,
		km:  keyManagerFactory(),
	}

	require.Equal(t, msgKey, dra.All()[indexKey][4])

}

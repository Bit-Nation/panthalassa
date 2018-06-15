package client

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"testing"

	deviceApi "github.com/Bit-Nation/panthalassa/api/device"
	keyManager "github.com/Bit-Nation/panthalassa/keyManager"
	keyStore "github.com/Bit-Nation/panthalassa/keyStore"
	mnemonic "github.com/Bit-Nation/panthalassa/mnemonic"
	require "github.com/stretchr/testify/require"
	dr "github.com/tiabc/doubleratchet"
)

type UpStreamTestImpl struct {
	f func(data string)
}

func (u *UpStreamTestImpl) Send(data string) {
	u.f(data)
}

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

	c := make(chan string)

	// fake client
	api := deviceApi.New(&UpStreamTestImpl{
		f: func(data string) {
			c <- data
		},
	})

	// answer request of fake client
	go func() {
		for {
			select {
			case data := <-c:

				rpcCall := deviceApi.ApiCall{}
				if err := json.Unmarshal([]byte(data), &rpcCall); err != nil {
					panic(err)
				}

				mustBeEqual(`{"key":"0000000000000000000100010000000000010001000000000000000100000001","msg_num":7}`, rpcCall.Data)

				requireNil(api.Receive(rpcCall.Id, `{"error":"","payload":"{\"key\":\"{\\\"iv\\\":\\\"3Zd2O1KxUz2OZnWQPrTgCg==\\\",\\\"cipher_text\\\":\\\"q4FO26h5TICATqwwp9RXXXes1jX8asn+0TkL5Khx8Oc=\\\",\\\"mac\\\":\\\"XwX884HeXuodY3vgoKvmZcGkW0oPu2fBvRafxAsMu/I=\\\",\\\"v\\\":2}\"}"}`))

			}
		}
	}()

	drk := DoubleRatchetKeyStore{
		api: api,
		km:  keyManagerFactory(),
	}

	k := dr.Key{
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x01, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x01, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01,
	}

	key, exist := drk.Get(k, 7)
	require.True(t, exist)
	require.Equal(t, dr.Key{
		0x10, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}, key)

}

func TestDoubleRatchetKeyStore_GetError(t *testing.T) {

	c := make(chan string)

	// fake client
	api := deviceApi.New(&UpStreamTestImpl{
		f: func(data string) {
			c <- data
		},
	})

	// answer request of fake client
	go func() {
		for {
			select {
			case data := <-c:

				rpcCall := deviceApi.ApiCall{}
				if err := json.Unmarshal([]byte(data), &rpcCall); err != nil {
					panic(err)
				}

				mustBeEqual(`{"key":"0000000000000000000100010000000000010001000000000000000100000001","msg_num":7}`, rpcCall.Data)

				requireNil(api.Receive(rpcCall.Id, `{"error":"I am an error","payload":""}`))

			}
		}
	}()

	drk := DoubleRatchetKeyStore{
		api: api,
		km:  keyManagerFactory(),
	}

	k := dr.Key{
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x01, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x01, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01,
	}

	key, exist := drk.Get(k, 7)
	require.False(t, exist)
	require.Equal(t, dr.Key{
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}, key)

}

func TestDoubleRatchetKeyStore_PutSuccess(t *testing.T) {

	c := make(chan string)

	km := keyManagerFactory()

	// fake client
	api := deviceApi.New(&UpStreamTestImpl{
		f: func(data string) {
			c <- data
		},
	})

	pubKey := dr.Key{
		0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}

	msgKey := dr.Key{
		0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01,
		0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01,
		0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01,
		0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01,
	}

	// answer request of fake client
	go func() {
		for {
			select {
			case data := <-c:

				rpcCall := deviceApi.ApiCall{}
				if err := json.Unmarshal([]byte(data), &rpcCall); err != nil {
					panic(err)
				}

				// check type
				mustBeEqual("DR:KEY_STORE:PUT", rpcCall.Type)

				var c DRKeyStorePutCall
				requireNil(json.Unmarshal([]byte(rpcCall.Data), &c))

				mustBeEqual("0100000000000000000000000000000000000000000000000000000000000000", c.IndexKey)
				mustBeEqual(uint(30), c.MsgNumber)

				plainKey, err := km.AESDecrypt(c.DoubleRatchetKey)
				requireNil(err)

				mustBeEqual(hex.EncodeToString(msgKey[:]), hex.EncodeToString(plainKey))

				requireNil(api.Receive(rpcCall.Id, `{"error":"","payload":""}`))

			}
		}
	}()

	drk := DoubleRatchetKeyStore{
		api: api,
		km:  keyManagerFactory(),
	}

	drk.Put(pubKey, 30, msgKey)

}

func TestDoubleRatchetKeyStore_DeleteMessageKeySuccess(t *testing.T) {

	c := make(chan string)

	// fake client
	api := deviceApi.New(&UpStreamTestImpl{
		f: func(data string) {
			c <- data
		},
	})

	// answer request of fake client
	go func() {
		for {
			select {
			case data := <-c:

				rpcCall := deviceApi.ApiCall{}
				if err := json.Unmarshal([]byte(data), &rpcCall); err != nil {
					panic(err)
				}
				mustBeEqual("DR:KEY_STORE:DELETE_MESSAGE_KEY", rpcCall.Type)

				mustBeEqual(`{"index_key":"0000000000000000000100010000000000010001000000000000000100000001","msg_num":3032}`, rpcCall.Data)

				requireNil(api.Receive(rpcCall.Id, `{"error":"","payload":""}`))

			}
		}
	}()

	drk := DoubleRatchetKeyStore{
		api: api,
		km:  keyManagerFactory(),
	}

	k := dr.Key{
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x01, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x01, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01,
	}

	drk.DeleteMk(k, 3032)

}

func TestDoubleRatchetKeyStore_DeleteIndexKeySuccess(t *testing.T) {

	c := make(chan string)

	// fake client
	api := deviceApi.New(&UpStreamTestImpl{
		f: func(data string) {
			c <- data
		},
	})

	// answer request of fake client
	go func() {
		for {
			select {
			case data := <-c:

				rpcCall := deviceApi.ApiCall{}
				if err := json.Unmarshal([]byte(data), &rpcCall); err != nil {
					panic(err)
				}
				mustBeEqual("DR:KEY_STORE:DELETE_INDEX_KEY", rpcCall.Type)

				mustBeEqual(`{"index_key":"0000000000000000000100010000000000010001000000000000000100000001"}`, rpcCall.Data)

				requireNil(api.Receive(rpcCall.Id, `{"error":"","payload":""}`))

			}
		}
	}()

	drk := DoubleRatchetKeyStore{
		api: api,
		km:  keyManagerFactory(),
	}

	k := dr.Key{
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x01, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x01, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01,
	}

	drk.DeletePk(k)

}

func TestDoubleRatchetKeyStore_CountSuccess(t *testing.T) {

	c := make(chan string)

	// fake client
	api := deviceApi.New(&UpStreamTestImpl{
		f: func(data string) {
			c <- data
		},
	})

	// answer request of fake client
	go func() {
		for {
			select {
			case data := <-c:

				rpcCall := deviceApi.ApiCall{}
				if err := json.Unmarshal([]byte(data), &rpcCall); err != nil {
					panic(err)
				}

				mustBeEqual("DR:KEY_STORE:COUNT_MESSAGES", rpcCall.Type)

				mustBeEqual(`{"index_key":"0000000000000000000100010000000000010001000000000000000100000001"}`, rpcCall.Data)

				requireNil(api.Receive(rpcCall.Id, `{"error":"","payload":"{\"key\":\"{\\\"count\\\":3}\"}"}`))

			}
		}
	}()

	drk := DoubleRatchetKeyStore{
		api: api,
		km:  keyManagerFactory(),
	}

	k := dr.Key{
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x01, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x01, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01,
	}

	drk.Count(k)

}

func TestDoubleRatchetKeyStore_FetchAllKeysSuccess(t *testing.T) {

	c := make(chan string)

	// fake client
	api := deviceApi.New(&UpStreamTestImpl{
		f: func(data string) {
			c <- data
		},
	})

	// answer request of fake client
	go func() {
		for {
			select {
			case data := <-c:

				rpcCall := deviceApi.ApiCall{}
				if err := json.Unmarshal([]byte(data), &rpcCall); err != nil {
					panic(err)
				}

				mustBeEqual("DR:KEY_STORE:FETCH_ALL_KEYS", rpcCall.Type)
				mustBeEqual(``, rpcCall.Data)

				requireNil(api.Receive(rpcCall.Id, `{"error":"","payload":"{\"0000000000000000000100010000000000010001000000000000000100000001\":{\"444\":\"{\\\"iv\\\":\\\"6IW62YJEo5o+Yc9eeRoyaA==\\\",\\\"cipher_text\\\":\\\"61eWJQS97l2t8ncDpSKKy567kidCKu9Po2cA/ZJZxUQ=\\\",\\\"mac\\\":\\\"v9n9PM68Wdj+g86hQUIS6ZyweDOtjhUyU34+xIos1q8=\\\",\\\"v\\\":1}\"}}"}`))

			}
		}
	}()

	drk := DoubleRatchetKeyStore{
		api: api,
		km:  keyManagerFactory(),
	}

	fmt.Println(drk.All())

}

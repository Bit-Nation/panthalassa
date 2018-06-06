package client

import (
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
		panic(fmt.Sprintf("Expected: %s to equal: %s", expected, got))
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

				requireNil(api.Receive(rpcCall.Id, `{"error":"","payload":"{\"key\":\"{\\\"iv\\\":\\\"f+ZgKGwhkz82bokcs7HI8A==\\\",\\\"cipher_text\\\":\\\"Qu0SeXs/ahOjqcwPRIrK0sR9ngirapvt33x3SNLayFs=\\\",\\\"mac\\\":\\\"vJhCluH/bWdkcaA3vtTzDDsYFVO0A7UcL7wbPbvYdG0=\\\",\\\"v\\\":1}\"}"}`))

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

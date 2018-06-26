package api

import (
	"testing"
	"time"

	pb "github.com/Bit-Nation/panthalassa/api/pb"
	dapp "github.com/Bit-Nation/panthalassa/dapp"
	proto "github.com/golang/protobuf/proto"
	require "github.com/stretchr/testify/require"
)

func TestAPI_ShowModal(t *testing.T) {

	c := make(chan string)

	api := New(&UpStreamTestImpl{
		f: func(data string) {
			c <- data
		},
	}, keyManagerFactory())

	go func() {

		select {
		case data := <-c:

			req := pb.Request{}
			requireNil(proto.Unmarshal([]byte(data), &req))

			if req.ShowModal.Title != "Request Money" {
				panic("Expected title to be 'Request Money'")
			}

			if req.ShowModal.Layout != "{}" {
				panic("Expected layout to be '{}'")
			}

			err := api.Respond(req.RequestID, &pb.Response{}, nil, time.Second*3)
			if err != nil {
				panic("expected error to be nil")
			}

		}

	}()

	err := api.ShowModal("Request Money", "{}")
	require.Nil(t, err)

}

func TestAPI_SendEthereumTransaction(t *testing.T) {

	c := make(chan string)

	api := New(&UpStreamTestImpl{
		f: func(data string) {
			c <- data
		},
	}, keyManagerFactory())

	go func() {

		select {
		case data := <-c:

			req := pb.Request{}
			requireNil(proto.Unmarshal([]byte(data), &req))

			tx := req.SendEthereumTransaction

			resp := &pb.Response_SendEthereumTransaction{
				Nonce:    3,
				GasPrice: "1000000000",
				GasLimit: "100000000000",
				To:       tx.To,
				Value:    tx.Value,
				Data:     tx.Data,
				V:        "v_of_tx",
				R:        "r_of_tx",
				S:        "s_of_tx",
				ChainID:  4,
				From:     "my_address",
				Hash:     "tx-hash",
			}

			err := api.Respond(req.RequestID, &pb.Response{
				SendEthereumTransaction: resp,
			}, nil, time.Second*5)
			if err != nil {
				panic(err)
			}
		}

	}()

	resp, err := api.SendEthereumTransaction("100", "0x1f75bb626ad018f3354259b10ab2e74bc0e0f267", "0xf3")
	require.Nil(t, err)
	require.Equal(t, `{"chainId":4,"data":"0xf3","from":"my_address","gasLimit":"100000000000","gasPrice":"1000000000","hash":"tx-hash","nonce":3,"r":"r_of_tx","s":"s_of_tx","to":"0x1f75bb626ad018f3354259b10ab2e74bc0e0f267","v":"v_of_tx","value":"100"}`, resp)
}

func TestAPI_SaveDApp(t *testing.T) {

	c := make(chan string)

	api := New(&UpStreamTestImpl{
		f: func(data string) {
			c <- data
		},
	}, keyManagerFactory())

	go func() {

		select {
		case data := <-c:

			req := pb.Request{}
			requireNil(proto.Unmarshal([]byte(data), &req))

			app := req.SaveDApp

			if "My DApp" != app.AppName {
				panic("Expected dapp name to be: My DApp")
			}
			
			if "var i = 3" != app.Code {
				panic("Expected dapp name to be: var i = 3")
			}
			
			if "7075626b6579" != app.SigningPublicKey {
				panic("Expected dapp name to be: My DApp")
			}
			
			if "7369676e6174757265" != app.Signature {
				panic("Expected dapp name to be: My DApp")
			}

			err := api.Respond(req.RequestID, &pb.Response{}, nil, time.Second*5)
			if err != nil {
				panic(err)
			}
		}

	}()

	err := api.SaveDApp(dapp.JsonRepresentation{
		Name:               "My DApp",
		Code:               "var i = 3",
		SignaturePublicKey: []byte("pubkey"),
		Signature:          []byte("signature"),
	})
	require.Nil(t, err)
}

package api

import (
	"time"

	"encoding/json"
	pb "github.com/Bit-Nation/panthalassa/api/pb"
	"github.com/kataras/iris/core/errors"
)

func (a *API) ShowModal(title, layout string) error {
	return a.dAppApi.ShowModal(title, layout)
}

func (a *API) SendEthereumTransaction(value, to, data string) (string, error) {
	return a.dAppApi.SendEthereumTransaction(value, to, data)
}

type DAppApi struct {
	api *API
}

// request to show a modal
func (a *DAppApi) ShowModal(title, layout string) error {

	// send request
	_, err := a.api.request(&pb.Request{
		ShowModal: &pb.Request_ShowModal{
			Title:  title,
			Layout: layout,
		},
	}, time.Second*20)

	return err

}

// send an ethereum transaction to api
func (a *DAppApi) SendEthereumTransaction(value, to, data string) (string, error) {

	// send request
	resp, err := a.api.request(&pb.Request{
		SendEthereumTransaction: &pb.Request_SendEthereumTransaction{
			Value: value,
			To:    to,
			Data:  data,
		},
	}, time.Second*120)
	if err != nil {
		return "", err
	}

	ethTx := resp.Msg.SendEthereumTransaction
	if ethTx == nil {
		resp.Closer <- errors.New("got nil response")
	}

	objTx := map[string]interface{}{
		"nonce":    ethTx.Nonce,
		"gasPrice": ethTx.GasPrice,
		"gasLimit": ethTx.GasLimit,
		"to":       ethTx.To,
		"value":    ethTx.Value,
		"data":     ethTx.Data,
		"v":        ethTx.V,
		"r":        ethTx.R,
		"s":        ethTx.S,
		"chainId":  ethTx.ChainID,
		"from":     ethTx.From,
		"hash":     ethTx.Hash,
	}

	raw, err := json.Marshal(objTx)
	if err != nil {
		resp.Closer <- err
		return "", err
	}

	return string(raw), nil

}

package api

import (
	"encoding/json"
	"errors"
	"time"

	pb "github.com/Bit-Nation/panthalassa/api/pb"
	ed25519 "golang.org/x/crypto/ed25519"
)

func (a *API) RenderModal(uiID, layout string, dAppIDKey ed25519.PublicKey) error {
	return a.dAppApi.RenderModal(uiID, layout, dAppIDKey)
}

func (a *API) SendEthereumTransaction(value, to, data string) (string, error) {
	return a.dAppApi.SendEthereumTransaction(value, to, data)
}

type DAppApi struct {
	api *API
}

// request to show a modal
func (a *DAppApi) RenderModal(uiID, layout string, dAppPubKey ed25519.PublicKey) error {

	// send request
	resp, err := a.api.request(&pb.Request{
		ShowModal: &pb.Request_RenderModal{
			DAppPublicKey: dAppPubKey,
			UiID:          uiID,
			Layout:        layout,
		},
	}, time.Second*20)
	if err != nil {
		return err
	}

	// close since we don't care about the response
	resp.Closer <- nil

	return nil
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
		resp.Closer <- errors.New("got nil ethTx response")
		return "", errors.New("got nil ethTx response")
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

	// Since closer is not passed further, we need to close it here to prevent timeout.
	resp.Closer <- nil
	return string(raw), nil

}

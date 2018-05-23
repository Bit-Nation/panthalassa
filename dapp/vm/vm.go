package vm

import (
	deviceApi "github.com/Bit-Nation/panthalassa/api/device"
	dapp "github.com/Bit-Nation/panthalassa/dapp"
	extension "github.com/Bit-Nation/panthalassa/dapp/vm/extension"
	otto "github.com/robertkrimen/otto"
)

type DAppVm struct {
	otto      *otto.Otto
	closeChan chan<- struct{}
}

func NewVM(api *deviceApi.Api, dApp *dapp.DApp) (*DAppVm, error) {

	dAppVm := &DAppVm{}

	vm := otto.New()
	dAppVm.otto = vm

	// register util to sign transaction's
	vm.Set("signTransaction", extension.SignTransaction(api))

	// run the code
	// @todo this must be an "async" process. 3rd party code could freeze the whole process.
	_, err := vm.Run(dApp.Code)
	if err != nil {
		return nil, err
	}

	return dAppVm, nil

}

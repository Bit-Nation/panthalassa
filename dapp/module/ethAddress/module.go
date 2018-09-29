package ethAddress

import (
	keyManager "github.com/Bit-Nation/panthalassa/keyManager"
	log "github.com/ipfs/go-log"
	otto "github.com/robertkrimen/otto"
)

var logger = log.Logger("eth address")

func New(km *keyManager.KeyManager) *Module {
	return &Module{
		km: km,
	}
}

type Module struct {
	km *keyManager.KeyManager
}

func (m *Module) Close() error {
	return nil
}

func (m *Module) Register(vm *otto.Otto) error {

	logger.Debug("get ethereum address")

	addr, err := m.km.GetEthereumAddress()
	if err != nil {
		return err
	}

	return vm.Set("ethereumAddress", addr)

}

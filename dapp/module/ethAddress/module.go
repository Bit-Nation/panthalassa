package ethAddress

import (
	keyManager "github.com/Bit-Nation/panthalassa/keyManager"
	otto "github.com/robertkrimen/otto"
)

func New(km *keyManager.KeyManager) *Module {
	return &Module{
		km: km,
	}
}

type Module struct {
	km *keyManager.KeyManager
}

func (m *Module) Name() string {
	return "ETH:ADDRESS"
}

func (m *Module) Register(vm *otto.Otto) error {

	addr, err := m.km.GetEthereumAddress()
	if err != nil {
		return err
	}

	return vm.Set("ethereumAddress", addr)

}

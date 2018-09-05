package ethAddress

import (
	keyManager "github.com/Bit-Nation/panthalassa/keyManager"
	log "github.com/ipfs/go-log"
	duktape "gopkg.in/olebedev/go-duktape.v3"
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

func (m *Module) Register(context *duktape.Context) error {

	logger.Debug("get ethereum address")

	addr, err := m.km.GetEthereumAddress()
	if err != nil {
		return err
	}
	context.PushString(addr)
	return err

}

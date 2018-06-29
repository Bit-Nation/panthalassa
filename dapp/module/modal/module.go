package modal

import (
	validator "github.com/Bit-Nation/panthalassa/dapp/validator"
	log "github.com/op/go-logging"
	otto "github.com/robertkrimen/otto"
	ed25519 "golang.org/x/crypto/ed25519"
)

type Device interface {
	ShowModal(title, layout string, dAppIDKey ed25519.PublicKey) error
}

type Module struct {
	device    Device
	logger    *log.Logger
	dAppIDKey ed25519.PublicKey
}

// create new Modal Module
func New(l *log.Logger, device Device, dAppIDKey ed25519.PublicKey) *Module {
	return &Module{
		logger:    l,
		device:    device,
		dAppIDKey: dAppIDKey,
	}
}

func (m *Module) Name() string {
	return "MODAL"
}

// showModal provides a way to display a modal
// the first parameter should be the modal title
// the second parameter should be the layout to render
// and the third parameter is an optional callback that
// will called with an optional error
func (m *Module) Register(vm *otto.Otto) error {

	return vm.Set("showModal", func(call otto.FunctionCall) otto.Value {

		// validate function call
		v := validator.New()
		v.Set(0, &validator.TypeString)
		v.Set(1, &validator.TypeString)
		v.Set(2, &validator.TypeFunction)
		if err := v.Validate(vm, call); err != nil {
			m.logger.Error(err.String())
			return *err
		}

		// call callback
		cb := call.Argument(2)

		// do the request to show modal async
		go func() {

			// request to show modal
			err := m.device.ShowModal(
				call.Argument(0).String(),
				call.Argument(1).String(),
				m.dAppIDKey,
			)

			if err != nil {
				cb.Call(cb, err.Error())
				return
			}

			cb.Call(cb)

		}()

		return otto.Value{}

	})

}

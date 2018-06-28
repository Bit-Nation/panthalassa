package modal

import (
	validator "github.com/Bit-Nation/panthalassa/dapp/validator"
	log "github.com/op/go-logging"
	otto "github.com/robertkrimen/otto"
)

type Device interface {
	ShowModal(title, layout string) error
}

type Module struct {
	device Device
	logger *log.Logger
}

// create new Modal Module
func New(l *log.Logger, device Device) *Module {
	return &Module{
		logger: l,
		device: device,
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

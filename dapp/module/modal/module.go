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

func (m *Module) Register(vm *otto.Otto) error {

	return vm.Set("showModal", func(call otto.FunctionCall) otto.Value {

		// validate function call
		v := validator.New()
		v.Set(0, &validator.TypeString)
		v.Set(1, &validator.TypeString)
		if err := v.Validate(vm, call); err != nil {
			return *err
		}

		// request to show modal
		err := m.device.ShowModal(
			call.Argument(0).String(),
			call.Argument(1).String(),
		)

		// call callback if passed in
		fn := call.Argument(2)
		if fn.IsFunction() {
			if err != nil {
				fn.Call(fn, err.Error())
				return otto.Value{}
			}
			fn.Call(fn)
		}

		return otto.Value{}

	})

}

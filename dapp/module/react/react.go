package react

import (
	logger "github.com/op/go-logging"
	otto "github.com/robertkrimen/otto"
)

// the client who displays
// the user interface
type Client interface {
	Render(jsonJSX string) error
}

type React struct {
	Client Client
	Logger *logger.Logger
}

// create new react module
func New(c Client, l *logger.Logger) *React {
	return &React{
		Client: c,
		Logger: l,
	}
}

func (r *React) Name() string {
	return "REACT"
}

func (r *React) Register(vm *otto.Otto) error {

	err := vm.Set("renderReact", func(call otto.FunctionCall) otto.Value {

		if !call.Argument(0).IsString() {
			v, err := otto.ToValue("expected element to be a string")
			if err != nil {
				r.Logger.Error(err.Error())
				return otto.Value{}
			}
			return v
		}

		err := r.Client.Render(call.Argument(0).String())
		if err != nil {
			v, err := otto.ToValue(err.Error())
			if err != nil {
				r.Logger.Error(err.Error())
				return otto.Value{}
			}
			return v
		}

		return otto.Value{}
	})
	if err != nil {
		return err
	}

	return nil
}

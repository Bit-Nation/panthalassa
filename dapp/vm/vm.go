package vm

import (
	"github.com/robertkrimen/otto"
)

func NewVM() *otto.Otto {

	vm := otto.New()

	vm.Set("signTransaction", func(call otto.FunctionCall) {

	})

	return vm

}

package module

import (
	otto "github.com/robertkrimen/otto"
)

type Module interface {
	Name() string
	Register(vm *otto.Otto) error
}
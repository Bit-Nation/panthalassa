package module

import (
	otto "github.com/robertkrimen/otto"
)

type Module interface {
	Register(vm *otto.Otto) error
	Close() error
}

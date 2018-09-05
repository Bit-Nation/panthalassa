package module

import (
	duktape "gopkg.in/olebedev/go-duktape.v3"
)

type Module interface {
	Register(vm *duktape.Context) error
	Close() error
}

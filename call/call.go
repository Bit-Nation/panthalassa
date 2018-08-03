package call

import (
	"fmt"
	"sync"
)

type Handler interface {
	Handle(params map[string]interface{}) (interface{}, error)
	Name() string
}

type Call struct {
	lock     sync.Mutex
	handlers map[string]Handler
}

func (c *Call) AddHandler(h Handler) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.handlers[h.Name()] = h
}

func (c *Call) Call(name string, params map[string]interface{}) (interface{}, error) {

	// fetch handler
	c.lock.Lock()
	h, exist := c.handlers[name]
	c.lock.Unlock()

	if !exist {
		return nil, fmt.Errorf("can't handle %s - handler has not been registered", name)
	}

	return h.Handle(params)

}

func NewCall() *Call {
	return &Call{
		lock: sync.Mutex{},
	}
}

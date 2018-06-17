package validator

import (
	"errors"
	"fmt"
	"sync"

	otto "github.com/robertkrimen/otto"
)

const (
	TypeFunction = iota
	TypeNumber   = iota
	TypeObject   = iota
)

var validators = map[int]func(call otto.FunctionCall, position int) error{
	TypeFunction: func(call otto.FunctionCall, position int) error {
		if !call.Argument(position).IsFunction() {
			return errors.New(fmt.Sprintf("expected parameter %d to be of type function", position))
		}
		return nil
	},
	TypeNumber: func(call otto.FunctionCall, position int) error {
		if !call.Argument(position).IsNumber() {
			return errors.New(fmt.Sprintf("expected parameter %d to be of type number", position))
		}
		return nil
	},
	TypeObject: func(call otto.FunctionCall, position int) error {
		if !call.Argument(position).IsObject() {
			return errors.New(fmt.Sprintf("expected parameter %d to be of type object", position))
		}
		return nil
	},
}

type CallValidator struct {
	lock  sync.Mutex
	rules map[int]int
}

// add validation rule
func (v *CallValidator) Set(index int, expectedType int) error {
	v.lock.Lock()
	defer v.lock.Unlock()
	_, exist := validators[expectedType]
	if !exist {
		return errors.New("type does not exist")
	}
	v.rules[index] = expectedType
	return nil
}

func (v *CallValidator) Validate(call otto.FunctionCall) error {

	v.lock.Lock()
	defer v.lock.Unlock()
	for index, expectedType := range v.rules {
		validator, exist := validators[expectedType]
		if !exist {
			return errors.New(fmt.Sprintf("couldn't find validator for type: %d", index))
		}
		if err := validator(call, index); err != nil {
			return err
		}
	}
	return nil

}

func New() *CallValidator {
	return &CallValidator{
		lock:  sync.Mutex{},
		rules: map[int]int{},
	}
}

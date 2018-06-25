package validator

import (
	"errors"
	"fmt"
	"sort"
	"sync"

	otto "github.com/robertkrimen/otto"
)

type Validator = func(call otto.FunctionCall, position int) error

var TypeFunction = func(call otto.FunctionCall, position int) error {
	if !call.Argument(position).IsFunction() {
		return errors.New(fmt.Sprintf("expected parameter %d to be of type function", position))
	}
	return nil
}

var TypeNumber = func(call otto.FunctionCall, position int) error {
	if !call.Argument(position).IsNumber() {
		return errors.New(fmt.Sprintf("expected parameter %d to be of type number", position))
	}
	return nil
}

var TypeString = func(call otto.FunctionCall, position int) error {
	if !call.Argument(position).IsString() {
		return errors.New(fmt.Sprintf("expected parameter %d to be of type string", position))
	}
	return nil
}

type CallValidator struct {
	lock  sync.Mutex
	rules map[int]*Validator
}

// add validation rule
func (v *CallValidator) Set(index int, validator *Validator) error {
	v.lock.Lock()
	defer v.lock.Unlock()
	v.rules[index] = validator
	return nil
}

func (v *CallValidator) Validate(vm *otto.Otto, call otto.FunctionCall) *otto.Value {

	v.lock.Lock()
	defer v.lock.Unlock()

	keys := make([]int, 0)
	for k, _ := range v.rules {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	for k := range keys {

		validator, exist := v.rules[k]
		if !exist {
			ve := vm.MakeCustomError("ValidationError", fmt.Sprintf("couldn't find validator for type: %d", k))
			return &ve
		}

		v := *validator
		if err := v(call, k); err != nil {
			ve := vm.MakeCustomError("ValidationError", err.Error())
			return &ve
		}

	}

	return nil

}

func New() *CallValidator {
	return &CallValidator{
		lock:  sync.Mutex{},
		rules: map[int]*Validator{},
	}
}

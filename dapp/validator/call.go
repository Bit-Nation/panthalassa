package validator

import (
	"errors"
	"fmt"
	"sync"

	otto "github.com/robertkrimen/otto"
	"sort"
)

type VMType int

const (
	TypeFunction VMType = 10
	TypeNumber   VMType = 20
	TypeObject   VMType = 30
	TypeString   VMType = 40
)

var validators = map[int]func(call otto.FunctionCall, position int) error{
	int(TypeFunction): func(call otto.FunctionCall, position int) error {
		if !call.Argument(position).IsFunction() {
			return errors.New(fmt.Sprintf("expected parameter %d to be of type function", position))
		}
		return nil
	},
	int(TypeNumber): func(call otto.FunctionCall, position int) error {
		if !call.Argument(position).IsNumber() {
			return errors.New(fmt.Sprintf("expected parameter %d to be of type number", position))
		}
		return nil
	},
	int(TypeObject): func(call otto.FunctionCall, position int) error {
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
func (v *CallValidator) Set(index int, expectedType VMType) error {
	v.lock.Lock()
	defer v.lock.Unlock()
	_, exist := validators[int(expectedType)]
	if !exist {
		return errors.New("type does not exist")
	}
	v.rules[index] = int(expectedType)
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
		validator, exist := validators[k]
		if !exist {
			ve := vm.MakeCustomError("ValidationError", fmt.Sprintf("couldn't find validator for type: %d", index))
			return &ve
		}
		if err := validator(call, k); err != nil {
			ve := vm.MakeCustomError("ValidationError", err.Error())
			return &ve
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

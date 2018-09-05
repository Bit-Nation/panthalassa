package validator

import (
	"errors"
	"fmt"
	"sort"
	"sync"

	duktape "gopkg.in/olebedev/go-duktape.v3"
)

type Validator = func(context *duktape.Context, position int) error

var TypeFunction = func(context *duktape.Context, position int) error {
	if !context.IsFunction(position) {
		return errors.New(fmt.Sprintf("expected parameter %d to be of type function", position))
	}
	return nil
}

var TypeNumber = func(context *duktape.Context, position int) error {
	if !context.IsNumber(position) {
		return errors.New(fmt.Sprintf("expected parameter %d to be of type number", position))
	}
	return nil
}

var TypeString = func(context *duktape.Context, position int) error {
	if !context.IsString(position) {
		return errors.New(fmt.Sprintf("expected parameter %d to be of type string", position))
	}
	return nil
}

var TypeObject = func(context *duktape.Context, position int) error {
	if !context.IsObject(position) {
		return errors.New(fmt.Sprintf("expected parameter %d to be of type object", position))
	}
	return nil
}

type CallValidator struct {
	lock  sync.Mutex
	rules map[int]*Validator
}

// add validation rule
func (v *CallValidator) Set(index int, validator *Validator) {
	v.lock.Lock()
	defer v.lock.Unlock()
	v.rules[index] = validator
}

func (v *CallValidator) Validate(vm *duktape.Context) error {

	v.lock.Lock()
	defer v.lock.Unlock()

	keys := make([]int, 0)
	for k, _ := range v.rules {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	for _, k := range keys {

		validator, exist := v.rules[k]
		if !exist {
			ve := fmt.Errorf("ValidationError: couldn't find validator for type: %d", k)
			return ve
		}

		v := *validator
		if err := v(vm, k); err != nil {
			ve := fmt.Errorf("ValidationError: %s", err.Error())
			return ve
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

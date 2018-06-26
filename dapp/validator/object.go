package validator

import (
	"errors"
	"fmt"
	"sync"

	eth "github.com/ethereum/go-ethereum/common"
	otto "github.com/robertkrimen/otto"
)

type ObjValueValidator func(value otto.Value, objKey string) error

// validate an address
var ObjTypeAddress = func(value otto.Value, objKey string) error {
	err := errors.New(fmt.Sprintf("Expected %s to be an ethereum address", objKey))
	if !value.IsString() {
		return err
	}
	if !eth.IsHexAddress(value.String()) {
		return err
	}
	return nil
}

// Validate if value is a string
var ObjTypeString = func(value otto.Value, objKey string) error {
	if !value.IsString() {
		return errors.New(fmt.Sprintf("Expected %s to be an string", objKey))
	}
	return nil
}

type objValidatorEntry struct {
	valueType ObjValueValidator
	required  bool
}

type ObjValidator struct {
	rules map[string]objValidatorEntry
	lock  sync.Mutex
}

func NewObjValidator() *ObjValidator {
	return &ObjValidator{
		rules: map[string]objValidatorEntry{},
		lock:  sync.Mutex{},
	}
}

// add an validator item
func (v *ObjValidator) Set(key string, validator ObjValueValidator, required bool) {

	v.lock.Lock()
	defer v.lock.Unlock()

	v.rules[key] = objValidatorEntry{
		valueType: validator,
		required:  required,
	}

}

func (v *ObjValidator) Validate(vm *otto.Otto, obj otto.Object) *otto.Value {

	v.lock.Lock()
	defer v.lock.Unlock()

	for objKey, rule := range v.rules {

		keyExist := false
		for _, k := range obj.Keys() {
			if k == objKey {
				keyExist = true
			}
		}

		// exit when key is required but missing
		if !keyExist && rule.required {
			e := vm.MakeCustomError(
				"ValidationError",
				fmt.Sprintf("Missing value for required key: %s", objKey),
			)
			return &e
		}

		// in the case the ke doesn't exist we should just move on
		if !keyExist {
			continue
		}

		// get value for key
		value, err := obj.Get(objKey)
		if err != nil {
			e := vm.MakeCustomError(
				"InternalError",
				err.Error(),
			)
			return &e
		}

		// validate value
		if err := rule.valueType(value, objKey); err != nil {
			e := vm.MakeCustomError(
				"ValidationError",
				err.Error(),
			)
			return &e
		}

	}

	return nil

}

package validator

import (
	"errors"
	"fmt"
	"sync"

	eth "github.com/ethereum/go-ethereum/common"
	duktape "gopkg.in/olebedev/go-duktape.v3"
)

type ObjValueValidator func(context *duktape.Context, position int) error

// validate an address
var ObjTypeAddress = func(context *duktape.Context, position int) error {
	err := errors.New(fmt.Sprintf("Expected %s to be an ethereum address", context.ToString(position)))
	if !context.IsString(position) {
		return err
	}
	if !eth.IsHexAddress(context.ToString(position)) {
		return err
	}
	return nil
}

// Validate if value is a string
var ObjTypeString = func(context *duktape.Context, position int) error {
	if !context.IsString(position) {
		return errors.New(fmt.Sprintf("Expected %s to be a string", context.ToString(position)))
	}
	return nil
}

// validate if value is an object
var ObjTypeObject = func(context *duktape.Context, position int) error {
	if !context.IsObject(position) {
		return fmt.Errorf("expected %s to be a object", context.ToString(position))
	}
	return nil
}

var ObjTypeBool = func(context *duktape.Context, position int) error {
	if !context.IsBoolean(position) {
		return fmt.Errorf("expected %s to be a bool", context.ToString(position))
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

func (v *ObjValidator) Validate(context *duktape.Context, position int) error {

	v.lock.Lock()
	defer v.lock.Unlock()

	for objKey, rule := range v.rules {
		keyExist := false
		if context.HasPropString(position, objKey) {
			keyExist = true
		}
		// exit when key is required but missing
		if !keyExist && rule.required {
			e := fmt.Errorf("ValidationError: Missing value for required key: %s", objKey)
			return e
		}

		// in the case the ke doesn't exist we should just move on
		if !keyExist {
			continue
		}

		// get value for key
		if !context.GetPropString(position, objKey) {
			e := fmt.Errorf("InternalError: Key doesn't exist  %s", objKey)
			return e
		}
		// validate value
		// @TODO find out why -1 works here and "position" doesnt
		if err := rule.valueType(context, -1); err != nil {
			//e := fmt.Errorf("ValidationError : Missing value for required key: %s", context.ToString(position))
			return err
		}

	}

	return nil

}

package validator

import (
	"testing"

	otto "github.com/robertkrimen/otto"
	require "github.com/stretchr/testify/require"
)

func TestNewObjValidatorOptionalFields(t *testing.T) {

	v := NewObjValidator()
	v.Set("to", ObjTypeAddress, false)

	vm := otto.New()
	ob, err := vm.Object(`({})`)
	require.Nil(t, err)

	// expect no error since it's ok omit the address
	validationErr := v.Validate(vm, *ob)
	require.Nil(t, validationErr)

}

func TestNewObjValidationRequiredFields(t *testing.T) {

	/**
	 * Test required field - no error
	 */
	v := NewObjValidator()
	v.Set("to", ObjTypeAddress, true)

	vm := otto.New()
	ob, err := vm.Object(`({"to": "0x0b078896a3d9166da5c37ae52a5809aca48630d4"})`)
	require.Nil(t, err)

	// expect no error since we provide the address
	require.Nil(t, v.Validate(vm, *ob))

	/**
	 * Test required field - with error
	 */
	v = NewObjValidator()
	v.Set("to", ObjTypeAddress, false)

	vm = otto.New()
	ob, err = vm.Object(`({})`)
	require.Nil(t, err)

	// expect error since we don't pass the required address in
	require.Nil(t, v.Validate(vm, *ob))

}

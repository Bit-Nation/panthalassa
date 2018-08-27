package validator

import (
	"testing"

	require "github.com/stretchr/testify/require"
	duktape "gopkg.in/olebedev/go-duktape.v3"
)

func TestNewObjValidatorOptionalFields(t *testing.T) {

	v := NewObjValidator()
	v.Set("to", ObjTypeAddress, false)

	vm := duktape.New()
	err := vm.PevalString(`({})`)
	require.Nil(t, err)

	// expect no error since it's ok omit the address
	validationErr := v.Validate(vm, -1)
	require.Nil(t, validationErr)

}

func TestNewObjValidationRequiredFields(t *testing.T) {

	/**
	 * Test required field - no error
	 */
	v := NewObjValidator()
	v.Set("to", ObjTypeAddress, true)

	vm := duktape.New()
	err := vm.PevalString(`({"to": "0x0b078896a3d9166da5c37ae52a5809aca48630d4"})`)
	require.Nil(t, err)
	// expect no error since we provide the address
	require.Nil(t, v.Validate(vm, -1))

	/**
	 * Test required field - with error
	 */
	v = NewObjValidator()
	v.Set("to", ObjTypeAddress, false)

	vm = duktape.New()
	err = vm.PevalString(`({})`)
	require.Nil(t, err)
	// expect error since we don't pass the required address in
	require.Nil(t, v.Validate(vm, -1))

}

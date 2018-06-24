package validator

import (
	"testing"

	require "github.com/stretchr/testify/require"
)

func TestValidator(t *testing.T) {

	v := New()

	v.Set(0, &TypeFunction)
	require.Equal(t, &TypeFunction, v.rules[0])

	v.Set(1, &TypeFunction)
	require.Equal(t, &TypeFunction, v.rules[1])

}

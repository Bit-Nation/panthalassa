package max_count

import (
	"errors"
	"testing"

	require "github.com/stretchr/testify/require"
)

func TestIncrease(t *testing.T) {
	mc := New(2, nil)
	require.Nil(t, mc.Increase())
	require.Equal(t, uint(1), mc.count)
}

func TestDecrease(t *testing.T) {
	mc := New(2, nil)
	mc.count = 1
	mc.Decrease()
	require.Equal(t, uint(0), mc.count)
}

func TestIncreaseError(t *testing.T) {
	mc := New(0, errors.New("test error"))
	require.EqualError(t, mc.Increase(), "test error")
}

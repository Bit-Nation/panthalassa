package request_limitation

import (
	"errors"
	"testing"

	require "github.com/stretchr/testify/require"
)

func TestCount_Increase(t *testing.T) {
	mc := NewCount(2, nil)
	require.Nil(t, mc.Increase())
	require.Equal(t, uint(1), mc.count)
}

func TestCount_Decrease(t *testing.T) {
	mc := NewCount(2, nil)
	mc.Increase()
	mc.Decrease()
	require.Equal(t, uint(0), mc.count)

	// should still be 0
	mc.Decrease()
	require.Equal(t, uint(0), mc.count)
}

func TestCount_Count(t *testing.T) {
	mc := NewCount(2, nil)
	mc.Increase()
	require.Equal(t, uint(1), mc.Count())
}

func TestCount_IncreaseError(t *testing.T) {
	mc := NewCount(0, errors.New("test error"))
	require.EqualError(t, mc.Increase(), "test error")
}

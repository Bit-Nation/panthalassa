package mnemonic

import (
	"encoding/hex"
	"testing"

	require "github.com/stretchr/testify/require"
	bip39 "github.com/tyler-smith/go-bip39"
)

func TestNewFromString(t *testing.T) {

	mne := "legal winner thank year wave sausage worth useful legal winner thank year wave sausage worth useful legal winner thank year wave sausage worth title"

	m, err := FromString(mne)
	require.Nil(t, err)
	require.Equal(t, mne, m.String())

}

func TestNew(t *testing.T) {

	m, err := New()
	require.Nil(t, err)

	bip39.IsMnemonicValid(m.String())

}

func TestSeed(t *testing.T) {

	mne := "legal winner thank year wave sausage worth useful legal winner thank year wave sausage worth useful legal winner thank year wave sausage worth title"
	expectedSeed, err := hex.DecodeString("bc09fca1804f7e69da93c2f2028eb238c227f2e9dda30cd63699232578480a4021b146ad717fbb7e451ce9eb835f43620bf5c514db0f8add49f5d121449d3e87")
	require.Nil(t, err)

	m, err := FromString(mne)
	require.Nil(t, err)

	newSeed, err := m.NewSeed("TREZOR")
	require.Nil(t, err)

	require.Equal(t, expectedSeed, newSeed)

}

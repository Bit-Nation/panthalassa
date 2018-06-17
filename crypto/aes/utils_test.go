package aes

import (
	"crypto/hmac"
	"crypto/sha256"
	"errors"
	"testing"

	require "github.com/stretchr/testify/require"
)

func TestVersionOneOfMac(t *testing.T) {

	secret := Secret{
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}

	ct := CipherText{
		CipherText: []byte("another fake value"),
		Version:    1,
	}

	// create mac of version one which consist only of the the cipher text (the byte value)
	h := hmac.New(sha256.New, secret[:])
	_, err := h.Write(ct.CipherText)
	require.Nil(t, err)

	// test version one mac function with expected value
	mac, err := vOneMac(ct, secret)
	require.Nil(t, err)
	require.Equal(t, h.Sum(nil), mac)

	// should exit if version is != 1
	mac, err = vOneMac(CipherText{}, secret)
	require.EqualError(t, errors.New("cipher text must be of version one"), err.Error())
	require.Nil(t, mac)

}

func TestVersionTwoOfMac(t *testing.T) {

	secret := Secret{
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}

	ct := CipherText{
		CipherText: []byte("another fake value"),
		Version:    2,
		IV:         []byte("fake IV"),
	}

	// create mac of version one which consist only of the the cipher text (the byte value)
	// cipher text + IV + Version
	h := hmac.New(sha256.New, secret[:])
	_, err := h.Write(ct.CipherText)
	require.Nil(t, err)

	_, err = h.Write(ct.IV)
	require.Nil(t, err)

	_, err = h.Write([]byte{ct.Version})
	require.Nil(t, err)

	// test version one mac function with expected value
	mac, err := vTwoMac(ct, secret)
	require.Nil(t, err)
	require.Equal(t, h.Sum(nil), mac)

	// should exit if version != 2
	mac, err = vTwoMac(CipherText{}, secret)
	require.EqualError(t, errors.New("cipher text must be of version two"), err.Error())
	require.Nil(t, mac)

}

// verify mac
func TestVerifyMACInvalidVersion(t *testing.T) {

	// should exit with error since version != 1 || != 2
	ct := CipherText{
		Version: uint8(5),
	}
	valid, err := ct.ValidMAC(Secret{})
	require.EqualError(t, err, "failed to verify MAC since we don't know how to handle version: 5")
	require.False(t, valid)

}

func TestVerifyMACVersionOne(t *testing.T) {

	// should be valid since we can handle verion one
	oldVOneMac := vOneMac
	sec := Secret{
		0x23,
	}
	ct := CipherText{
		Version: 1,
		Mac:     []byte("i am a version one mac"),
	}
	// we override the mac function for better testing
	vOneMac = func(ct CipherText, secret Secret) ([]byte, error) {
		require.Equal(t, sec, secret)
		require.Equal(t, []byte("i am a version one mac"), ct.Mac)
		return []byte("i am a version one mac"), nil
	}
	valid, err := ct.ValidMAC(sec)
	require.Nil(t, err)
	require.True(t, valid)

	// reset vOneMac function
	vOneMac = oldVOneMac

}

func TestVerifyMACVersionTwo(t *testing.T) {

	// should be valid since we can handle verion one
	oldVTwoMac := vOneMac
	sec := Secret{
		0x23,
	}
	ct := CipherText{
		Version: 2,
		Mac:     []byte("i am a version one mac"),
	}

	// we override the mac function for better testing
	vTwoMac = func(ct CipherText, secret Secret) ([]byte, error) {
		require.Equal(t, sec, secret)
		require.Equal(t, []byte("i am a version one mac"), ct.Mac)
		return []byte("i am a version one mac"), nil
	}
	valid, err := ct.ValidMAC(sec)
	require.Nil(t, err)
	require.True(t, valid)

	// reset vOneMac function
	vTwoMac = oldVTwoMac

}

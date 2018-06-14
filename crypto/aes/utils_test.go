package aes

import (
	"crypto/hmac"
	"crypto/sha256"
	"github.com/kataras/iris/core/errors"
	"github.com/stretchr/testify/require"
	"testing"
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

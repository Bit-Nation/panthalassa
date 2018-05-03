package cid

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCIDSha256(t *testing.T) {
	cid, err := CIDSha256("HI")
	require.Nil(t, err)
	require.Equal(t, "mAVUWIDAwU3dua7xpbnFfMuJGB3ydE8z6o/Oz8+bihw3lVhT1", cid)
}

func TestCIDSha512(t *testing.T) {
	cid, err := CIDSha512("HI")
	require.Nil(t, err)
	require.Equal(t, "UAVUUQEG8zpe5V7XbE9ZuG_0-q4q0ujKrksvGUoJTCs2fTjXSUvfOE6Z4f9ePhYante_RzbqCZT5RADYOegO_Yiz8WJo=", cid)
}

func TestIsValidCid(t *testing.T) {
	require.True(t, IsValidCid("mAVUWIDAwU3dua7xpbnFfMuJGB3ydE8z6o/Oz8+bihw3lVhT1"))
	require.True(t, IsValidCid("UAVUUQEG8zpe5V7XbE9ZuG_0-q4q0ujKrksvGUoJTCs2fTjXSUvfOE6Z4f9ePhYante_RzbqCZT5RADYOegO_Yiz8WJo="))
	require.False(t, IsValidCid("I should be invalid"))
}

package keyStore

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGetKey(t *testing.T) {
	ks := KeyStore{
		keys: map[string]string{
			"key": "value",
		},
	}

	value, err := ks.GetKey("key")
	require.Equal(t, "value", value)
	require.Nil(t, err)

	value, err = ks.GetKey("i_do_not_exist_in_the_map")
	require.Equal(t, "key does not exist", err.Error())
}

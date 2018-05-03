package ethereum

import (
	"errors"
	"github.com/Bit-Nation/panthalassa/mnemonic"
	"github.com/stretchr/testify/require"
	"testing"
)

var testMnemonic = "panda eyebrow bullet gorilla call smoke muffin taste mesh discover soft ostrich alcohol speed nation flash devote level hobby quick inner drive ghost inside"

func TestUpSuccess(t *testing.T) {

	//mnemonic
	mnemo, err := mnemonic.FromString(testMnemonic)
	require.Nil(t, err)

	//Migration
	m := Migration{}

	keys := make(map[string]string)

	//Migrate up
	migratedKeys, err := m.Up(mnemo, keys)
	require.Nil(t, err)

	require.Equal(t, "f84d5d4808521ae7330607cbbd0503959659b927f24db70421fc551e05b50409", migratedKeys[EthereumKey])

}

func TestUpFailOnWrongValue(t *testing.T) {

	//mnemonic
	mnemo, err := mnemonic.FromString(testMnemonic)
	require.Nil(t, err)

	//Migration
	m := Migration{}

	keys := map[string]string{
		EthereumKey: "i_am_not_an_private_key",
	}

	//Migrate up
	_, err = m.Up(mnemo, keys)

	//Should fail since the derived value will NOT match with the key that already exist
	require.Error(t, errors.New("private key already exist BUT does not match derived private key"), err)

}

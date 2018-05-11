package encryption_key

import (
	"github.com/Bit-Nation/panthalassa/mnemonic"
	"github.com/stretchr/testify/require"
	"testing"
)

const testMnemonic = "panda eyebrow bullet gorilla call smoke muffin taste mesh discover soft ostrich alcohol speed nation flash devote level hobby quick inner drive ghost inside"
const testMnemonicWrong = "eyebrow eyebrow bullet gorilla call smoke muffin taste mesh discover soft ostrich alcohol speed nation flash devote level hobby quick inner drive ghost inside"

func TestSuccessMigration(t *testing.T) {

	m := Migration{}

	//Test mnemonic
	mne, err := mnemonic.FromString(testMnemonic)
	require.Nil(t, err)

	//Migrate
	keys := map[string]string{}
	keys, err = m.Up(mne, keys)
	require.Nil(t, err)

	require.Equal(t, "4b08e01c2b0d41c056bd862be4b953f67d2a4d15272fa7c0419c6d897c016790", keys[BIP39Password])

}

func TestFailMigration(t *testing.T) {

	m := Migration{}

	//Test mnemonic
	mne, err := mnemonic.FromString(testMnemonic)
	require.Nil(t, err)

	//Migrate
	keys := map[string]string{}
	keys, err = m.Up(mne, keys)
	require.Nil(t, err)
	require.Equal(t, "4b08e01c2b0d41c056bd862be4b953f67d2a4d15272fa7c0419c6d897c016790", keys[BIP39Password])

	//Test Fail migration with wrong mnemonic
	mne, err = mnemonic.FromString(testMnemonicWrong)
	keys, err = m.Up(mne, keys)
	require.EqualError(t, err, "migration - derived key miss match with existing key")

}

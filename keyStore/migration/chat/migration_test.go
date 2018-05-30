package chat

import (
	"testing"

	"github.com/Bit-Nation/panthalassa/mnemonic"
	require "github.com/stretchr/testify/require"
)

const testMnemonic = "panda eyebrow bullet gorilla call smoke muffin taste mesh discover soft ostrich alcohol speed nation flash devote level hobby quick inner drive ghost inside"
const testMnemonicWrong = "eyebrow eyebrow bullet gorilla call smoke muffin taste mesh discover soft ostrich alcohol speed nation flash devote level hobby quick inner drive ghost inside"

func TestSuccessfulMigration(t *testing.T) {

	m := Migration{}
	mne, err := mnemonic.FromString(testMnemonic)
	require.Nil(t, err)

	keys := map[string]string{}

	migrated, err := m.Up(mne, keys)
	require.Nil(t, err)

	require.Equal(t, "c0d4b1ea7b04c2d6e3c95b63c0dff3ce841bf5b14b2904910bf82ae86623d57a", migrated[MigrationPrivPrefix])
	require.Equal(t, "3c78ae7aa62692b5840ebdb40cb12a2283b754a11d61c0c3e5a64a28d6607d75", migrated[MigrationPubPrefix])

	// migrate again to make sure we still get the same output with a second run
	migrated, err = m.Up(mne, keys)

	require.Equal(t, "c0d4b1ea7b04c2d6e3c95b63c0dff3ce841bf5b14b2904910bf82ae86623d57a", migrated[MigrationPrivPrefix])
	require.Equal(t, "3c78ae7aa62692b5840ebdb40cb12a2283b754a11d61c0c3e5a64a28d6607d75", migrated[MigrationPubPrefix])

}

func TestWrongMigration(t *testing.T) {

	// first migrate with the correct mnemonic
	m := Migration{}
	mne, err := mnemonic.FromString(testMnemonic)
	require.Nil(t, err)

	keys := map[string]string{}

	migrated, err := m.Up(mne, keys)
	require.Nil(t, err)

	require.Equal(t, "c0d4b1ea7b04c2d6e3c95b63c0dff3ce841bf5b14b2904910bf82ae86623d57a", migrated[MigrationPrivPrefix])
	require.Equal(t, "3c78ae7aa62692b5840ebdb40cb12a2283b754a11d61c0c3e5a64a28d6607d75", migrated[MigrationPubPrefix])

	// now migrate again with wrong mnemonic to trigger error
	mne, err = mnemonic.FromString(testMnemonicWrong)
	require.Nil(t, err)

	m = Migration{}
	_, err = m.Up(mne, migrated)
	require.EqualError(t, err, "migration (chat) derivation miss match of private key")

}

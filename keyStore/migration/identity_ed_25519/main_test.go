package identity_ed_25519

import (
	"github.com/Bit-Nation/panthalassa/mnemonic"
	"github.com/stretchr/testify/require"
	"testing"
)

const testMnemonic = "panda eyebrow bullet gorilla call smoke muffin taste mesh discover soft ostrich alcohol speed nation flash devote level hobby quick inner drive ghost inside"

func TestMigration_Up(t *testing.T) {

	//Create mnemonic
	mne, err := mnemonic.FromString(testMnemonic)
	require.Nil(t, err)

	keys := map[string]string{}

	mig := Migration{}
	keys, err = mig.Up(mne, keys)
	require.Nil(t, err)

	require.Equal(t, "898e73f4f72e15fe9acf80129ed0a16c500ef2656e0b26369fef6728424b2f68d368d5c5089abca526663524bdac8b8210ec82000ca5e084a15a45a27a9f8666", keys[Ed25519PrivateKey])
	require.Equal(t, "d368d5c5089abca526663524bdac8b8210ec82000ca5e084a15a45a27a9f8666", keys[Ed25519PublicKey])

}

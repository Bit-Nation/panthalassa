package ed25519

import (
	"github.com/Bit-Nation/panthalassa/mnemonic"
	"github.com/stretchr/testify/require"
	"github.com/Bit-Nation/panthalassa/keyStore/migration/identity"
	"testing"
)

const testMnemonic = "panda eyebrow bullet gorilla call smoke muffin taste mesh discover soft ostrich alcohol speed nation flash devote level hobby quick inner drive ghost inside"
const testMnemonicWrong = "eyebrow eyebrow bullet gorilla call smoke muffin taste mesh discover soft ostrich alcohol speed nation flash devote level hobby quick inner drive ghost inside"

func TestMigration_UpSuccess(t *testing.T) {

	//Create mnemonic
	mne, err := mnemonic.FromString(testMnemonic)
	require.Nil(t, err)

	keys := map[string]string{}

	mig := Migration{}
	keys, err = mig.Up(mne, keys)
	require.Nil(t, err)

	require.Equal(t, "898e73f4f72e15fe9acf80129ed0a16c500ef2656e0b26369fef6728424b2f68d368d5c5089abca526663524bdac8b8210ec82000ca5e084a15a45a27a9f8666", keys[identity.Ed25519PrivateKey])
	require.Equal(t, "d368d5c5089abca526663524bdac8b8210ec82000ca5e084a15a45a27a9f8666", keys[identity.Ed25519PublicKey])

}

func TestMigration_UpDerivationMissMatch(t *testing.T) {

	//Create mnemonic
	mne, err := mnemonic.FromString(testMnemonic)
	require.Nil(t, err)

	keys := map[string]string{}

	//Migrate success
	mig := Migration{}
	keys, err = mig.Up(mne, keys)
	require.Nil(t, err)

	//Create wrong mnemonic. This will be used with an key store
	//that was migrated with another mnemonic
	wrongMne, err := mnemonic.FromString(testMnemonicWrong)
	require.Nil(t, err)

	//Test migration with different mnemonic -> should result in an error for private key
	_, err = mig.Up(wrongMne, keys)
	require.EqualError(t, err, "migration - ed25519 private key derivation miss match")

	//Test migration with different mnemonic -> should result in an error for public key
	//Delete private key from key store. That will make the check for the private key pass since it doesn't exist
	delete(keys, identity.Ed25519PrivateKey)
	_, err = mig.Up(wrongMne, keys)
	require.EqualError(t, err, "migration - ed25519 public key derivation miss match")

}

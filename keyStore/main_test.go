package keyStore

import (
	"errors"
	"testing"

	"github.com/Bit-Nation/panthalassa/keyStore/migration"
	"github.com/Bit-Nation/panthalassa/mnemonic"
	"github.com/stretchr/testify/require"
)

var oldMigrations = migrations

var testMnemonic = "panda eyebrow bullet gorilla call smoke muffin taste mesh discover soft ostrich alcohol speed nation flash devote level hobby quick inner drive ghost inside"

//Test fail migration
type errorMigration struct{}

func (m errorMigration) Up(mnemonic mnemonic.Mnemonic, keys map[string]string) (map[string]string, error) {
	return keys, errors.New("I am a migration test error")
}

//Test success migration
type successMigration struct{}

func (m successMigration) Up(mnemonic mnemonic.Mnemonic, keys map[string]string) (map[string]string, error) {

	if mnemonic.String() != testMnemonic {
		panic("Got invalid mnemonic")
	}

	keys["new_key"] = "new_value"
	return keys, nil
}

func TestMigrateUpSuccess(t *testing.T) {

	//Mnemonic
	m, err := mnemonic.FromString(testMnemonic)
	require.Nil(t, err)

	//Override migration
	migrations = []keys.Migration{
		successMigration{},
	}

	s := Store{
		mnemonic: m,
		keys:     make(map[string]string),
		version:  1,
	}

	//Migrate up
	s, err = migrateUp(s)
	require.Nil(t, err)

	//Keys should now contain key "new_key" with value "new_value"
	require.Equal(t, s.keys["new_key"], "new_value")
	require.True(t, s.changed)

	//Reset migrations
	migrations = oldMigrations

}

//Exit on error in migration
func TestMigrateUpError(t *testing.T) {

	//Mnemonic
	m, err := mnemonic.FromString(testMnemonic)
	require.Nil(t, err)

	//Override migration
	migrations = []keys.Migration{
		errorMigration{},
	}

	s := Store{
		mnemonic: m,
		keys:     make(map[string]string),
		version:  1,
	}

	//Migrate up should exit with an error from the migration
	s, err = migrateUp(s)
	require.Error(t, errors.New("I am a migration test error"), err)

	//Reset migrations
	migrations = oldMigrations

}

func TestNewFromMnemonic(t *testing.T) {

	//mnemonic
	m, err := mnemonic.FromString(testMnemonic)
	require.Nil(t, err)

	//Create store
	store, err := NewFromMnemonic(m)

	//Do some assertions
	require.Equal(t, testMnemonic, store.mnemonic.String())
	require.Equal(t, uint8(1), store.version)
	require.Equal(t, "f84d5d4808521ae7330607cbbd0503959659b927f24db70421fc551e05b50409", store.keys["ethereum_private_key"])

}

func TestGetKey(t *testing.T) {

	//mnemonic
	m, err := mnemonic.FromString(testMnemonic)
	require.Nil(t, err)

	s := Store{
		mnemonic: m,
		keys: map[string]string{
			"key": "value",
		},
		version: 1,
	}

	//Get key that exist
	v, err := s.GetKey("key")
	require.Nil(t, err)
	require.Equal(t, "value", v)

	//Get key that doesn't exist
	v, err = s.GetKey("not_present")
	require.Error(t, errors.New(""), err)
	require.Equal(t, "", v)

}

func TestMarshal(t *testing.T) {

	expectedStore := `{"mnemonic":"panda eyebrow bullet gorilla call smoke muffin taste mesh discover soft ostrich alcohol speed nation flash devote level hobby quick inner drive ghost inside","keys":{"key":"value"},"version":1}`

	//mnemonic
	m, err := mnemonic.FromString(testMnemonic)
	require.Nil(t, err)

	s := Store{
		mnemonic: m,
		keys: map[string]string{
			"key": "value",
		},
		version: 1,
	}

	store, err := s.Marshal()
	require.Nil(t, err)

	require.Equal(t, expectedStore, string(store))

}

func TestUnmarshal(t *testing.T) {

	//Test json key store
	jsonStore := `{"mnemonic":"panda eyebrow bullet gorilla call smoke muffin taste mesh discover soft ostrich alcohol speed nation flash devote level hobby quick inner drive ghost inside","keys":{"key":"value"},"version":1}`

	//Unmarshal the store
	s, err := UnmarshalStore(jsonStore)
	require.Nil(t, err)

	//Assert it's the same
	require.Equal(t, "panda eyebrow bullet gorilla call smoke muffin taste mesh discover soft ostrich alcohol speed nation flash devote level hobby quick inner drive ghost inside", s.mnemonic.String())
	require.Equal(t, "value", s.keys["key"])
	require.Equal(t, uint8(1), s.version)

}

package keyStore

import (
	"encoding/json"
	"errors"
	"reflect"

	migration "github.com/Bit-Nation/panthalassa/keyStore/migration"
	chatMigration "github.com/Bit-Nation/panthalassa/keyStore/migration/chat"
	encryptionKeyMigr "github.com/Bit-Nation/panthalassa/keyStore/migration/encryption_key"
	ethereumMigration "github.com/Bit-Nation/panthalassa/keyStore/migration/ethereum"
	ed25519Migration "github.com/Bit-Nation/panthalassa/keyStore/migration/identity/ed25519"
	mnemonic "github.com/Bit-Nation/panthalassa/mnemonic"
)

//Migrations to run
var migrations = []migration.Migration{
	ethereumMigration.Migration{},
	ed25519Migration.Migration{},
	chatMigration.Migration{},
	encryptionKeyMigr.Migration{},
}

type Store struct {
	mnemonic mnemonic.Mnemonic
	keys     map[string]string
	version  uint8
	changed  bool
}

type jsonStore struct {
	Mnemonic string            `json:"mnemonic"`
	Keys     map[string]string `json:"keys"`
	Version  uint8             `json:"version"`
}

//Return the plain keys store
func (s Store) Marshal() (string, error) {

	//Json representation
	js := jsonStore{
		Mnemonic: s.mnemonic.String(),
		Keys:     s.keys,
		Version:  s.version,
	}

	//Marshal the whole thing
	b, err := json.Marshal(js)
	if err != nil {
		return "", err
	}

	//transform to string
	return string(b), nil

}

//Get a value from the keystore
func (s Store) GetKey(key string) (string, error) {

	_, exist := s.keys[key]
	if !exist {
		return "", errors.New("key does not exist")
	}

	return s.keys[key], nil
}

//Get the mnemonic
func (s Store) GetMnemonic() mnemonic.Mnemonic {
	return s.mnemonic
}

//Did the keystore changed (happen after a migration)
func (s Store) WasMigrated() bool {
	return s.changed
}

//Migrate keystore up
func migrateUp(s Store) (Store, error) {

	oldKeys := s.keys

	//Mutate key store
	for _, m := range migrations {

		//Migrate Up
		//@todo change to by value
		newKeys, err := m.Up(s.mnemonic, s.keys)
		if err != nil {
			return Store{}, err
		}

		//Assign new key's
		s.keys = newKeys

	}

	//Check if there are any new key's
	s.changed = reflect.DeepEqual(oldKeys, s.keys)

	return s, nil

}

func UnmarshalStore(keyStore string) (Store, error) {

	var js jsonStore

	//Unmarshal key store
	if err := json.Unmarshal([]byte(keyStore), &js); err != nil {
		return Store{}, err
	}

	//create mnemotic from string representation
	m, err := mnemonic.FromString(js.Mnemonic)
	if err != nil {
		return Store{}, err
	}

	//Usable keystore
	s := Store{
		mnemonic: m,
		keys:     js.Keys,
		version:  js.Version,
	}

	//Migrate the keystore
	s, err = migrateUp(s)
	if err != nil {
		return Store{}, err
	}

	return s, err

}

//Create a new store from the mnemonic and migrate it
func NewFromMnemonic(mnemonic mnemonic.Mnemonic) (Store, error) {

	//Store
	s := Store{
		mnemonic: mnemonic,
		keys:     make(map[string]string),
		version:  1,
	}

	//Migrate store
	s, err := migrateUp(s)
	if err != nil {
		return Store{}, err
	}

	return s, nil

}

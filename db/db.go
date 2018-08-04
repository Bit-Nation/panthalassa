package db

import (
	"os"
	"path/filepath"

	migration "github.com/Bit-Nation/panthalassa/db/migration"
	km "github.com/Bit-Nation/panthalassa/keyManager"
	bolt "github.com/coreos/bbolt"
)

// get database path for key manager
func KMToDBPath(dir string, km *km.KeyManager) (string, error) {

	idPubKey, err := km.IdentityPublicKey()
	if err != nil {
		return "", err
	}

	return filepath.Abs(filepath.Join(dir, idPubKey+".db"))

}

// open a database
func Open(path, userPath string, mode os.FileMode, options *bolt.Options) (*bolt.DB, error) {

	migrations := []migration.Migration{}

	// check if production database exist
	if _, err := os.Stat(path); err == nil {
		// migrate the database
		err := migration.Migrate(path, userPath, migrations)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, err
	}

	// open database
	db, err := bolt.Open(path, mode, options)
	if err != nil {
		return nil, err
	}

	return db, nil

}

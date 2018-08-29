package db

import (
	"os"
	"path/filepath"

	migration "github.com/Bit-Nation/panthalassa/db/migration"
	km "github.com/Bit-Nation/panthalassa/keyManager"
	storm "github.com/asdine/storm"
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
func Open(path string, mode os.FileMode, options *bolt.Options) (*storm.DB, error) {

	migrations := []migration.Migration{}

	// migrate the database
	err := migration.Migrate(path, migrations)
	if err != nil {
		return nil, err
	}

	// open database
	db, err := storm.Open(path, storm.BoltOptions(mode, options))
	if err != nil {
		return nil, err
	}

	return db, nil

}

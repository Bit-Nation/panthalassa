package db

import (
	"os"

	migration "github.com/Bit-Nation/panthalassa/db/migration"
	bolt "github.com/coreos/bbolt"
)

type DB struct {
	bolt *bolt.DB
}

func (db *DB) Close() error {
	return db.bolt.Close()
}

// open a database
func Open(path string, mode os.FileMode, options *bolt.Options) (*DB, error) {

	migrations := []migration.Migration{}

	// migrate the database
	err := migration.Migrate(path, migrations)
	if err != nil {
		return nil, err
	}

	// open database
	db, err := bolt.Open(path, mode, options)
	if err != nil {
		return nil, err
	}

	return &DB{
		bolt: db,
	}, err

}

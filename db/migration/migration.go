package migration

import (
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"

	bolt "github.com/coreos/bbolt"
)

const dbFileMode = 0600

type Migration interface {
	Migrate(db *bolt.DB) error
	Version() uint32
}

func randomTempDBPath() (string, error) {
	file := make([]byte, 50)
	_, err := rand.Read(file)
	if err != nil {
		return "", err
	}
	return filepath.Abs(filepath.Join(os.TempDir(), hex.EncodeToString(file)+".bolt"))
}

func findNextMigration(currentVersion uint32, migrations []Migration) (Migration, error) {
	mig := doubleMigration(migrations)
	if mig != nil {
		return nil, errors.New(fmt.Sprintf("migration %d is already registered", mig.Version()))
	}
	// sort the migrations based on there version
	sort.SliceStable(migrations, func(i, j int) bool {
		return migrations[i].Version() < migrations[j].Version()
	})
	for _, m := range migrations {
		if currentVersion < m.Version() {
			return m, nil
		}
	}
	return nil, nil
}

// returns a migration if a version with the same ID already exist
func doubleMigration(migrations []Migration) Migration {
	existingMigrations := map[uint32]bool{}
	for _, m := range migrations {
		_, exist := existingMigrations[m.Version()]
		if exist {
			return m
		}
		existingMigrations[m.Version()] = true
	}
	return nil
}

// migrate migrates a bold db database
func Migrate(prodDBPath string, migrations []Migration) error {

	// make sure there are no double migrations
	mig := doubleMigration(migrations)
	if mig != nil {
		return errors.New(fmt.Sprintf("a migration with id (%d) was already registered", mig.Version()))
	}

	// check if production database exist
	if _, err := os.Stat(prodDBPath); err != nil {
		return err
	}

	// open production database
	prodDB, err := bolt.Open(prodDBPath, dbFileMode, &bolt.Options{Timeout: time.Second})
	if err != nil {
		return err
	}

	// migration database
	// @todo maybe a more random value should be used
	migrationDBFile, err := randomTempDBPath()
	if err != nil {
		return err
	}
	// copy production database over to migration database
	err = prodDB.View(func(tx *bolt.Tx) error {
		return tx.CopyFile(migrationDBFile, dbFileMode)
	})
	defer prodDB.Close()
	if err != nil {
		return err
	}

	// open migration database
	migrationDB, err := bolt.Open(migrationDBFile, dbFileMode, bolt.DefaultOptions)
	if err != nil {
		return err
	}

	for {
		var version uint32

		// read version of the migration database
		err := migrationDB.View(func(tx *bolt.Tx) error {
			// fetch system bucket
			buck := tx.Bucket([]byte("system"))
			if buck == nil {
				return errors.New("system bucket does not exist")
			}
			// fetch database version
			ver := buck.Get([]byte("database_version"))
			if ver == nil {
				return errors.New("no version present in this schema")
			}
			version = binary.BigEndian.Uint32(ver)
			return nil
		})
		if err != nil {
			return err
		}

		// find next migration
		nextMigr, err := findNextMigration(version, migrations)
		if err != nil {
			return err
		}

		// exit if there is no next migration
		if nextMigr == nil {
			break
		}

		// migrate up
		if err := nextMigr.Migrate(migrationDB); err != nil {
			return err
		}

		// update version of last migration on database
		err = migrationDB.Update(func(tx *bolt.Tx) error {
			// fetch bucket
			buck := tx.Bucket([]byte("system"))
			if buck == nil {
				return errors.New("system bucket does not exist")
			}
			v := make([]byte, 4)
			binary.BigEndian.PutUint32(v, nextMigr.Version())
			return buck.Put([]byte("database_version"), v)
		})
		if err != nil {
			return err
		}
	}

	// create path for production backup DB
	prodDBBackupPath, err := filepath.Abs(filepath.Join(os.TempDir(), time.Now().String()))
	if err != nil {
		return err
	}
	prodBackupDB, err := bolt.Open(prodDBBackupPath, dbFileMode, &bolt.Options{Timeout: time.Second})
	if err != nil {
		return err
	}

	// backup production database
	err = prodDB.View(func(tx *bolt.Tx) error {
		return tx.CopyFile(prodDBBackupPath, dbFileMode)
	})
	if err != nil {
		return err
	}

	// remove old production database
	if err := prodDB.Close(); err != nil {
		return err
	}
	if err := os.Remove(prodDBPath); err != nil {
		return err
	}

	// copy migrated database to production
	err = migrationDB.View(func(tx *bolt.Tx) error {
		return tx.CopyFile(prodDBPath, dbFileMode)
	})
	if err != nil {
		// try to recover production database
		recErr := prodBackupDB.Update(func(tx *bolt.Tx) error {
			return tx.CopyFile(prodDBPath, dbFileMode)
		})
		if recErr != nil {
			return errors.New(fmt.Sprintf(
				"failed to recover database after attempt to copy over migrated db. original error: %s. database recover error: %s",
				err.Error(),
				recErr.Error(),
			))
		}
		return err
	}

	// delete migration database
	if err := migrationDB.Close(); err != nil {
		return err
	}
	return os.Remove(migrationDBFile)

}

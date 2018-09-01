package db

import (
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"sort"
	"time"

	storm "github.com/asdine/storm"
	bolt "github.com/coreos/bbolt"
)

type Migration interface {
	Migrate(db *storm.DB) error
	Version() uint32
}

var systemBucketName = []byte("system")

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

var setupSystemBucket = func(db *storm.DB) error {
	return db.Bolt.Update(func(tx *bolt.Tx) error {

		// fetch system bucket
		systemBucket, err := tx.CreateBucketIfNotExists(systemBucketName)
		if err != nil {
			return err
		}

		// get database version
		dbVersion := systemBucket.Get([]byte("database_version"))
		if dbVersion == nil {
			v := make([]byte, 4)
			binary.BigEndian.PutUint32(v, 0)
			return systemBucket.Put([]byte("database_version"), v)
		}
		return nil
	})
}

// create a migration file and create the directory
// structure needed for the file
var prepareMigration = func(file string) (string, error) {

	// create migration file
	currentDir := path.Dir(file)
	migFileName := make([]byte, 50)
	if _, err := rand.Read(migFileName); err != nil {
		return "", err
	}
	migrationFile, err := filepath.Abs(filepath.Join(currentDir, hex.EncodeToString(migFileName)+".db.backup"))
	if err != nil {
		return "", err
	}

	return migrationFile, nil
}

// migrate migrates a bold db database
func Migrate(prodDBFile string, migrations []Migration) error {

	// make sure there are no double migrations
	mig := doubleMigration(migrations)
	if mig != nil {
		return errors.New(fmt.Sprintf("a migration with id (%d) was already registered", mig.Version()))
	}

	// check if production database exist
	if _, err := os.Stat(prodDBFile); err != nil {
		return err
	}

	// open production database
	prodDB, err := storm.Open(prodDBFile, storm.BoltOptions(0644, &bolt.Options{Timeout: time.Second}))
	defer prodDB.Close()
	if err != nil {
		return err
	}

	// make sure the system bucket exist and is setup correct
	if err := setupSystemBucket(prodDB); err != nil {
		return err
	}

	// backup production database
	dbBackupFile, err := prepareMigration(prodDBFile)
	if err != nil {
		return err
	}
	err = prodDB.Bolt.View(func(tx *bolt.Tx) error {
		return tx.CopyFile(dbBackupFile, 0644)
	})
	if err != nil {
		return err
	}
	dbBackup, err := bolt.Open(dbBackupFile, 0644, &bolt.Options{Timeout: time.Second})
	defer dbBackup.Close()
	if err != nil {
		return err
	}

	for {
		var version uint32

		// read version of the migration database
		err := prodDB.Bolt.View(func(tx *bolt.Tx) error {
			// fetch system bucket
			buck := tx.Bucket(systemBucketName)
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
		if migErr := nextMigr.Migrate(prodDB); migErr != nil {

			// recover old database on error
			// which requires to copy the backup to production
			// and delete the migration DB
			if err := prodDB.Close(); err != nil {
				// @todo is it really clever to return at this point? Maybe we can continue the recovery process even with an close error.
				return err
			}
			err := dbBackup.View(func(tx *bolt.Tx) error {
				return tx.CopyFile(prodDBFile, 0644)
			})
			if err != nil {
				return err
			}

			return migErr

		}

		// update version of last migration on database
		err = prodDB.Bolt.Update(func(tx *bolt.Tx) error {
			// fetch bucket
			buck := tx.Bucket(systemBucketName)
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

	return os.Remove(dbBackupFile)

}

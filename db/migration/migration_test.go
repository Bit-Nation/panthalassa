package migration

import (
	"encoding/binary"
	"errors"
	"testing"
	"time"

	bolt "github.com/coreos/bbolt"
	require "github.com/stretchr/testify/require"
)

type migration struct {
	version           uint32
	migrationFunction func(m *bolt.DB) error
}

func (m *migration) Version() uint32 {
	return m.version
}
func (m *migration) Migrate(db *bolt.DB) error {
	return m.migrationFunction(db)
}

func TestDoubleMigration(t *testing.T) {

	type vector struct {
		migrations []Migration
		assertion  func(migration Migration)
	}

	vectors := []vector{
		vector{
			migrations: []Migration{
				&migration{version: 3},
				&migration{version: 1},
				&migration{version: 7},
				&migration{version: 1},
				&migration{version: 2},
			},
			assertion: func(migration Migration) {
				require.NotNil(t, migration)
				require.Equal(t, uint32(1), migration.Version())
			},
		},
		vector{
			migrations: []Migration{
				&migration{version: 1},
				&migration{version: 2},
			},
			assertion: func(migration Migration) {
				require.Nil(t, migration)
			},
		},
	}

	for _, v := range vectors {
		v.assertion(doubleMigration(v.migrations))
	}

}

func TestFindNextMigration(t *testing.T) {

	type vector struct {
		migrations     []Migration
		assertion      func(migration Migration, err error)
		currentVersion uint32
	}

	vectors := []vector{
		vector{
			migrations: []Migration{
				&migration{version: 3},
				&migration{version: 1},
				&migration{version: 7},
			},
			assertion: func(migration Migration, err error) {
				// expect to be nil since there is
				// no migration that is greater than 8
				require.Nil(t, migration)
			},
			currentVersion: 8,
		},
		vector{
			migrations: []Migration{
				&migration{version: 3},
				&migration{version: 2},
				&migration{version: 1},
			},
			currentVersion: 1,
			assertion: func(migration Migration, err error) {
				require.NotNil(t, migration)
				// Expect to be 2 since our current version is one
				require.Equal(t, uint32(2), migration.Version())
			},
		},
	}

	for _, v := range vectors {
		v.assertion(findNextMigration(v.currentVersion, v.migrations))
	}

}

func TestMigrateDoubleMigration(t *testing.T) {

	migrations := []Migration{
		&migration{version: 3},
		&migration{version: 3},
		&migration{version: 1},
	}

	require.EqualError(t, Migrate("", migrations), "a migration with id (3) was already registered")

}

func TestMigrationDBDoesNotExist(t *testing.T) {

	require.EqualError(
		t,
		Migrate("i_do_not_exist.db", []Migration{}),
		"stat i_do_not_exist.db: no such file or directory",
	)

}

func TestMigrateTimeoutOnOpen(t *testing.T) {

	path, err := randomTempDBPath()
	require.Nil(t, err)

	db, err := bolt.Open(path, dbFileMode, &bolt.Options{Timeout: time.Second})
	require.Nil(t, err)

	migrations := []Migration{
		&migration{version: 1},
	}

	require.EqualError(t, Migrate(db.Path(), migrations), "timeout")

}

func TestMigrateSystemBucketError(t *testing.T) {

	dbPath, err := randomTempDBPath()
	require.Nil(t, err)

	db, err := bolt.Open(dbPath, dbFileMode, &bolt.Options{Timeout: time.Second})
	require.Nil(t, err)

	migrations := []Migration{
		&migration{version: 1},
	}

	require.Nil(t, db.Close())

	require.EqualError(t, Migrate(dbPath, migrations), "system bucket does not exist")

}

func TestMigrateVersionError(t *testing.T) {

	dbPath, err := randomTempDBPath()
	require.Nil(t, err)

	db, err := bolt.Open(dbPath, dbFileMode, &bolt.Options{Timeout: time.Second})
	require.Nil(t, err)

	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucket([]byte("system"))
		return err
	})
	require.Nil(t, err)

	migrations := []Migration{
		&migration{version: 1},
	}

	require.Nil(t, db.Close())

	require.EqualError(t, Migrate(dbPath, migrations), "no version present in this schema")

}

func TestMigrateSuccess(t *testing.T) {

	dbPath, err := randomTempDBPath()
	require.Nil(t, err)

	db, err := bolt.Open(dbPath, dbFileMode, &bolt.Options{Timeout: time.Second})
	require.Nil(t, err)

	err = db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucket([]byte("system"))
		if err != nil {
			return err
		}
		v := make([]byte, 4)
		binary.BigEndian.PutUint32(v, 0)
		b.Put([]byte("database_version"), v)
		return err
	})
	require.Nil(t, err)

	migrations := []Migration{
		&migration{
			version: 1,
			migrationFunction: func(m *bolt.DB) error {
				return m.Update(func(tx *bolt.Tx) error {
					b, err := tx.CreateBucket([]byte("key_value_store"))
					if err != nil {
						return err
					}
					return b.Put([]byte("profile"), []byte("Florian"))
				})
			},
		},
	}

	require.Nil(t, db.Close())
	require.Nil(t, Migrate(dbPath, migrations))

	db, err = bolt.Open(dbPath, dbFileMode, &bolt.Options{Timeout: time.Second})
	require.Nil(t, err)

	err = db.View(func(tx *bolt.Tx) error {
		kvs := tx.Bucket([]byte("key_value_store"))
		require.Equal(t, string(kvs.Get([]byte("profile"))), "Florian")
		systemBucket := tx.Bucket([]byte("system"))
		require.Equal(t, uint32(1), binary.BigEndian.Uint32(systemBucket.Get([]byte("database_version"))))
		return nil
	})
	require.Nil(t, err)
}

// in the case of an error the database should not change
func TestMigrateRevertChangesOnError(t *testing.T) {

	dbPath, err := randomTempDBPath()
	require.Nil(t, err)

	db, err := bolt.Open(dbPath, dbFileMode, &bolt.Options{Timeout: time.Second})
	require.Nil(t, err)

	err = db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucket([]byte("system"))
		if err != nil {
			return err
		}
		v := make([]byte, 4)
		binary.BigEndian.PutUint32(v, 0)
		b.Put([]byte("database_version"), v)
		return err
	})
	require.Nil(t, err)

	migrations := []Migration{
		&migration{
			version: 1,
			migrationFunction: func(m *bolt.DB) error {
				err := m.Update(func(tx *bolt.Tx) error {
					b, err := tx.CreateBucket([]byte("key_value_store"))
					if err != nil {
						return err
					}
					return b.Put([]byte("profile"), []byte("Florian"))
				})
				require.Nil(t, err)
				return errors.New("i am an error returned by a migration")
			},
		},
	}

	require.Nil(t, db.Close())
	require.EqualError(t, Migrate(dbPath, migrations), "i am an error returned by a migration")

	db, err = bolt.Open(dbPath, dbFileMode, &bolt.Options{Timeout: time.Second})
	require.Nil(t, err)

	err = db.View(func(tx *bolt.Tx) error {
		// b should be nil since we don't accept changes made by
		// migrations when they return an error
		b := tx.Bucket([]byte("key_value_store"))
		require.Nil(t, b)

		// database_version should still be 0 since the migration failed
		systemBucket := tx.Bucket([]byte("system"))
		require.Equal(t, uint32(0), binary.BigEndian.Uint32(systemBucket.Get([]byte("database_version"))))
		return nil
	})
	require.Nil(t, err)
}

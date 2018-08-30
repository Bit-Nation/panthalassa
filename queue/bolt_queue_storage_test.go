package queue

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	storm "github.com/asdine/storm"
	require "github.com/stretchr/testify/require"
)

func createStorm() *storm.DB {
	dbPath, err := filepath.Abs(os.TempDir() + "/" + time.Now().String())
	if err != nil {
		panic(err)
	}
	db, err := storm.Open(dbPath)
	if err != nil {
		panic(err)
	}
	return db
}

func TestBoltQueueStorage_PersistAndDeleteJob(t *testing.T) {

	boltDB := createStorm()
	s := NewStorage(boltDB)

	jobToPersist := Job{
		Type: "SEND_MONEY",
		Data: map[string]interface{}{
			"key": "value",
		},
	}

	// persist job
	err := s.PersistJob(&jobToPersist)
	require.Nil(t, err)
	var jobs []Job
	require.Nil(t, s.db.All(&jobs))
	require.Equal(t, 1, len(jobs))

	// delete job
	require.Nil(t, s.DeleteJob(jobs[0]))
	require.Nil(t, s.db.All(&jobs))
	require.Equal(t, 0, len(jobs))

}

func TestBoltQueueStorage_Map(t *testing.T) {

	st := createStorm()
	s := NewStorage(st)

	// persist first job
	firstJob := Job{
		ID:   1,
		Type: "SEND_MONEY",
		Data: map[string]interface{}{
			"key": "value",
		},
	}
	err := s.PersistJob(&firstJob)
	require.Nil(t, err)

	// persist second job
	secondJob := Job{
		ID:   2,
		Type: "SEND_MONEY",
		Data: map[string]interface{}{
			"key": "value",
		},
	}
	err = s.PersistJob(&secondJob)
	require.Nil(t, err)

	stack := make(chan Job, 2)

	s.Map(stack)

	// make sure first job is ok
	job := <-stack
	require.Equal(t, firstJob, job)

	// make sure second job is ok
	job = <-stack
	require.Equal(t, secondJob, job)

}

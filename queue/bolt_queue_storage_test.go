package queue

import (
	"encoding/json"
	"errors"

	"github.com/coreos/bbolt"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestBoltQueueStorage_PersistAndDeleteJob(t *testing.T) {

	boltDB := createDB()

	s := NewStorage(boltDB)

	jobToPersist := Job{
		ID:   "my_id",
		Type: "SEND_MONEY",
		Data: map[string]interface{}{
			"key": "value",
		},
	}

	// persist job
	err := s.PersistJob(jobToPersist)
	require.Nil(t, err)

	// check if job exist
	err = s.db.View(func(tx *bolt.Tx) error {

		// queue bucket
		jobBucket := tx.Bucket(queueStorageBucketName)
		if err != nil {
			return err
		}

		rawJob := jobBucket.Get([]byte("my_id"))
		if rawJob == nil {
			return errors.New("failed to fetch job")
		}

		job := Job{}
		require.Nil(t, json.Unmarshal(rawJob, &job))

		require.Equal(t, jobToPersist, job)

		return nil

	})
	require.Nil(t, err)

	// delete job
	require.Nil(t, s.DeleteJob("my_id"))
	err = s.db.View(func(tx *bolt.Tx) error {

		// queue bucket
		jobBucket := tx.Bucket(queueStorageBucketName)
		if err != nil {
			return err
		}

		rawJob := jobBucket.Get([]byte("my_id"))
		require.Nil(t, rawJob)

		return nil

	})
	require.Nil(t, err)

}

func TestBoltQueueStorage_Map(t *testing.T) {

	boltDB := createDB()

	s := NewStorage(boltDB)

	// persist first job
	firstJob := Job{
		ID:   "job_id_first",
		Type: "SEND_MONEY",
		Data: map[string]interface{}{
			"key": "value",
		},
	}
	err := s.PersistJob(firstJob)
	require.Nil(t, err)

	// persist second job
	secondJob := Job{
		ID:   "job_id_second",
		Type: "SEND_MONEY",
		Data: map[string]interface{}{
			"key": "value",
		},
	}
	err = s.PersistJob(secondJob)
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

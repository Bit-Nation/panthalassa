package queue

import (
	"bytes"
	"encoding/json"

	bolt "github.com/coreos/bbolt"
)

var queueStorageBucketName = []byte("queue_storage")

type BoltQueueStorage struct {
	db *bolt.DB
}

func (s *BoltQueueStorage) PersistJob(j Job) error {
	return s.db.Update(func(tx *bolt.Tx) error {

		// queue bucket
		jobBucket, err := tx.CreateBucketIfNotExists(queueStorageBucketName)
		if err != nil {
			return err
		}

		// marshal job
		rawJob, err := json.Marshal(j)
		if err != nil {
			return err
		}

		return jobBucket.Put([]byte(j.ID), rawJob)

	})
}

func (s *BoltQueueStorage) DeleteJob(id string) error {
	return s.db.Update(func(tx *bolt.Tx) error {

		// queue bucket
		jobBucket, err := tx.CreateBucketIfNotExists(queueStorageBucketName)
		if err != nil {
			return err
		}

		return jobBucket.Delete([]byte(id))

	})
}

func (s *BoltQueueStorage) Map(queue chan Job) {
	go func() {

		err := s.db.View(func(tx *bolt.Tx) error {

			// queue bucket
			jobBucket := tx.Bucket(queueStorageBucketName)
			if jobBucket == nil {
				return nil
			}

			// map over job bucket
			return jobBucket.ForEach(func(_, job []byte) error {
				j := Job{}
				d := json.NewDecoder(bytes.NewReader(job))
				d.UseNumber()
				if err := d.Decode(&j); err != nil {
					logger.Error(err)
					return nil
				}
				queue <- j
				return nil
			})

		})

		if err != nil {
			logger.Error(err)
		}

	}()
}

func NewStorage(db *bolt.DB) *BoltQueueStorage {
	return &BoltQueueStorage{
		db: db,
	}
}

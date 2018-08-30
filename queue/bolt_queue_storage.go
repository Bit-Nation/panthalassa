package queue

import (
	storm "github.com/asdine/storm"
)

type BoltQueueStorage struct {
	db *storm.DB
}

func (s *BoltQueueStorage) PersistJob(j *Job) error {
	return s.db.Save(j)
}

func (s *BoltQueueStorage) DeleteJob(j Job) error {
	return s.db.DeleteStruct(&j)
}

// map over all jobs and send them into the given
func (s *BoltQueueStorage) Map(queue chan Job) {
	go func() {

		q := s.db.Select()
		err := q.Each(new(Job), func(i interface{}) error {
			j := i.(*Job)
			queue <- *j
			return nil
		})

		if err != nil {
			logger.Error(err)
		}

	}()
}

func NewStorage(db *storm.DB) *BoltQueueStorage {
	return &BoltQueueStorage{
		db: db,
	}
}

package queue

import (
	"os"
	"path/filepath"
	"time"

	bolt "github.com/coreos/bbolt"
)

type testProcessor struct {
	processorType string
	validJob      func(j Job) error
	process       func(j Job) error
}

func (p *testProcessor) Type() string {
	return p.processorType
}

func (p *testProcessor) ValidJob(j Job) error {
	return p.validJob(j)
}

func (p *testProcessor) Process(j Job) error {
	return p.process(j)
}

type testStorage struct {
	persistJob func(j Job) error
	deleteJob  func(j string) error
	mapFunc    func(queue chan Job)
}

func (s *testStorage) PersistJob(j Job) error {
	return s.persistJob(j)
}

func (s *testStorage) DeleteJob(id string) error {
	return s.deleteJob(id)
}

func (s *testStorage) Map(queue chan Job) {
	s.mapFunc(queue)
}

func createDB() *bolt.DB {
	dbPath, err := filepath.Abs(os.TempDir() + "/" + time.Now().String())
	if err != nil {
		panic(err)
	}
	db, err := bolt.Open(dbPath, 0600, &bolt.Options{Timeout: time.Second})
	if err != nil {
		panic(err)
	}
	return db
}

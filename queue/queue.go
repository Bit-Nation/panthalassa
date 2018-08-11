package queue

import (
	"fmt"
	"sync"

	log "github.com/ipfs/go-log"
)

// this is a very simple queue
// the purpose is to process a job till it succeed

var logger = log.Logger("queue")

type Processor interface {
	Type() string
	ValidJob(j Job) error
	Process(j Job) error
}

type Storage interface {
	PersistJob(j Job) error
	DeleteJob(id string) error
	Map(queue chan Job)
}

type Job struct {
	ID   string                 `json:"id"`
	Type string                 `json:"type"`
	Data map[string]interface{} `json:"data"`
}

type Queue struct {
	processors map[string]Processor
	storage    Storage
	lock       sync.Mutex
	jobStack   chan Job
}

// close the queue
func (q *Queue) Close() error {
	close(q.jobStack)
	return nil
}

// register a new processor
func (q *Queue) RegisterProcessor(p Processor) error {
	q.lock.Lock()
	defer q.lock.Unlock()
	if _, exist := q.processors[p.Type()]; exist {
		return fmt.Errorf("processor %s has already been registered", p.Type())
	}
	q.processors[p.Type()] = p
	return nil
}

// persist job to queue
func (q *Queue) AddJob(j Job) error {
	// lock
	q.lock.Lock()
	defer q.lock.Unlock()
	// fetch processor
	p, exist := q.processors[j.Type]
	if !exist {
		return fmt.Errorf("can't add job of type: %s since no processor has been registered", j.Type)
	}
	// validate job
	err := p.ValidJob(j)
	if err != nil {
		return err
	}
	// persist job
	if err := q.storage.PersistJob(j); err != nil {
		return err
	}
	// add job to stack so that queue will pick it up
	q.jobStack <- j
	return nil
}

func (q *Queue) fetchProcessor(processor string) (Processor, error) {
	// lock
	q.lock.Lock()
	q.lock.Unlock()
	// fetch processor
	p, exist := q.processors[processor]
	if !exist {
		return nil, fmt.Errorf("processor: %s does not exist", processor)
	}
	return p, nil
}

func (q *Queue) DeleteJob(j Job) error {
	return q.storage.DeleteJob(j.ID)
}

func New(s Storage, jobStackSize uint, concurrency uint) *Queue {

	// construct queue
	q := &Queue{
		processors: map[string]Processor{},
		storage:    s,
		lock:       sync.Mutex{},
		jobStack:   make(chan Job, jobStackSize),
	}

	// register all processors
	for {
		// exit when we registered all handlers
		if concurrency == 0 {
			break
		}
		concurrency--
		go func(q *Queue) {
			for {
				// exit if job stack go closed
				if q.jobStack == nil {
					return
				}
				select {
				case j := <-q.jobStack:
					// fetch processor
					p, err := q.fetchProcessor(j.Type)
					if err != nil {
						logger.Error(err)
						continue
					}
					// process error
					if err := p.Process(j); err != nil {
						logger.Error(err)
					}
				}
			}
		}(q)
	}

	// load past job's and add them to job stack
	go func() {
		s.Map(q.jobStack)
	}()

	return q
}

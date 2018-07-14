package request_limitation

import (
	"sync"
	"time"
)

type Throttling struct {
	concurrency    uint
	coolDown       time.Duration
	stack          chan func()
	lock           sync.Mutex
	maxQueue       uint
	inWork         uint
	queueFullError error
}

// create new throttling request limitation
func NewThrottling(concurrency uint, coolDown time.Duration, maxQueue uint, queueFullError error) *Throttling {

	t := &Throttling{
		concurrency:    concurrency,
		coolDown:       coolDown,
		stack:          make(chan func(), maxQueue),
		lock:           sync.Mutex{},
		maxQueue:       maxQueue,
		queueFullError: queueFullError,
	}

	go func() {
		for {
			t.lock.Lock()
			if t.inWork >= t.concurrency {
				t.lock.Unlock()
				continue
			}
			t.lock.Unlock()
			select {
			case job := <-t.stack:
				t.lock.Lock()
				t.inWork++
				t.lock.Unlock()
				go job()
				go func() {
					<-time.After(t.coolDown)
					t.lock.Lock()
					t.inWork--
					t.lock.Unlock()
				}()
			}
		}
	}()

	return t
}

func (t *Throttling) Exec(cb func()) error {
	if len(t.stack) >= int(t.maxQueue) {
		return t.queueFullError
	}
	t.stack <- cb
	return nil
}

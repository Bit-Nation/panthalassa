package request_limitation

import (
	"sync"
	"time"
)

// count throttling provide the combined functionality of
// throttling and count. It will have a count of current jobs
// that must be decreased AND the throttling based on time
type CountThrottling struct {
	concurrency    uint
	coolDown       time.Duration
	stack          chan func()
	lock           sync.Mutex
	maxQueue       uint
	inWork         uint
	queueFullError error
	current        uint
}

var countThrottlingSleep = time.Second

// create new throttling request limitation
func NewCountThrottling(concurrency uint, coolDown time.Duration, maxQueue uint, queueFullError error) *CountThrottling {

	t := &CountThrottling{
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
			if t.inWork >= t.concurrency || t.current >= t.concurrency {
				t.lock.Unlock()
				time.Sleep(countThrottlingSleep)
				continue
			}
			t.lock.Unlock()
			select {
			case job := <-t.stack:
				t.lock.Lock()
				t.inWork++
				t.current++
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

func (t *CountThrottling) Decrease() {
	t.lock.Lock()
	defer t.lock.Unlock()
	if t.current == 0 {
		return
	}
	t.current--
}

func (t *CountThrottling) Exec(cb func()) error {
	if len(t.stack) >= int(t.maxQueue) {
		return t.queueFullError
	}
	t.stack <- cb
	return nil
}

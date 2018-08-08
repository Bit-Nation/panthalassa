package request_limitation

import (
	"time"
)

// count throttling provide the combined functionality of
// throttling and count. It will have a count of current jobs
// that must be decreased AND the throttling based on time
type CountThrottling struct {
	concurrency uint
	coolDown    time.Duration
	// the channel should be called when the function succeed
	stack          chan func(chan struct{})
	maxQueue       uint
	queueFullError error
	currentChan    chan chan uint
	inWorkChan     chan chan uint
	decCurrent     chan struct{}
}

var countThrottlingSleep = time.Second

// create new throttling request limitation
func NewCountThrottling(concurrency uint, coolDown time.Duration, maxQueue uint, queueFullError error) *CountThrottling {

	t := &CountThrottling{
		concurrency:    concurrency,
		coolDown:       coolDown,
		stack:          make(chan func(chan struct{}), maxQueue),
		maxQueue:       maxQueue,
		queueFullError: queueFullError,
		currentChan:    make(chan chan uint),
		inWorkChan:     make(chan chan uint),
		decCurrent:     make(chan struct{}),
	}

	// "in work" related channels
	incInWork := make(chan struct{})
	decInWork := make(chan struct{})

	// "current" related channels
	incCurrent := make(chan struct{})

	// state
	go func() {

		var inWork uint
		var current uint

		for {

			// exit when stack got closed
			if t.stack == nil {
				return
			}

			select {

			case <-incInWork:
				inWork++
			case <-decInWork:
				if inWork > 0 {
					inWork--
				}
			case <-incCurrent:
				current++
			case <-t.decCurrent:
				if current > 0 {
					current--
				}
			case resp := <-t.inWorkChan:
				resp <- inWork
			case resp := <-t.currentChan:
				resp <- current
			}

		}

	}()

	// worker
	go func() {

		for {

			// exit when stack got closed
			if t.stack == nil {
				return
			}

			if t.inWork() >= t.concurrency || t.Current() >= t.concurrency {
				time.Sleep(countThrottlingSleep)
				continue
			}

			select {
			case job := <-t.stack:
				incInWork <- struct{}{}
				incCurrent <- struct{}{}
				go job(t.decCurrent)
				go func() {
					<-time.After(t.coolDown)
					decInWork <- struct{}{}
				}()
			}
		}
	}()

	return t
}

func (t *CountThrottling) Close() error {
	close(t.stack)
	return nil
}

func (t *CountThrottling) inWork() uint {
	iw := make(chan uint)
	t.inWorkChan <- iw
	return <-iw
}

func (t *CountThrottling) Decrease() {
	t.decCurrent <- struct{}{}
}

func (t *CountThrottling) Current() uint {
	cq := make(chan uint)
	t.currentChan <- cq
	return <-cq
}

func (t *CountThrottling) Exec(cb func(dec chan struct{})) error {
	if len(t.stack) >= int(t.maxQueue) {
		return t.queueFullError
	}
	t.stack <- cb
	return nil
}

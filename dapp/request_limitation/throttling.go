package request_limitation

import (
	"time"
)

type Throttling struct {
	concurrency    uint
	coolDown       time.Duration
	stack          chan func(chan struct{})
	maxQueue       uint
	queueFullError error
	vmDone         chan struct{}
}

func (t *Throttling) Close() error {
	close(t.stack)
	return nil
}

// create new throttling request limitation
func NewThrottling(concurrency uint, coolDown time.Duration, maxQueue uint, queueFullError error) *Throttling {

	t := &Throttling{
		concurrency:    concurrency,
		coolDown:       coolDown,
		stack:          make(chan func(chan struct{}), maxQueue),
		maxQueue:       maxQueue,
		queueFullError: queueFullError,
		vmDone:         make(chan struct{}),
	}

	inWorkChan := make(chan chan uint)
	incInWork := make(chan struct{})
	decInWork := make(chan struct{})

	// state
	go func() {

		var inWork uint

		for {

			// exit when stack got closed
			if t.stack == nil {
				return
			}

			select {
			case respChan := <-inWorkChan:
				respChan <- inWork
			case <-incInWork:
				inWork++
			case <-decInWork:
				if inWork == 0 {
					continue
				}
				inWork--
			}

		}

	}()

	go func() {
		for {

			// exit when stack got closed
			if t.stack == nil {
				return
			}

			iw := make(chan uint)
			inWorkChan <- iw

			if <-iw >= t.concurrency {
				time.Sleep(time.Second)
				continue
			}

			select {
			case job := <-t.stack:
				incInWork <- struct{}{}
				go job(t.vmDone)
				go func() {
					<-time.After(t.coolDown)
					decInWork <- struct{}{}
				}()
			}
		}
	}()

	return t
}

func (t *Throttling) Exec(cb func(vmDone chan struct{})) error {
	if len(t.stack) >= int(t.maxQueue) {
		return t.queueFullError
	}
	t.stack <- cb
	t.vmDone <- struct{}{}
	return nil
}

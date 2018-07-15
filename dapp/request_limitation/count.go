package request_limitation

import "sync"

// a simple increase and decrease
// request limitation.
type Count struct {
	lock              sync.Mutex
	count             uint
	max               uint
	canNotIncreaseErr error
}

func NewCount(max uint, canNotIncreaseErr error) *Count {
	return &Count{
		lock:              sync.Mutex{},
		max:               max,
		count:             0,
		canNotIncreaseErr: canNotIncreaseErr,
	}
}

// get current count of request limitation
func (c *Count) Count() uint {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.count
}

func (c *Count) Increase() error {
	c.lock.Lock()
	defer c.lock.Unlock()
	if c.count+1 > c.max {
		return c.canNotIncreaseErr
	}
	c.count++
	return nil
}

func (c *Count) Decrease() {
	c.lock.Lock()
	defer c.lock.Unlock()
	if c.count == 0 {
		return
	}
	c.count--
}

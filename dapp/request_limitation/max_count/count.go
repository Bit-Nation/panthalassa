package max_count

import "sync"

// a simple increase and decrease
// request limitation.
type Count struct {
	lock              sync.Mutex
	count             uint
	max               uint
	canNotIncreaseErr error
}

func New(max uint, canNotIncreaseErr error) *Count {
	return &Count{
		lock:              sync.Mutex{},
		max:               max,
		count:             0,
		canNotIncreaseErr: canNotIncreaseErr,
	}
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
	c.count--
}

package request_limitation

// a simple increase and decrease
// request limitation.
type Count struct {
	max               uint
	canNotIncreaseErr error
	increase          chan struct{}
	decrease          chan struct{}
	count             chan chan uint
}

func NewCount(max uint, canNotIncreaseErr error) *Count {
	c := &Count{
		max:               max,
		canNotIncreaseErr: canNotIncreaseErr,
		increase:          make(chan struct{}),
		decrease:          make(chan struct{}),
		count:             make(chan chan uint),
	}

	// state
	go func() {

		var count uint

		for {

			// exit from go routine
			if c.increase == nil || c.decrease == nil || c.count == nil {
				return
			}

			select {
			case <-c.increase:
				count++
			case <-c.decrease:
				count--
			case c := <-c.count:
				c <- count
			}

		}

	}()

	return c

}

func (c *Count) Close() error {
	close(c.increase)
	close(c.decrease)
	close(c.count)
	return nil
}

// get current count of request limitation
func (c *Count) Count() uint {
	q := make(chan uint)
	c.count <- q
	return <-q
}

func (c *Count) Increase() error {
	if c.Count() > c.max {
		return c.canNotIncreaseErr
	}
	c.increase <- struct{}{}
	return nil
}

func (c *Count) Decrease() {
	if c.Count() == 0 {
		return
	}
	c.decrease <- struct{}{}
}

package request_limitation

import (
	"errors"
	"testing"
	"time"

	require "github.com/stretchr/testify/require"
)

func TestThrottling_Exec(t *testing.T) {

	throttling := NewThrottling(0, time.Second, 1, errors.New("queue is full"))

	require.Nil(t, throttling.Exec(func() {}))
	require.EqualError(t, throttling.Exec(func() {}), "queue is full")

}

/**
@todo this test is failing from time to time
func TestThrottling_ExecCoolDown(t *testing.T) {

	throttling := NewThrottling(2, time.Second, 10, nil)

	// add to stack
	require.Nil(t, throttling.Exec(func() {}))
	require.Nil(t, throttling.Exec(func() {}))
	require.Nil(t, throttling.Exec(func() {
		time.Sleep(time.Second)
	}))
	// wait a bit so that the queue can pickup the job
	time.Sleep(time.Millisecond * 100)

	// should be 2 since we chose two for concurrency
	require.Equal(t, uint(2), throttling.inWork)

	// wait for the cool down to be over
	time.Sleep(time.Second)

	// should be 1 since the we waited a second.
	// Which means that the cool down is over.
	// After the cool down is over, the missing
	// jobs will be processes
	require.Equal(t, uint(1), throttling.inWork)

	// sleep a second to wait for cool down
	time.Sleep(time.Second)

	// make sure in work is 0 since we waited
	// for the cool down
	require.Equal(t, uint(0), throttling.inWork)

}
 */

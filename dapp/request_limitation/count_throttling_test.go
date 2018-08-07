package request_limitation

import (
	"errors"
	"testing"
	"time"

	require "github.com/stretchr/testify/require"
)

func TestCountThrottling(t *testing.T) {

	ct := NewCountThrottling(
		3,
		time.Second*1,
		10,
		errors.New("queue full error"),
	)

	// override sleep
	countThrottlingSleep = time.Microsecond

	// decrease channel
	var decChan chan struct{}

	require.Nil(t, ct.Exec(func(dec chan struct{}) {
		decChan = dec
	}))
	require.Nil(t, ct.Exec(func(dec chan struct{}) {

	}))
	require.Nil(t, ct.Exec(func(dec chan struct{}) {

	}))
	require.Nil(t, ct.Exec(func(dec chan struct{}) {

	}))
	require.Nil(t, ct.Exec(func(dec chan struct{}) {

	}))
	require.Nil(t, ct.Exec(func(dec chan struct{}) {

	}))
	require.Nil(t, ct.Exec(func(dec chan struct{}) {

	}))
	require.Nil(t, ct.Exec(func(dec chan struct{}) {

	}))
	require.Nil(t, ct.Exec(func(dec chan struct{}) {

	}))
	require.Nil(t, ct.Exec(func(dec chan struct{}) {

	}))

	// wait for queue to pick up the jobs
	time.Sleep(time.Millisecond * 10)
	require.Equal(t, uint(3), ct.inWork())
	require.Equal(t, uint(3), ct.Current())

	// wait a second to make sure throttling is over
	time.Sleep(time.Second)

	require.Equal(t, uint(0), ct.inWork())
	// current must be 3 since we didn't call Decrease
	// to decrease current. As long as current is
	// greater or equal then the concurrency,
	// the worker can't pick up new jobs
	require.Equal(t, uint(3), ct.Current())

	// decrease the amount of current
	// to make sure worker will pick up new job
	decChan <- struct{}{}

	// wait for queue to pick up new jobs
	time.Sleep(time.Millisecond * 10)
	require.Equal(t, uint(1), ct.inWork())

}

func TestCountThrottlingFullError(t *testing.T) {

	ct := NewCountThrottling(
		0,
		time.Second*1,
		5,
		errors.New("queue full error"),
	)

	require.Nil(t, ct.Exec(func(dec chan struct{}) {
		dec <- struct{}{}
	}))
	require.Nil(t, ct.Exec(func(dec chan struct{}) {
		dec <- struct{}{}
	}))
	require.Nil(t, ct.Exec(func(dec chan struct{}) {
		dec <- struct{}{}
	}))
	require.Nil(t, ct.Exec(func(dec chan struct{}) {
		dec <- struct{}{}
	}))
	require.Nil(t, ct.Exec(func(dec chan struct{}) {
		dec <- struct{}{}
	}))

	require.EqualError(t, ct.Exec(func(dec chan struct{}) {}), "queue full error")

}

/**
func TestCountThrottling_Decrease(t *testing.T) {

	ct := NewCountThrottling(
		1,
		time.Second*1,
		5,
		errors.New("queue full error"),
	)

	ct.current = 1

	// must decrease current by one
	ct.Decrease()
	require.Equal(t, uint(0), ct.current)

	// if current is zero, calling Decrease shouldn't decrease by one
	ct.Decrease()
	require.Equal(t, uint(0), ct.current)

}

*/

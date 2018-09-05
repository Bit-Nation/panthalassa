package queue

import (
	"errors"
	"testing"

	require "github.com/stretchr/testify/require"
)

func TestRegisterProcessorError(t *testing.T) {

	st := createStorm()
	s := NewStorage(st)
	queue := New(s, 10, 3)

	// register the first time should be valid
	err := queue.RegisterProcessor(&testProcessor{
		processorType: "MY_PROCESSOR",
		validJob: func(j Job) error {
			return nil
		},
		process: func(j Job) error {
			return nil
		},
	})
	require.Nil(t, err)

	// register second time should be invalid
	err = queue.RegisterProcessor(&testProcessor{
		processorType: "MY_PROCESSOR",
		validJob: func(j Job) error {
			return nil
		},
		process: func(j Job) error {
			return nil
		},
	})
	require.EqualError(t, err, "processor MY_PROCESSOR has already been registered")

}

func TestFetchProcessor(t *testing.T) {

	st := createStorm()
	s := NewStorage(st)
	queue := New(s, 10, 3)

	_, err := queue.fetchProcessor("I_DO_NOT_EXIST")
	require.EqualError(t, err, "processor: I_DO_NOT_EXIST does not exist")

	// register the first time should be valid
	err = queue.RegisterProcessor(&testProcessor{
		processorType: "MY_PROCESSOR",
		validJob: func(j Job) error {
			return nil
		},
		process: func(j Job) error {
			return nil
		},
	})
	require.Nil(t, err)

	// fetch processor should succeed
	p, err := queue.fetchProcessor("MY_PROCESSOR")
	require.Nil(t, err)
	require.NotNil(t, p)

}

func TestQueue_AddJobError(t *testing.T) {

	st := createStorm()
	s := NewStorage(st)

	queue := New(s, 10, 3)

	// add job
	err := queue.AddJob(Job{
		Type: "MY_JOB",
		Data: map[string]interface{}{
			"key": "value",
		},
	})
	require.EqualError(t, err, "can't add job of type: MY_JOB since no processor has been registered")

}

func TestQueue_AddInvalidJob(t *testing.T) {

	st := createStorm()
	s := NewStorage(st)
	queue := New(s, 10, 3)

	// register processor
	err := queue.RegisterProcessor(&testProcessor{
		processorType: "MY_JOB",
		validJob: func(j Job) error {
			return errors.New("got invalid job")
		},
		process: func(j Job) error {
			return nil
		},
	})
	require.Nil(t, err)

	// add job
	err = queue.AddJob(Job{
		Type: "MY_JOB",
		Data: map[string]interface{}{
			"key": "value",
		},
	})
	require.EqualError(t, err, "got invalid job")

}

func TestQueue_AddJob(t *testing.T) {

	db := &BoltQueueStorage{db: createStorm()}

	queue := New(db, 10, 3)

	// register processor
	err := queue.RegisterProcessor(&testProcessor{
		processorType: "MY_JOB",
		validJob: func(j Job) error {
			return nil
		},
		process: func(j Job) error {
			return nil
		},
	})
	require.Nil(t, err)

	// add job
	err = queue.AddJob(Job{
		Type: "MY_JOB",
		Data: map[string]interface{}{
			"key": "value",
		},
	})
	require.Nil(t, err)

	amount, err := db.db.Count(&Job{})
	require.Nil(t, err)
	require.Equal(t, 1, amount)

}

//
func TestQueue_ProcessJob(t *testing.T) {

	db := &BoltQueueStorage{db: createStorm()}

	queue := New(db, 10, 3)

	// channel
	wait := make(chan struct{}, 1)

	// register test processor
	err := queue.RegisterProcessor(&testProcessor{
		processorType: "SEND_MONEY",
		validJob: func(j Job) error {
			return nil
		},
		process: func(j Job) error {
			if j.Type != "SEND_MONEY" {
				panic("got invalid type")
			}
			wait <- struct{}{}
			return nil
		},
	})

	require.Nil(t, err)

	// add job to stack
	queue.jobStack <- Job{
		Type: "SEND_MONEY",
	}

	// wait till the job has been processed
	<-wait

}

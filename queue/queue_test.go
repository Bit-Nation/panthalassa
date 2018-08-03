package queue

import (
	"errors"
	"testing"

	require "github.com/stretchr/testify/require"
)

func TestRegisterProcessorError(t *testing.T) {

	s := testStorage{}

	queue := New(&s, 10, 3)

	// register the first time should be valid
	err := queue.RegisterProcessor(&testProcessor{
		processorType: "MY_PROCESSOR",
	})
	require.Nil(t, err)

	// register second time should be invalid
	err = queue.RegisterProcessor(&testProcessor{
		processorType: "MY_PROCESSOR",
	})
	require.EqualError(t, err, "processor MY_PROCESSOR has already been registered")

}

func TestFetchProcessor(t *testing.T) {

	s := testStorage{}

	queue := New(&s, 10, 3)

	_, err := queue.fetchProcessor("I_DO_NOT_EXIST")
	require.EqualError(t, err, "processor: I_DO_NOT_EXIST does not exist")

	// register the first time should be valid
	err = queue.RegisterProcessor(&testProcessor{
		processorType: "MY_PROCESSOR",
	})
	require.Nil(t, err)

	// fetch processor should succeed
	p, err := queue.fetchProcessor("MY_PROCESSOR")
	require.Nil(t, err)
	require.NotNil(t, p)

}

func TestQueue_AddJobError(t *testing.T) {

	s := testStorage{}

	queue := New(&s, 10, 3)

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

	s := testStorage{}

	queue := New(&s, 10, 3)

	// register processor
	err := queue.RegisterProcessor(&testProcessor{
		processorType: "MY_JOB",
		validJob: func(j Job) error {
			return errors.New("got invalid job")
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

	calledPersistJob := false
	s := testStorage{
		persistJob: func(j Job) error {
			require.Equal(t, "MY_JOB", j.Type)
			require.Equal(t, map[string]interface{}{"key": "value"}, j.Data)
			calledPersistJob = true
			return nil
		},
	}

	queue := New(&s, 10, 3)

	// register processor
	err := queue.RegisterProcessor(&testProcessor{
		processorType: "MY_JOB",
		validJob: func(j Job) error {
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
	require.True(t, calledPersistJob)

}

//
func TestQueue_ProcessJob(t *testing.T) {

	queue := New(&testStorage{}, 10, 3)

	// channel
	wait := make(chan struct{}, 1)

	// register test processor
	err := queue.RegisterProcessor(&testProcessor{
		processorType: "SEND_MONEY",
		validJob: func(j Job) error {
			return nil
		},
		process: func(j Job) error {
			require.Equal(t, "<job-id>", j.ID)
			require.Equal(t, "SEND_MONEY", j.Type)
			wait <- struct{}{}
			return nil
		},
	})

	require.Nil(t, err)

	// add job to stack
	queue.jobStack <- Job{
		ID:   "<job-id>",
		Type: "SEND_MONEY",
	}

	// wait till the job has been processed
	<-wait

}

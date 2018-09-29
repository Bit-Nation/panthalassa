package chat

import (
	"testing"

	queue "github.com/Bit-Nation/panthalassa/queue"
	require "github.com/stretchr/testify/require"
)

func TestSubmitMessagesProcessor_ValidJob(t *testing.T) {

	p := SubmitMessagesProcessor{}
	err := p.ValidJob(queue.Job{
		Type: "MESSAGE:SUBMIT",
		Data: map[string]interface{}{
			"db_message_id": int64(3),
			"chat_id":       3,
		},
	})
	require.Nil(t, err)

}

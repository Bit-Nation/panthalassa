package chat

import (
	"testing"
	"crypto/rand"
	
	queue "github.com/Bit-Nation/panthalassa/queue"
	require "github.com/stretchr/testify/require"
	ed25519 "golang.org/x/crypto/ed25519"
)

func TestSubmitMessagesProcessor_ValidJob(t *testing.T) {
	
	pub, _, err := ed25519.GenerateKey(rand.Reader)
	require.Nil(t, err)
	
	p := SubmitMessagesProcessor{}
	err = p.ValidJob(queue.Job{
		Type: "MESSAGE:SUBMIT",
		Data: map[string]interface{}{
			"db_message_id": int64(3),
			"partner": pub,
		},
	})
	require.Nil(t, err)
	
}
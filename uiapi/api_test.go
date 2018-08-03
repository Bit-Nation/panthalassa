package stapi

import (
	"testing"
	"time"
	
	require "github.com/stretchr/testify/require"
)

type upstream struct {
	send func(data string)
}

func (u *upstream) Send(data string) {
	u.send(data)
}

func TestApi_Send(t *testing.T) {
	
	signal := make(chan string, 1)
	up := upstream{
		send: func(data string) {
			signal <- data
		},
	}
	
	a := New(&up)
	a.Send("TEST:CALL", map[string]interface{}{
		"key": "value",
	})
	
	select {
		case data := <-signal:
			require.Equal(t, `{"name":"TEST:CALL","payload":{"key":"value"}}`, data)
		case <- time.After(time.Second):
			require.Fail(t, "time out")
	}
	
}
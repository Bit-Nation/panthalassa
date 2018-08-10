package backend

import (
	"net/http"
	"testing"
	"time"

	bpb "github.com/Bit-Nation/protobuffers"
	proto "github.com/gogo/protobuf/proto"
	mux "github.com/gorilla/mux"
	gws "github.com/gorilla/websocket"
	require "github.com/stretchr/testify/require"
)

func TestWSTransport_Send(t *testing.T) {

	// setup test websocket server
	router := mux.Router{}
	server := &http.Server{Addr: ":3857", Handler: &router}
	defer server.Close()
	upgrader := gws.Upgrader{}
	reader := make(chan []byte, 1)
	router.HandleFunc("/ws", func(writer http.ResponseWriter, request *http.Request) {
		// connection upgrade
		conn, err := upgrader.Upgrade(writer, request, nil)
		if err != nil {
			panic(err)
		}
		// read message
		_, msg, err := conn.ReadMessage()
		if err != nil {
			panic(err)
		}
		// send message to our reader channel
		reader <- msg
		if err := conn.Close(); err != nil {
			panic(err)
		}
	})

	// start websocket server
	go func() {
		server.ListenAndServe()
	}()

	// setup new transport
	trans := NewWSTransport("ws://127.0.0.1:3857/ws", "")

	// send test message to transport
	err := trans.Send(&bpb.BackendMessage{
		RequestID: "request-id",
	})
	require.Nil(t, err)

	// assertion
	select {
	case rawProtoMsg := <-reader:
		// make sure the message we got has the same request id
		msg := bpb.BackendMessage{}
		require.Nil(t, proto.Unmarshal(rawProtoMsg, &msg))
		require.Equal(t, "request-id", msg.RequestID)
	case <-time.After(time.Second * 2):
		require.Fail(t, "timed out")
	}

}

func TestWSTransport_NextMessage(t *testing.T) {

	// setup test websocket server
	router := mux.Router{}
	server := &http.Server{Addr: ":3857", Handler: &router}
	defer server.Close()
	upgrader := gws.Upgrader{}
	router.HandleFunc("/ws", func(writer http.ResponseWriter, request *http.Request) {
		// connection upgrade
		conn, err := upgrader.Upgrade(writer, request, nil)
		if err != nil {
			panic(err)
		}
		// send message to client
		rawMsg, err := proto.Marshal(&bpb.BackendMessage{
			RequestID: "request-id",
		})
		if err != nil {
			panic(err)
		}
		conn.WriteMessage(gws.BinaryMessage, rawMsg)
		conn.Close()
	})

	// start websocket server
	go func() {
		server.ListenAndServe()
	}()

	trans := NewWSTransport("ws://127.0.0.1:3857/ws", "")

	msg := trans.NextMessage()
	require.Equal(t, "request-id", msg.RequestID)

}

package backend

import (
	"net/http"
	"testing"
	"time"

	keyManager "github.com/Bit-Nation/panthalassa/keyManager"
	keyStore "github.com/Bit-Nation/panthalassa/keyStore"
	mnemonic "github.com/Bit-Nation/panthalassa/mnemonic"
	bpb "github.com/Bit-Nation/protobuffers"
	proto "github.com/gogo/protobuf/proto"
	mux "github.com/gorilla/mux"
	gws "github.com/gorilla/websocket"
	require "github.com/stretchr/testify/require"
)

func TestWSTransport_Send(t *testing.T) {

	// key manager setup
	mne, err := mnemonic.New()
	require.Nil(t, err)
	ks, err := keyStore.NewFromMnemonic(mne)
	require.Nil(t, err)
	km := keyManager.CreateFromKeyStore(ks)

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
	trans := NewWSTransport("ws://127.0.0.1:3857/ws", "", km)

	// send test message to transport
	err = trans.Send(&bpb.BackendMessage{
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

	// key manager setup
	mne, err := mnemonic.New()
	require.Nil(t, err)
	ks, err := keyStore.NewFromMnemonic(mne)
	require.Nil(t, err)
	km := keyManager.CreateFromKeyStore(ks)

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
		// @TODO Research more on .CloseNormalClosure and .CloseMessage for expected way to close connection
		// @TODO This should be done before calling conn.Close()
		// @TODO Client should be handling different connection closure types, including .CloseNormalClosure
		//closureMessage := gws.FormatCloseMessage(gws.CloseNormalClosure, "Closing Connection")
		//if err := conn.WriteMessage(gws.CloseMessage, cm); err != nil {
		//	panic(err)
		//}
		conn.Close()
	})

	// start websocket server
	go func() {
		server.ListenAndServe()
	}()

	trans := NewWSTransport("ws://127.0.0.1:3857/ws", "", km)

	msg, err := trans.NextMessage()
	require.Nil(t, err)
	require.Equal(t, "request-id", msg.RequestID)

}

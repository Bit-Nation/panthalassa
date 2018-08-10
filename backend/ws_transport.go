package backend

import (
	"fmt"
	"net/http"
	"time"

	bpb "github.com/Bit-Nation/protobuffers"
	proto "github.com/gogo/protobuf/proto"
	gws "github.com/gorilla/websocket"
	log "github.com/ipfs/go-log"
)

var wsTransLogger = log.Logger("ws transport")

type WSTransport struct {
	closer chan struct{}
	conn   *gws.Conn
	write  chan *bpb.BackendMessage
	read   chan *bpb.BackendMessage
}

func NewWSTransport(endpoint, bearerToken string) *WSTransport {

	// construct ws transport
	wst := &WSTransport{
		closer: make(chan struct{}, 2),
		write:  make(chan *bpb.BackendMessage, 100),
		read:   make(chan *bpb.BackendMessage, 100),
	}

	go func() {

		// dial to endpoint
		d := gws.Dialer{}

		// try to connect till success
		for {
			conn, _, err := d.Dial(endpoint, http.Header{
				"Bearer": []string{bearerToken},
			})
			if err != nil {
				wsTransLogger.Error(err)
				time.Sleep(time.Second)
				continue
			}
			wst.conn = conn
			break
		}

		wst.conn.SetCloseHandler(func(code int, text string) error {
			panic("closed")
			wsTransLogger.Warning("closed websocket, code: %d - message: %s", code, text)
			return nil
		})

		// start reader
		go func() {

			for {
				// exit when channel got closed
				if wst.read == nil {
					break
				}

				// react message
				mt, msg, err := wst.conn.ReadMessage()
				if err != nil {
					wsTransLogger.Error(err)
					continue
				}
				wsTransLogger.Debugf(
					`got message of type: %d - content: %s`,
					mt,
					string(msg),
				)

				// unmarshal message into protobuf
				m := &bpb.BackendMessage{}
				proto.Unmarshal(msg, m)
				if err != nil {
					wsTransLogger.Error(err)
					continue
				}

				// send to read channel so that it can be fetched from the NextMessage function
				wst.read <- m

			}

		}()

		// start writer
		go func() {

			for {
				// exit when channel got closed
				if wst.write == nil {
					break
				}

				select {
				case msg := <-wst.write:
					fmt.Println("write")
					wsTransLogger.Debugf(
						"going to write backend message with id: %s to ws",
						msg.RequestID,
					)
					rawMsg, err := proto.Marshal(msg)
					if err != nil {
						wsTransLogger.Error(err)
						continue
					}
					if err := wst.conn.WriteMessage(gws.BinaryMessage, rawMsg); err != nil {
						wsTransLogger.Error(err)
					}
				}
			}

		}()

	}()

	return wst

}

func (t *WSTransport) Send(msg *bpb.BackendMessage) error {
	t.write <- msg
	return nil

}

func (t *WSTransport) NextMessage() *bpb.BackendMessage {
	return <-t.read
}

func (t *WSTransport) Close() error {
	close(t.write)
	close(t.read)
	return t.conn.Close()
}

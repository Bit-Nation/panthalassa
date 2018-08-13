package backend

import (
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
	conn   *conn
	write  chan *bpb.BackendMessage
	read   chan *bpb.BackendMessage
}

// connection is kind of a extension of the gws.Conn
// it has additional state + some utils we need
type conn struct {
	closer chan struct{}
	wsConn *gws.Conn
}

func (c *conn) Close() error {
	c.closer <- struct{}{}
	return nil
}

func (t *WSTransport) newConn(closed chan struct{}, endpoint, bearerToken string) *conn {

	c := &conn{
		closer: make(chan struct{}, 2),
	}

	// ask this for the closed state
	isClosed := make(chan chan bool)

	// connection state routine
	go func() {
		var closed bool
		for {
			select {
			// query for closed
			case isClosedResp := <-isClosed:
				isClosedResp <- closed
			case <-c.closer:
				closed = true
			}
		}
	}()

	go func() {

		// dial to endpoint
		d := gws.Dialer{}

		// try to connect till success
		var wsConn *gws.Conn
		for {
			conn, _, err := d.Dial(endpoint, http.Header{
				"Bearer": []string{bearerToken},
			})
			if err != nil {
				wsTransLogger.Error(err)
				time.Sleep(time.Second)
				continue
			}
			wsConn = conn
			break
		}

		wsConn.SetCloseHandler(func(code int, text string) error {
			wsTransLogger.Warning("closed websocket, code: %d - message: %s", code, text)
			return nil
		})

		// start reader
		go func() {

			for {
				// exit when connection got closed
				isClosedRespChan := make(chan bool)
				isClosed <- isClosedRespChan
				if <-isClosedRespChan {
					logger.Debug("stop reading from websocket")
					break
				}

				// react message
				mt, msg, err := t.conn.wsConn.ReadMessage()
				if err != nil {
					wsTransLogger.Error(err)
					time.Sleep(5 * time.Second)
					closed <- struct{}{}
					break
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
				t.read <- m

			}

		}()

		// start writer
		go func() {

			for {

				// exit when connection got closed
				isClosedRespChan := make(chan bool)
				isClosed <- isClosedRespChan
				if <-isClosedRespChan {
					logger.Debug("stop writing to websocket")
					break
				}

				msg := <-t.write

				wsTransLogger.Debugf(
					"going to write backend message with id: %s to ws",
					msg.RequestID,
				)
				rawMsg, err := proto.Marshal(msg)
				if err != nil {
					wsTransLogger.Error(err)
					continue
				}
				if err := t.conn.wsConn.WriteMessage(gws.BinaryMessage, rawMsg); err != nil {
					wsTransLogger.Error(err)
				}

			}

		}()

	}()

	return c

}

func NewWSTransport(endpoint, bearerToken string) *WSTransport {

	// construct ws transport
	wst := &WSTransport{
		closer: make(chan struct{}),
		write:  make(chan *bpb.BackendMessage, 100),
		read:   make(chan *bpb.BackendMessage, 100),
	}

	// routine that keeps track of the connection
	// close and re connect
	connClosed := make(chan struct{}, 5)
	go func() {
		for {
			select {
			case <-wst.closer:
				break
			case <-connClosed:
				wst.conn = wst.newConn(connClosed, endpoint, bearerToken)
			}
		}
	}()

	// create initial connection
	go func() {
		wst.conn = wst.newConn(connClosed, endpoint, bearerToken)
	}()

	return wst

}

func (t *WSTransport) Send(msg *bpb.BackendMessage) error {
	t.write <- msg
	return nil
}

func (t *WSTransport) NextMessage() (*bpb.BackendMessage, error) {
	return <-t.read, nil
}

func (t *WSTransport) Close() error {
	t.closer <- struct{}{}
	if t.conn == nil {
		return nil
	}
	return t.conn.Close()
}

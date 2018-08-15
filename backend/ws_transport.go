package backend

import (
	"net/http"

	bpb "github.com/Bit-Nation/protobuffers"
	proto "github.com/gogo/protobuf/proto"
	log "github.com/ipfs/go-log"
	rws "github.com/mariuspass/recws"
)

var wsTransLogger = log.Logger("ws transport")

type WSTransport struct {
	write chan *bpb.BackendMessage
	read  chan *bpb.BackendMessage
	Conn  *rws.RecConn
}

func NewWSTransport(endpoint, bearerToken string) *WSTransport {

	// construct ws transport
	wst := &WSTransport{
		write: make(chan *bpb.BackendMessage, 100),
		read:  make(chan *bpb.BackendMessage, 100),
	}

	c := &rws.RecConn{}
	c.Dial(endpoint, http.Header{
		"Bearer": []string{bearerToken},
	})
	wst.Conn = c

	// writer
	go func() {
		for {
			msg := <-wst.write
			rawMsg, err := proto.Marshal(msg)
			if err != nil {
				logger.Error(err)
				continue
			}
			err = c.WriteMessage(2, rawMsg)
			if err == nil {
				continue
			}
			if err == rws.ErrNotConnected {
				continue
			}
			logger.Error(err)
			break
		}
	}()

	// reader
	go func() {
		for {
			_, rawMsg, err := c.ReadMessage()
			if err != nil {
				if err == rws.ErrNotConnected {
					continue
				}
				logger.Error(err)
				break
			}
			msg := &bpb.BackendMessage{}
			if err := proto.Unmarshal(rawMsg, msg); err != nil {
				logger.Error(err)
				continue
			}
			wst.read <- msg
		}
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

func (t *WSTransport) RegisterConnectionCloseListener(listener chan struct{}) {
	go func() {
		lastInformedAbout := true
		for {
			if !t.Conn.IsConnected() {
				if lastInformedAbout == false {
					continue
				}
				lastInformedAbout = false
				listener <- struct{}{}
				continue
			}
		}
	}()
	return
}

func (t *WSTransport) Close() error {
	return nil
}

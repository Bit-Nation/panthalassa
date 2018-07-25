package chat

import (
	"encoding/hex"

	db "github.com/Bit-Nation/panthalassa/db"
	stapi "github.com/Bit-Nation/panthalassa/uiapi"
	proto "github.com/gogo/protobuf/proto"
)

func NewMessageDBListener(api *stapi.Api, chanBuff uint32) chan db.PlainMessagePersistedEvent {

	listener := make(chan db.PlainMessagePersistedEvent, chanBuff)

	go func() {

		for {
			// exit if closed
			if listener == nil {
				return
			}
			select {
			case event := <-listener:

				if event.DBMsg.Status == db.StatusPersisted {
					rawMsg, err := proto.Marshal(&event.Msg)
					if err != nil {
						logger.Error(err)
					}
					api.Send("MESSAGE:PERSISTED", map[string]interface{}{
						"message_id": event.MsgID,
						"message":    string(rawMsg),
						"partner":    hex.EncodeToString(event.Partner),
					})
				}

				if event.DBMsg.Status == db.StatusDelivered {
					rawMsg, err := proto.Marshal(&event.Msg)
					if err != nil {
						logger.Error(err)
					}
					api.Send("MESSAGE:DELIVERED", map[string]interface{}{
						"message_id": event.MsgID,
						"message":    string(rawMsg),
						"partner":    hex.EncodeToString(event.Partner),
					})
				}

				if event.DBMsg.Received {
					rawMsg, err := proto.Marshal(&event.Msg)
					if err != nil {
						logger.Error(err)
					}
					api.Send("MESSAGE:RECEIVED", map[string]interface{}{
						"message_id": event.MsgID,
						"message":    string(rawMsg),
						"partner":    hex.EncodeToString(event.Partner),
					})
				}

			}
		}

	}()

	return listener

}

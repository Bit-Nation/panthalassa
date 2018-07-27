package chat

import (
	"encoding/hex"

	db "github.com/Bit-Nation/panthalassa/db"
	stapi "github.com/Bit-Nation/panthalassa/uiapi"
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
					api.Send("MESSAGE:PERSISTED", map[string]interface{}{
						"message_id": event.MsgID,
						"message": map[string]interface{}{
							"content":    string(event.Msg.Message),
							"created_at": event.Msg.CreatedAt,
						},
						"partner": hex.EncodeToString(event.Partner),
					})
				}

				if event.DBMsg.Status == db.StatusDelivered {
					api.Send("MESSAGE:DELIVERED", map[string]interface{}{
						"message_id": event.MsgID,
						"message": map[string]interface{}{
							"content":    string(event.Msg.Message),
							"created_at": event.Msg.CreatedAt,
						},
						"partner": hex.EncodeToString(event.Partner),
					})
				}

				if event.DBMsg.Received {
					api.Send("MESSAGE:RECEIVED", map[string]interface{}{
						"message_id": event.MsgID,
						"message": map[string]interface{}{
							"content":    string(event.Msg.Message),
							"created_at": event.Msg.CreatedAt,
						},
						"partner": hex.EncodeToString(event.Partner),
					})
				}

			}
		}

	}()

	return listener

}

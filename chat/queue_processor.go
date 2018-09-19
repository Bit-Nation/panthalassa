package chat

import (
	"encoding/hex"
	"errors"

	db "github.com/Bit-Nation/panthalassa/db"
	queue "github.com/Bit-Nation/panthalassa/queue"
)

// processor that submits messages from the queue to the backend
type SubmitMessagesProcessor struct {
	chat   *Chat
	chatDB db.ChatStorage
	queue  *queue.Queue
}

func (p *SubmitMessagesProcessor) Type() string {
	return "MESSAGE:SUBMIT"
}

func (p *SubmitMessagesProcessor) ValidJob(j queue.Job) error {

	if p.Type() != j.Type {
		return errors.New("invalid job type")
	}

	// job to data will validate data as well
	_, _, err := p.jobToData(j)
	if err != nil {
		return err
	}

	return nil

}

// get data from job
func (p *SubmitMessagesProcessor) jobToData(j queue.Job) (int, int64, error) {

	// fetch chatID
	chatID, oki := j.Data["chat_id"].(int)
	if !oki {
		return 0, 0, errors.New("expected chat id to be a float64")
	}

	messageID, oki := j.Data["db_message_id"].(int64)
	if !oki {
		return 0, 0, errors.New("expected message id to be a float64")
	}

	return int(chatID), messageID, nil

}

func (p *SubmitMessagesProcessor) Process(j queue.Job) error {

	// make sure type is correct
	if p.Type() != j.Type {
		return errors.New("invalid job type")
	}

	// get data from job map
	chatId, messageID, err := p.jobToData(j)
	if err != nil {
		return err
	}

	chat, err := p.chatDB.GetChat(chatId)
	if err != nil {
		return err
	}

	// fetch message
	msg, err := chat.GetMessage(messageID)
	if err != nil {
		return err
	}
	if msg == nil {
		return errors.New("failed to fetch message")
	}

	if chat.IsGroupChat() {

		for _, groupMember := range chat.Partners {

			// make sure we don't send to our self
			idPubKey, err := p.chat.km.IdentityPublicKey()
			if err != nil {
				return err
			}
			if idPubKey == hex.EncodeToString(groupMember) {
				continue
			}

			msg.GroupChatID = chat.GroupChatRemoteID
			if err := p.chat.SendMessage(groupMember, *msg); err != nil {
				return err
			}

		}
	} else {
		err = p.chat.SendMessage(chat.Partner, *msg)
		if err != nil {
			return err
		}
	}

	// delete job
	return p.queue.DeleteJob(j)

}

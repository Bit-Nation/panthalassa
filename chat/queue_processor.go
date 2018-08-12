package chat

import (
	"encoding/json"
	"errors"

	db "github.com/Bit-Nation/panthalassa/db"
	queue "github.com/Bit-Nation/panthalassa/queue"
	ed25519 "golang.org/x/crypto/ed25519"
)

// processor that submits messages from the queue to the backend
type SubmitMessagesProcessor struct {
	chat  *Chat
	msgDB db.ChatMessageStorage
	queue *queue.Queue
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
func (p *SubmitMessagesProcessor) jobToData(j queue.Job) (ed25519.PublicKey, int64, error) {
	// validate message id
	messageIdNumber, oki := j.Data["db_message_id"].(json.Number)
	if !oki {
		return nil, 0, errors.New("message id is not of type json.Number")
	}

	messageId, messageIdErr := messageIdNumber.Int64()

	if messageIdErr != nil {
		return nil, 0, messageIdErr
	}

	if messageId == 0 {
		return nil, 0, errors.New("message id is 0")
	}

	// check partner
	partner, oki := j.Data["partner"].(ed25519.PublicKey)
	if !oki {
		return nil, 0, errors.New("partner is not of type []byte")
	}
	if len(partner) != 32 {
		return nil, 0, errors.New("invalid partner id length")
	}
	return partner, messageId, nil
}

func (p *SubmitMessagesProcessor) Process(j queue.Job) error {

	// make sure type is correct
	if p.Type() != j.Type {
		return errors.New("invalid job type")
	}

	// get data from job map
	partner, messageID, err := p.jobToData(j)
	if err != nil {
		return err
	}

	// fetch message
	msg, err := p.msgDB.GetMessage(partner, messageID)
	if err != nil {
		return err
	}
	if msg == nil {
		return errors.New("failed to fetch message")
	}

	// send message
	err = p.chat.SendMessage(partner, *msg)
	if err != nil {
		return err
	}

	// delete job
	return p.queue.DeleteJob(j)

}

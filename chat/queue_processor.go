package chat

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"

	db "github.com/Bit-Nation/panthalassa/db"
	queue "github.com/Bit-Nation/panthalassa/queue"
	ed25519 "golang.org/x/crypto/ed25519"
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
func (p *SubmitMessagesProcessor) jobToData(j queue.Job) (ed25519.PublicKey, int64, error) {
	// validate message id
	var messageId int64
	var messageIdErr error
	messageIdInterface, oki := j.Data["db_message_id"]
	if !oki {
		return nil, 0, errors.New("db_message_id is missing")
	}
	switch msgId := messageIdInterface.(type) {
	case json.Number:
		messageId, messageIdErr = msgId.Int64()
		if messageIdErr != nil {
			return nil, 0, messageIdErr
		}
	case int64:
		messageId = msgId
	default:
		return nil, 0, errors.New("Expected : json.Number/int64, Got : " + fmt.Sprint(msgId))
	}
	if messageId == 0 {
		return nil, 0, errors.New("message id is 0")
	}
	// check partner
	var partner []byte
	var partnerErr error
	partnerInterface, oki := j.Data["partner"]
	if !oki {
		return nil, 0, errors.New("partner is missing")
	}
	switch potentialPartner := partnerInterface.(type) {
	case string:
		partner, partnerErr = base64.StdEncoding.DecodeString(potentialPartner)
		if partnerErr != nil {
			return nil, 0, partnerErr
		}
	case ed25519.PublicKey:
		partner = potentialPartner
	default:
		return nil, 0, errors.New("Expected : base64 string/ed25519.PublicKey, Got : " + fmt.Sprint(potentialPartner))
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

	chat, err := p.chatDB.GetChat(partner)
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

	// send message
	err = p.chat.SendMessage(partner, *msg)
	if err != nil {
		return err
	}

	// delete job
	return p.queue.DeleteJob(j)

}

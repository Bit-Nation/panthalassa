package chat

import (
	"time"

	x3dh "github.com/Bit-Nation/x3dh"
)

type Info struct {
	SharedSecret x3dh.SharedSecret
}

// send a message to a chat partner
func (c *Chat) CreateHumanMessage(msg string, sec x3dh.SharedSecret) (Message, error) {

	// create doublerachet cipher text
	encryptedMessage, err := c.encryptMessage(sec, []byte(msg))
	if err != nil {
		return Message{}, err
	}

	m := Message{
		Type:                 "HUMAN_MESSAGE",
		SendAt:               time.Now(),
		DoubleratchetMessage: encryptedMessage,
	}

	// sign message
	err = m.Sign(c.km)
	if err != nil {
		return Message{}, err
	}

	return m, err

}

package chat

import (
	"time"

	profile "github.com/Bit-Nation/panthalassa/profile"
	x3dh "github.com/Bit-Nation/x3dh"
)

type Info struct {
	SharedSecret x3dh.SharedSecret
}

// send a message to a chat partner
func (c *Chat) SendHumanMessage(msg string, profile profile.Profile, sec x3dh.SharedSecret) (Message, error) {

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

	// publish the message
	err = c.publishMessage(m)

	return m, err

}

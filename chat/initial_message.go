package chat

import (
	"encoding/hex"

	x3dh "github.com/Bit-Nation/x3dh"
	"golang.org/x/crypto/ed25519"
)

type Initialisation struct {
	Msg    Message `json:"message"`
	Secret string  `json:"secret"`
}

func (c *Chat) InitializeChat(idPubKey ed25519.PublicKey) (Message, x3dh.InitializedProtocol, error) {

	// init the x3dh protocol
	preKeyBundle, err := c.FetchPreKeyBundle(idPubKey)
	if err != nil {
		return Message{}, x3dh.InitializedProtocol{}, err
	}

	// init the protocol
	ip, err := c.CreateSharedSecret(preKeyBundle)
	if err != nil {
		return Message{}, x3dh.InitializedProtocol{}, err
	}

	// create encrypted message
	msg, err := c.encryptMessage(ip.SharedSecret, []byte("hi"))
	if err != nil {
		return Message{}, x3dh.InitializedProtocol{}, err
	}

	// construct message
	m := Message{
		Type: "PROTOCOL_INITIALISATION",
		AdditionalData: map[string]string{
			"used_one_time_pre_key": hex.EncodeToString(ip.UsedOneTimePreKey[:]),
			"used_signed_pre_key":   hex.EncodeToString(ip.UsedSignedPreKey[:]),
			"ephemeral_key":         hex.EncodeToString(ip.EphemeralKey[:]),
		},
		DoubleratchetMessage: msg,
	}

	// sign message
	err = m.Sign(c.km)
	if err != nil {
		return Message{}, x3dh.InitializedProtocol{}, err
	}

	// publish message to the network
	err = c.publishMessage(m)

	return m, ip, err
}

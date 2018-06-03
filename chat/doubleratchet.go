package chat

import (
	"encoding/json"
	
	profile "github.com/Bit-Nation/panthalassa/profile"
	x3dh "github.com/Bit-Nation/x3dh"
	doubleratchet "github.com/tiabc/doubleratchet"
)

// local double ratchet key pair
// as a local helper
type doubleratchetKeyPair struct {
	kp x3dh.KeyPair
}
func (p doubleratchetKeyPair) PrivateKey() doubleratchet.Key {
	var byt [32]byte = p.kp.PrivateKey
	return byt
}
func (p doubleratchetKeyPair) PublicKey() doubleratchet.Key {
	var byt [32]byte = p.kp.PublicKey
	return byt
}

// encrypt a double rachet message
func (c *Chat) encryptMessage(secret x3dh.SharedSecret, data []byte) (doubleratchet.Message, error) {

	// fetch chat ID key pair from key manager
	chatIdKey, err := c.km.ChatIdKeyPair()
	if err != nil {
		return doubleratchet.Message{}, err
	}

	// decoded secret
	var secBytes [32]byte = secret

	// create double rachet session
	s, err := doubleratchet.New(secBytes, doubleratchetKeyPair{
		kp: chatIdKey,
	}, doubleratchet.WithKeysStorage(c.keyStore))
	if err != nil {
		return doubleratchet.Message{}, err
	}

	// encrypt message
	return s.RatchetEncrypt(data, []byte{}), nil

}

// decrypt a message
func (c *Chat) DecryptMessage(secret x3dh.SharedSecret, profile profile.Profile, msg string) (string, error) {

	// chat partner chat id public key
	chatIdKey := profile.GetChatIDPublicKey()

	// unmarshal the message
	m := doubleratchet.Message{}
	err := json.Unmarshal([]byte(msg), m)
	if err != nil {
		return "", err
	}

	var secBytes [32]byte = secret
	var remotePub [32]byte = chatIdKey

	// create double rachet instance
	s, err := doubleratchet.NewWithRemoteKey(
		secBytes,
		remotePub,
		doubleratchet.WithKeysStorage(c.keyStore),
	)
	if err != nil {
		return "", err
	}

	// decrypt
	dec, err := s.RatchetDecrypt(m, nil)

	return string(dec), err
}

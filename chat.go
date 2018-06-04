package panthalassa

import (
	"encoding/hex"
	"encoding/json"
	"errors"

	chat "github.com/Bit-Nation/panthalassa/chat"
	profile "github.com/Bit-Nation/panthalassa/profile"
)

// create new pre key bundle
func NewPreKeyBundle() (string, error) {

	if panthalassaInstance == nil {
		return "", errors.New("please start panthalassa first")
	}

	bundle, err := panthalassaInstance.chat.NewPreKeyBundle()
	if err != nil {
		return "", err
	}

	bundleRaw, err := bundle.Marshal()
	if err != nil {
		return "", err
	}

	return string(bundleRaw), nil

}

// initialize chat with given identity key and pre key bundle
func InitializeChat(identityPublicKey, preKeyBundle string) (string, error) {

	if panthalassaInstance == nil {
		return "", errors.New("please start panthalassa first")
	}

	// decode public key
	pubKey, err := hex.DecodeString(identityPublicKey)
	if err != nil {
		return "", err
	}
	if len(pubKey) != 32 {
		return "", errors.New("public key must have length of 32 bytes")
	}

	// decode pre key bundle
	bundle, err := chat.UnmarshalPreKeyBundle([]byte(preKeyBundle))
	if err != nil {
		return "", err
	}

	msg, initializedProtocol, err := panthalassaInstance.chat.InitializeChat(pubKey, bundle)
	if err != nil {
		return "", err
	}

	exportedSecret, err := chat.EncryptX3DHSecret(initializedProtocol.SharedSecret, panthalassaInstance.km)
	if err != nil {
		return "", err
	}

	initialProtocol, err := json.Marshal(struct {
		Message chat.Message `json:"message"`
		Secret  string       `json:"secret"`
	}{
		Message: msg,
		Secret:  exportedSecret,
	})

	return string(initialProtocol), err

}

// create message
func CreateHumanMessage(rawMsg, rawProfile, secret string) (string, error) {

	if panthalassaInstance == nil {
		return "", errors.New("please start panthalassa first")
	}

	// unmarshal profile
	prof, err := profile.Unmarshal(rawProfile)
	if err != nil {
		return "", err
	}

	// shared secret
	sharedSecret, err := chat.DecryptX3DHSecret(secret, panthalassaInstance.km)
	if err != nil {
		return "", err
	}

	// create message
	msg, err := panthalassaInstance.chat.CreateHumanMessage(rawMsg, prof, sharedSecret)
	if err != nil {
		return "", err
	}

	// marshal message
	m, err := msg.Marshal()
	if err != nil {
		return "", err
	}

	return string(m), nil

}

// decrypt a chat message
func DecryptMessage(encryptedMessage, rawProfile, secret string) (string, error) {

	if panthalassaInstance == nil {
		return "", errors.New("please start panthalassa first")
	}

	// unmarshal profile
	prof, err := profile.Unmarshal(rawProfile)
	if err != nil {
		return "", err
	}

	// shared secret
	sharedSecret, err := chat.DecryptX3DHSecret(secret, panthalassaInstance.km)
	if err != nil {
		return "", err
	}

	return panthalassaInstance.chat.DecryptMessage(sharedSecret, prof, encryptedMessage)

}

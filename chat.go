package panthalassa

import (
	"encoding/hex"
	"encoding/json"
	"errors"

	chat "github.com/Bit-Nation/panthalassa/chat"
	aes "github.com/Bit-Nation/panthalassa/crypto/aes"
)

// create new pre key bundle
func NewPreKeyBundle() (string, error) {

	if panthalassaInstance == nil {
		return "", errors.New("please start panthalassa first")
	}

	// create new per key bundle
	bundle, err := panthalassaInstance.chat.NewPreKeyBundle()
	if err != nil {
		return "", err
	}

	// marshal public part
	publicPart, err := bundle.PublicPart.Marshal()
	if err != nil {
		return "", err
	}

	// marshal private part and encrypt
	privatePart, err := bundle.PrivatePart.Marshal()
	if err != nil {
		return "", err
	}

	// encrypt private part
	privtePartCipherText, err := panthalassaInstance.km.AESEncrypt(privatePart)
	if err != nil {
		return "", err
	}

	// marshal pre key bundle
	preKeyBundle, err := json.Marshal(struct {
		PublicPart  string         `json:"public_part"`
		PrivatePart aes.CipherText `json:"private_part"`
	}{
		PublicPart:  string(publicPart),
		PrivatePart: privtePartCipherText,
	})

	if err != nil {
		return "", err
	}

	return string(preKeyBundle), nil

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
		Message chat.Message   `json:"message"`
		Secret  aes.CipherText `json:"secret"`
	}{
		Message: msg,
		Secret:  exportedSecret,
	})

	return string(initialProtocol), err

}

// create message
// secret should be a aes cipher text as string
func CreateHumanMessage(rawMsg, secret string) (string, error) {

	if panthalassaInstance == nil {
		return "", errors.New("please start panthalassa first")
	}

	// unmarshal raw secret (secret is a cipher text)
	cipherText, err := aes.Unmarshal([]byte(rawMsg))
	if err != nil {
		return "", err
	}

	// shared secret
	sharedSecret, err := chat.DecryptX3DHSecret(cipherText, panthalassaInstance.km)
	if err != nil {
		return "", err
	}

	// create message
	msg, err := panthalassaInstance.chat.CreateHumanMessage(rawMsg, sharedSecret)
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
// secret should be a cipher text as string
func DecryptMessage(message, secret string) (string, error) {

	if panthalassaInstance == nil {
		return "", errors.New("please start panthalassa first")
	}

	ct, err := aes.Unmarshal([]byte(secret))
	if err != nil {
		return "", err
	}

	// shared secret
	sharedSecret, err := chat.DecryptX3DHSecret(ct, panthalassaInstance.km)
	if err != nil {
		return "", err
	}

	// unmarshal message
	var m chat.Message
	if err := json.Unmarshal([]byte(message), &m); err != nil {
		return "", err
	}

	return panthalassaInstance.chat.DecryptMessage(sharedSecret, m)

}

// return a encrypted shared secret used by the double rachet
func HandleInitialMessage(message, preKeyBundlePrivatePart string) (string, error) {

	if panthalassaInstance == nil {
		return "", errors.New("please start panthalassa first")
	}

	// unmarshal message
	var m chat.Message
	if err := json.Unmarshal([]byte(message), &m); err != nil {
		return "", err
	}

	// unmarshal pre key bundle private part
	var p chat.PreKeyBundlePrivate
	if err := json.Unmarshal([]byte(preKeyBundlePrivatePart), &p); err != nil {
		return "", err
	}

	sec, err := panthalassaInstance.chat.HandleInitialMessage(m, p)
	if err != nil {
		return "", err
	}

	ct, err := chat.EncryptX3DHSecret(sec, panthalassaInstance.km)
	if err != nil {
		return "", err
	}

	ctRaw, err := ct.Marshal()
	if err != nil {
		return "", err
	}

	return string(ctRaw), err

}

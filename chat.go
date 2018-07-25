package panthalassa

import (
	"encoding/hex"
	"encoding/json"
	"errors"

	db "github.com/Bit-Nation/panthalassa/db"
	bpb "github.com/Bit-Nation/protobuffers"
	proto "github.com/gogo/protobuf/proto"
)

func SendMessage(partner, message string) error {

	// make sure panthalassa has been started
	if panthalassaInstance == nil {
		return errors.New("you have to start panthalassa first")
	}

	// partner public key
	partnerPub, err := hex.DecodeString(partner)
	if err != nil {
		return err
	}

	// make sure public key has the right length
	if len(partnerPub) != 32 {
		return errors.New("partner must have a length of 32 bytes")
	}

	// unmarshal the plain message
	plainMsg := bpb.PlainChatMessage{}
	if err := proto.Unmarshal([]byte(message), &plainMsg); err != nil {
		return err
	}

	// persist private message
	return panthalassaInstance.chat.SavePrivateMessage(partnerPub, plainMsg)

}

func AllChats() (string, error) {

	// make sure panthalassa has been started
	if panthalassaInstance == nil {
		return "", errors.New("you have to start panthalassa first")
	}

	chats, err := panthalassaInstance.chat.AllChats()
	if err != nil {
		return "", err
	}

	chatsStr := []string{}
	for _, chat := range chats {
		chatsStr = append(chatsStr, hex.EncodeToString(chat))
	}

	chatList, err := json.Marshal(chatsStr)
	if err != nil {
		return "", err
	}

	return string(chatList), nil
}

func Messages(partner string, start int64, amount uint) (string, error) {

	// make sure panthalassa has been started
	if panthalassaInstance == nil {
		return "", errors.New("you have to start panthalassa first")
	}

	// partner public key
	partnerPub, err := hex.DecodeString(partner)
	if err != nil {
		return "", err
	}

	// make sure public key has the right length
	if len(partnerPub) != 32 {
		return "", errors.New("partner must have a length of 32 bytes")
	}

	// database messages
	databaseMessages, err := panthalassaInstance.chat.Messages(partnerPub, start, amount)
	if err != nil {
		return "", err
	}

	// plain messages
	plainMessages := map[int64]string{}

	// decrypt message
	for key, dbMessage := range databaseMessages {
		plainMessage, err := panthalassaInstance.km.AESDecrypt(dbMessage.Message)
		if err != nil {
			return "", err
		}
		plainMessages[key] = string(plainMessage)
	}

	// marshal messages
	messages, err := json.Marshal(plainMessages)
	if err != nil {
		return "", err
	}

	return string(messages), nil

}

func ImportOldSentMessage(partner, message string) error {

	// make sure panthalassa has been started
	if panthalassaInstance == nil {
		return errors.New("you have to start panthalassa first")
	}

	// partner public key
	partnerPub, err := hex.DecodeString(partner)
	if err != nil {
		return err
	}

	// make sure public key has the right length
	if len(partnerPub) != 32 {
		return errors.New("partner must have a length of 32 bytes")
	}

	// unmarshal the plain message
	plainMsg := bpb.PlainChatMessage{}
	if err := proto.Unmarshal([]byte(message), &plainMsg); err != nil {
		return err
	}

	return panthalassaInstance.msgDB.PersistMessage(partnerPub, plainMsg, false, db.StatusSent)

}

func ImportOldReceivedMessage(partner, message string) error {

	// make sure panthalassa has been started
	if panthalassaInstance == nil {
		return errors.New("you have to start panthalassa first")
	}

	// partner public key
	partnerPub, err := hex.DecodeString(partner)
	if err != nil {
		return err
	}

	// make sure public key has the right length
	if len(partnerPub) != 32 {
		return errors.New("partner must have a length of 32 bytes")
	}

	// unmarshal the plain message
	plainMsg := bpb.PlainChatMessage{}
	if err := proto.Unmarshal([]byte(message), &plainMsg); err != nil {
		return err
	}

	return panthalassaInstance.msgDB.PersistReceivedMessage(partnerPub, plainMsg)

}

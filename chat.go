package panthalassa

import (
	"encoding/hex"
	"encoding/json"
	"errors"

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

	// persist private message
	return panthalassaInstance.chat.SavePrivateMessage(partnerPub, []byte(message))

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
	plainMessages := []map[string]interface{}{}

	// decrypt message
	for key, dbMessage := range databaseMessages {
		plainMessage, err := panthalassaInstance.km.AESDecrypt(dbMessage.Message)
		if err != nil {
			return "", err
		}
		msg := bpb.PlainChatMessage{}
		if err := proto.Unmarshal(plainMessage, &msg); err != nil {
			return "", err
		}
		plainMessages = append(plainMessages, map[string]interface{}{
			"message_id": key,
			"message": map[string]interface{}{
				"content":    string(msg.Message),
				"created_at": msg.CreatedAt,
			},
		})
	}

	// marshal messages
	messages, err := json.Marshal(plainMessages)
	if err != nil {
		return "", err
	}

	return string(messages), nil

}

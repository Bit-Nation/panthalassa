package panthalassa

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"strconv"

	ed25519 "golang.org/x/crypto/ed25519"
)

func SendMessage(chatID int, message string) error {

	// make sure panthalassa has been started
	if panthalassaInstance == nil {
		return errors.New("you have to start panthalassa first")
	}

	// persist message
	return panthalassaInstance.chat.SaveMessage(chatID, []byte(message))

}

func AddUsersToGroupChat(users string, chatID int) error {

	// make sure panthalassa has been started
	if panthalassaInstance == nil {
		return errors.New("you have to start panthalassa first")
	}

	// json unmarshal partners
	rawPartners := []string{}
	if err := json.Unmarshal([]byte(users), &rawPartners); err != nil {
		return err
	}

	// unmarshal hex partners
	partners := []ed25519.PublicKey{}
	for _, hexPartner := range rawPartners {
		rawPartner, err := hex.DecodeString(hexPartner)
		if err != nil {
			return err
		}
		partners = append(partners, rawPartner)
	}

	return panthalassaInstance.chat.AddUserToGroupChat(partners, chatID)

}

func CreatePrivateChat(partnerStr string) (int, error) {

	// make sure panthalassa has been started
	if panthalassaInstance == nil {
		return 0, errors.New("you have to start panthalassa first")
	}

	partner, err := hex.DecodeString(partnerStr)
	if err != nil {
		return 0, err
	}

	return panthalassaInstance.chatDB.CreateChat(partner)

}

// return chatID
func CreateGroupChat(users string, name string) (int, error) {

	// make sure panthalassa has been started
	if panthalassaInstance == nil {
		return 0, errors.New("you have to start panthalassa first")
	}

	// json unmarshal partners
	rawPartners := []string{}
	if err := json.Unmarshal([]byte(users), &rawPartners); err != nil {
		return 0, err
	}

	// unmarshal hex partners
	partners := []ed25519.PublicKey{}
	for _, hexPartner := range rawPartners {
		rawPartner, err := hex.DecodeString(hexPartner)
		if err != nil {
			return 0, err
		}
		partners = append(partners, rawPartner)
	}

	return panthalassaInstance.chat.CreateGroupChat(partners, name)

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

	chatsRep := []map[string]interface{}{}
	for _, chat := range chats {
		chatsRep = append(chatsRep, map[string]interface{}{
			"chat_id":         chat.ID,
			"chat_partner":    chat.Partner,
			"unread_messages": chat.UnreadMessages,
			"group_chat_name": chat.GroupChatName,
			"partners":        chat.Partners,
		})
	}

	chatList, err := json.Marshal(chatsRep)
	if err != nil {
		return "", err
	}

	return string(chatList), nil
}

func Messages(chatID int, startStr string, amount int) (string, error) {

	// unmarshal start
	start, err := strconv.ParseInt(startStr, 10, 64)
	if err != nil {
		return "", err
	}

	// make sure panthalassa has been started
	if panthalassaInstance == nil {
		return "", errors.New("you have to start panthalassa first")
	}

	// fetch chat
	chat, err := panthalassaInstance.chatDB.GetChat(chatID)
	if err != nil {
		return "", err
	}

	if chat == nil {
		return "", errors.New("chat doesn't exit")
	}

	// database messages
	databaseMessages, err := chat.Messages(start, uint(amount))
	if err != nil {
		return "", err
	}

	// plain messages
	plainMessages := []map[string]interface{}{}

	// decrypt message
	for _, msg := range databaseMessages {
		dapp := ""
		if msg.DApp != nil {
			dappBytes, err := json.Marshal(msg.DApp)
			if err != nil {
				logger.Error(err)
			} else {
				dapp = string(dappBytes)
			}
		}
		plainMessages = append(plainMessages, map[string]interface{}{
			"db_id":      strconv.FormatInt(msg.UniqueMsgID, 10),
			"content":    string(msg.Message),
			"created_at": msg.CreatedAt,
			"received":   msg.Received,
			"dapp":       dapp,
		})
	}

	// marshal messages
	messages, err := json.Marshal(plainMessages)
	if err != nil {
		return "", err
	}

	return string(messages), nil

}

func MarkMessagesAsRead(chatID int) error {

	// make sure panthalassa has been started
	if panthalassaInstance == nil {
		return errors.New("you have to start panthalassa first")
	}

	return panthalassaInstance.chat.MarkMessagesAsRead(chatID)

}

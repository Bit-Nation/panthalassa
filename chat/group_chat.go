package chat

import (
	"errors"

	db "github.com/Bit-Nation/panthalassa/db"
	ed25519 "golang.org/x/crypto/ed25519"
)

func (c *Chat) AddUserToGroupChat(partners []ed25519.PublicKey, chatID int) error {

	chat, err := c.chatStorage.GetChat(chatID)
	if err != nil {
		return err
	}

	if chat == nil {
		return errors.New("couldn't find chat")
	}

	if !chat.IsGroupChat() {
		return errors.New("chat must be a group chat")
	}

	// update chat partners
	if err := chat.AddChatPartners(partners); err != nil {
		return err
	}

	// fetch chat
	for _, partner := range partners {

		// fetch chat
		partnerChat, err := c.chatStorage.GetChatByPartner(partner)
		if err != nil {
			return err
		}
		if partnerChat == nil {
			if _, err := c.chatStorage.CreateChat(partner); err != nil {
				return err
			}
		}
		partnerChat, err = c.chatStorage.GetChatByPartner(partner)
		if err != nil {
			return err
		}
		if partnerChat == nil {
			return errors.New("chat with partner should exist at this point in time")
		}

		// persist message
		msg := db.Message{
			AddUserToChat: &db.AddUserToChat{
				Users:          partners,
				SharedSecret:   chat.GroupChatSharedSecret,
				SharedSecretID: chat.GroupSharedSecretID,
				ChatID:         chat.GroupChatRemoteID,
			},
		}
		if err := partnerChat.PersistMessage(msg); err != nil {
			return err
		}

	}

	return nil

}

func (c *Chat) CreateGroupChat(partners []ed25519.PublicKey) (int, error) {

	// create chat
	chatID, err := c.chatStorage.CreateGroupChat(partners)
	if err != nil {
		return 0, err
	}

	// fetch group chat
	groupChat, err := c.chatStorage.GetChat(chatID)
	if err != nil {
		return 0, err
	}

	// fetch chat
	for _, partner := range partners {

		// fetch chat
		partnerChat, err := c.chatStorage.GetChatByPartner(partner)
		if err != nil {
			return 0, err
		}
		if partnerChat == nil {
			if _, err := c.chatStorage.CreateChat(partner); err != nil {
				return 0, err
			}
		}
		partnerChat, err = c.chatStorage.GetChatByPartner(partner)
		if err != nil {
			return 0, err
		}
		if partnerChat == nil {
			return 0, errors.New("chat with partner should exist at this point in time")
		}

		// persist message
		msg := db.Message{
			AddUserToChat: &db.AddUserToChat{
				Users:          partners,
				SharedSecret:   groupChat.GroupChatSharedSecret,
				SharedSecretID: groupChat.GroupSharedSecretID,
				ChatID:         groupChat.GroupChatRemoteID,
			},
		}
		if err := partnerChat.PersistMessage(msg); err != nil {
			return 0, err
		}

	}

	return chatID, nil

}

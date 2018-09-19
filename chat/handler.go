package chat

import (
	"crypto/rand"
	"encoding/json"
	"errors"
	"strconv"
	"sync"

	preKey "github.com/Bit-Nation/panthalassa/chat/prekey"
	db "github.com/Bit-Nation/panthalassa/db"
	queue "github.com/Bit-Nation/panthalassa/queue"
	bpb "github.com/Bit-Nation/protobuffers"
	x3dh "github.com/Bit-Nation/x3dh"
)

// handles a set of protobuf messages
func (c *Chat) messagesHandler(req *bpb.BackendMessage_Request) (*bpb.BackendMessage_Response, error) {

	wg := sync.WaitGroup{}
	if len(req.Messages) > 0 {
		for _, msg := range req.Messages {
			wg.Add(1)
			go func(msg *bpb.ChatMessage) {
				defer wg.Done()
				err := c.handleReceivedMessage(msg)
				if err != nil {
					logger.Error(err)
				}
			}(msg)
		}
		wg.Wait()
		return &bpb.BackendMessage_Response{}, nil
	}

	return nil, nil

}

// handle new one time pre keys request
func (c *Chat) oneTimePreKeysHandler(req *bpb.BackendMessage_Request) (*bpb.BackendMessage_Response, error) {

	// exit if not a request to fetch new one time pre keys
	if req.NewOneTimePreKeys == 0 {
		return nil, nil
	}

	// make sure request is valid
	if req.NewOneTimePreKeys > 100 {
		return nil, errors.New("requested more then the max allowed pre keys")
	}

	curve := x3dh.NewCurve25519(rand.Reader)

	// generate key pairs
	keyPairs := []x3dh.KeyPair{}
	for {
		if len(keyPairs) == int(req.NewOneTimePreKeys) {
			break
		}
		keyPair, err := curve.GenerateKeyPair()
		if err != nil {
			return nil, err
		}
		keyPairs = append(keyPairs, keyPair)
	}

	// persist all key pairs
	if err := c.oneTimePreKeyStorage.Put(keyPairs); err != nil {
		logger.Error(err)
		return nil, errors.New("failed to persist generated key pairs")
	}

	preKeys := []*bpb.PreKey{}

	// convert and sign one time pre keys
	for _, oneTimePreKey := range keyPairs {
		pk := preKey.PreKey{}
		pk.PublicKey = oneTimePreKey.PublicKey
		if err := pk.Sign(*c.km); err != nil {
			logger.Error(err)
			return nil, errors.New("failed to sign one time pre key")
		}
		pkProto, err := pk.ToProtobuf()
		if err != nil {
			logger.Error(err)
			return nil, errors.New("failed to convert pre key to protobuf")
		}
		preKeys = append(preKeys, &pkProto)
	}

	return &bpb.BackendMessage_Response{
		OneTimePrekeys: preKeys,
	}, nil

}

func (c *Chat) handlePersistedMessage(e db.MessagePersistedEvent) {

	// if the message was received we want to
	// to switch the chat boolean flag (UnreadMessages) to true
	// and tell the client that we got new unread messages
	if e.Message.Received {
		if err := c.chatStorage.UnreadMessages(e.Chat); err != nil {
			logger.Error(err)
			return
		}
		c.uiApi.Send("CHAT:UNREAD", map[string]interface{}{
			"chat": e.Chat.Partner,
		})
	}

	// when the handled message was not received we would like to send it
	if !e.Message.Received {
		// add to queue
		err := c.queue.AddJob(queue.Job{
			Type: "MESSAGE:SUBMIT",
			Data: map[string]interface{}{
				"chat_id":       e.Chat.ID,
				"db_message_id": e.Message.UniqueMsgID,
			},
		})
		if err != nil {
			logger.Error(err)
		}
	}

	dapp := ""
	if e.Message.DApp != nil {
		dappBytes, err := json.Marshal(e.Message.DApp)
		if err != nil {
			logger.Error(err)
		} else {
			dapp = string(dappBytes)
		}
	}

	if e.Message.Status == db.StatusPersisted {
		c.uiApi.Send("MESSAGE:PERSISTED", map[string]interface{}{
			"db_id":      strconv.FormatInt(e.Message.UniqueMsgID, 10),
			"content":    string(e.Message.Message),
			"created_at": e.Message.CreatedAt,
			"chat":       e.Chat.ID,
			"received":   e.Message.Received,
			"dapp":       dapp,
		})
	}

	if e.Message.Status == db.StatusDelivered {
		c.uiApi.Send("MESSAGE:DELIVERED", map[string]interface{}{
			"db_id":      strconv.FormatInt(e.Message.UniqueMsgID, 10),
			"content":    string(e.Message.Message),
			"created_at": e.Message.CreatedAt,
			"chat":       e.Chat.ID,
			"received":   e.Message.Received,
			"dapp":       dapp,
		})
	}

	if e.Message.Received {
		c.uiApi.Send("MESSAGE:RECEIVED", map[string]interface{}{
			"db_id":      strconv.FormatInt(e.Message.UniqueMsgID, 10),
			"content":    string(e.Message.Message),
			"created_at": e.Message.CreatedAt,
			"chat":       e.Chat.ID,
			"received":   e.Message.Received,
			"dapp":       dapp,
		})
	}

}

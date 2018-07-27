package db

import (
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"testing"

	bpb "github.com/Bit-Nation/protobuffers"
	bolt "github.com/coreos/bbolt"
	proto "github.com/gogo/protobuf/proto"
	require "github.com/stretchr/testify/require"
	ed25519 "golang.org/x/crypto/ed25519"
)

func TestBoltChatMessageStorage_PersistMessage(t *testing.T) {

	// setup
	db := createDB()
	km := createKeyManager()
	partner, _, err := ed25519.GenerateKey(rand.Reader)
	require.Nil(t, err)

	// event listener
	eventListener := make(chan PlainMessagePersistedEvent, 1)

	// persist message
	storage := NewChatMessageStorage(db, []chan PlainMessagePersistedEvent{eventListener}, km)
	err = storage.persistMessage(partner, bpb.PlainChatMessage{
		Message:   []byte("hi there"),
		CreatedAt: 33,
		MessageID: "the_message_id",
	}, true, 300)
	require.Nil(t, err)

	// assertions
	allChatPartners, err := storage.AllChats()
	require.Nil(t, err)
	require.Equal(t, hex.EncodeToString(partner), hex.EncodeToString(allChatPartners[0]))

	err = db.View(func(tx *bolt.Tx) error {
		// private messages
		privateMessage := tx.Bucket(privateChatBucketName)
		require.NotNil(t, privateMessage)

		// partner
		partnerBucket := privateMessage.Bucket(partner)
		require.NotNil(t, partnerBucket)

		// fetch message
		messageKey := make([]byte, 8)
		binary.BigEndian.PutUint64(messageKey, 33)
		fetchedMessage := partnerBucket.Get(messageKey)

		// unmarshal message
		message := Message{}
		if err := json.Unmarshal(fetchedMessage, &message); err != nil {
			return err
		}

		// decrypted message
		plainMessage, err := km.AESDecrypt(message.Message)
		require.Nil(t, err)

		protoMsg := bpb.PlainChatMessage{}
		if err := proto.Unmarshal(plainMessage, &protoMsg); err != nil {
			return err
		}

		// plain message
		plainMsg := bpb.PlainChatMessage{}
		require.Nil(t, proto.Unmarshal(plainMessage, &plainMsg))

		// test database message
		require.Equal(t, uint32(1), message.Version)
		require.Equal(t, Status(300), message.Status)
		require.Equal(t, "the_message_id", message.ID)
		require.Equal(t, "hi there", string(plainMsg.Message))
		require.Equal(t, true, message.Received)

		// test event listener
		persistedMessage := <-eventListener
		require.Equal(t, "hi there", string(persistedMessage.Msg.Message))
		require.Equal(t, int64(33), persistedMessage.Msg.CreatedAt)
		require.Equal(t, "the_message_id", persistedMessage.Msg.MessageID)
		require.Equal(t, "the_message_id", persistedMessage.DBMsg.ID)
		require.Equal(t, Status(300), persistedMessage.DBMsg.Status)
		require.Equal(t, uint32(1), persistedMessage.DBMsg.Version)
		require.Equal(t, true, persistedMessage.DBMsg.Received)

		return nil
	})
	require.Nil(t, err)

	// test the fetch messages
	messages, err := storage.Messages(partner, int64(33), 10)
	require.Nil(t, err)
	message := messages[33]
	require.Equal(t, "the_message_id", message.ID)
	require.Equal(t, true, message.Received)
	require.Equal(t, uint32(1), message.Version)
	require.Equal(t, Status(300), message.Status)
}

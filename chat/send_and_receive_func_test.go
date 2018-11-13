package chat

import (
	"encoding/hex"
	"testing"
	"time"

	ed25519 "golang.org/x/crypto/ed25519"

	backend "github.com/Bit-Nation/panthalassa/backend"
	db "github.com/Bit-Nation/panthalassa/db"
	profile "github.com/Bit-Nation/panthalassa/profile"
	queue "github.com/Bit-Nation/panthalassa/queue"
	uiapi "github.com/Bit-Nation/panthalassa/uiapi"
	bpb "github.com/Bit-Nation/protobuffers"
	require "github.com/stretchr/testify/require"
)

type chatTestBackendTransport struct {
	sendChan    chan *bpb.BackendMessage
	receiveChan chan *bpb.BackendMessage
}

func (t *chatTestBackendTransport) Send(msg *bpb.BackendMessage) error {
	t.sendChan <- msg
	return nil
}

// will return the next message from the transport
func (t *chatTestBackendTransport) NextMessage() (*bpb.BackendMessage, error) {
	return <-t.receiveChan, nil
}

// close the transport
func (t *chatTestBackendTransport) Close() error {
	close(t.sendChan)
	close(t.receiveChan)
	return nil
}

type upStream struct{}

func (u *upStream) Send(data string) {}

func createAliceAndBob() (alice *Chat, aliceTrans *chatTestBackendTransport, bob *Chat, bobTrans *chatTestBackendTransport, err error) {

	// creates a chat. Don't forget to set the backend.
	createChat := func() (*Chat, *chatTestBackendTransport, error) {

		storm := createStorm()
		km := createKeyManager()

		chatStorage := db.NewChatStorage(storm, []func(e db.MessagePersistedEvent){}, km)
		ssStorage := db.NewBoltSharedSecretStorage(storm, km)
		drKeyStorage := db.NewBoltDRKeyStorage(storm, km)
		signedPreKeyStorage := db.NewBoltSignedPreKeyStorage(storm, km)
		oneTimePreKeyStorage := db.NewBoltOneTimePreKeyStorage(storm, km)
		userStorage := db.NewBoltUserStorage(storm)
		q := queue.New(queue.NewStorage(storm), 5, 1)

		trans := chatTestBackendTransport{
			sendChan:    make(chan *bpb.BackendMessage, 10),
			receiveChan: make(chan *bpb.BackendMessage, 10),
		}

		up := uiapi.New(&upStream{})
		b, err := backend.NewBackend(&trans, km, signedPreKeyStorage)
		if err != nil {
			return nil, nil, err
		}

		c, err := NewChat(Config{
			ChatStorage:          chatStorage,
			Backend:              b,
			SharedSecretDB:       ssStorage,
			KM:                   km,
			DRKeyStorage:         drKeyStorage,
			SignedPreKeyStorage:  signedPreKeyStorage,
			OneTimePreKeyStorage: oneTimePreKeyStorage,
			UserStorage:          userStorage,
			UiApi:                up,
			Queue:                q,
		})
		return c, &trans, err
	}

	// create alice
	alice, aliceTrans, err = createChat()
	if err != nil {
		return nil, nil, nil, nil, err
	}

	// create bob
	bob, bobTrans, err = createChat()
	if err != nil {
		return nil, nil, nil, nil, err
	}

	return

}

func TestChatBetweenAliceAndBob(t *testing.T) {

	alice, aliceTrans, bob, bobTrans, err := createAliceAndBob()

	// listen for bob's messages
	bobReceivedMsgChan := make(chan db.MessagePersistedEvent, 10)
	bob.chatStorage.AddListener(func(e db.MessagePersistedEvent) {
		bobReceivedMsgChan <- e
	})
	bobIDKeyStr, err := bob.km.IdentityPublicKey()
	require.Nil(t, err)
	bobIDKey, err := hex.DecodeString(bobIDKeyStr)
	require.Nil(t, err)

	// listen for alice messages
	aliceReceivedMsgChan := make(chan db.MessagePersistedEvent, 10)
	alice.chatStorage.AddListener(func(e db.MessagePersistedEvent) {
		aliceReceivedMsgChan <- e
	})
	aliceIDKeyStr, err := alice.km.IdentityPublicKey()
	require.Nil(t, err)
	aliceIDKey, err := hex.DecodeString(aliceIDKeyStr)
	require.Nil(t, err)

	// bob go
	go func() {

		bobSignedPreKey := new(bpb.PreKey)
		bobSignedPreKey = nil

		aliceSignedPreKey := new(bpb.PreKey)
		aliceSignedPreKey = nil

		for {
			select {
			case msg := <-bobTrans.sendChan:

				if msg.Request != nil {

					// handle uploaded pre key
					if msg.Request.NewSignedPreKey != nil {
						bobSignedPreKey = msg.Request.NewSignedPreKey
						bobTrans.receiveChan <- &bpb.BackendMessage{
							RequestID: msg.RequestID,
						}
						continue
					}

					if len(msg.Request.PreKeyBundle) != 0 {

						aliceProfile, err := profile.SignProfile("bob", "", "", *bob.km)
						if err != nil {
							panic(err)
						}
						aliceProfileProto, err := aliceProfile.ToProtobuf()
						if err != nil {
							panic(err)
						}

						bobTrans.receiveChan <- &bpb.BackendMessage{
							RequestID: msg.RequestID,
							Response: &bpb.BackendMessage_Response{
								PreKeyBundle: &bpb.BackendMessage_PreKeyBundle{
									SignedPreKey: aliceSignedPreKey,
									Profile:      aliceProfileProto,
								},
							},
						}
						continue
					}

					// handle message send to bob
					// (the nice thing is that we can write it from alice to bob :))
					if len(msg.Request.Messages) != 0 {
						aliceTrans.receiveChan <- msg
						bobTrans.receiveChan <- &bpb.BackendMessage{
							RequestID: msg.RequestID,
						}
						continue
					}

				}
				// alice
			case msg := <-aliceTrans.sendChan:

				if msg.Request != nil {
					// handle the pre key bundle request for bob
					if len(msg.Request.PreKeyBundle) != 0 {

						bobProfile, err := profile.SignProfile("bob", "", "", *bob.km)
						if err != nil {
							panic(err)
						}
						bobProfileProto, err := bobProfile.ToProtobuf()
						if err != nil {
							panic(err)
						}

						<-time.After(time.Second)

						aliceTrans.receiveChan <- &bpb.BackendMessage{
							RequestID: msg.RequestID,
							Response: &bpb.BackendMessage_Response{
								PreKeyBundle: &bpb.BackendMessage_PreKeyBundle{
									SignedPreKey: bobSignedPreKey,
									Profile:      bobProfileProto,
								},
							},
						}
						continue
					}

					// handle upload of our new pre key bundle
					if msg.Request.NewSignedPreKey != nil {
						aliceSignedPreKey = msg.Request.NewSignedPreKey
						aliceTrans.receiveChan <- &bpb.BackendMessage{
							RequestID: msg.RequestID,
						}
						continue
					}

					// handle message send to bob
					// (the nice thing is that we can write it from alice to bob :))
					if len(msg.Request.Messages) != 0 {
						bobTrans.receiveChan <- msg
						aliceTrans.receiveChan <- &bpb.BackendMessage{
							RequestID: msg.RequestID,
						}
						continue
					}
				}

			}

		}
	}()

	require.Nil(t, err)

	// persist private message for bob
	_, err = alice.chatStorage.CreateChat(bobIDKey)
	require.Nil(t, err)
	bobChat, err := alice.chatStorage.GetChatByPartner(bobIDKey)
	require.Nil(t, err)
	require.NotNil(t, bobChat)

	require.Nil(t, alice.SaveMessage(bobChat.ID, []byte("hi bob")))

	// done signal
	done := make(chan struct{}, 1)

	for {
		select {
		case msgEv := <-aliceReceivedMsgChan:

			msg := msgEv.Message

			if msg.Received {

				// make sure message is as we expect it to be
				require.Equal(t, "hi alice", string(msgEv.Message.Message))
				require.Equal(t, hex.EncodeToString(bobIDKey), hex.EncodeToString(msg.Sender))
				require.Equal(t, uint(1), msg.Version)
				require.Equal(t, db.StatusPersisted, msg.Status)

				// we need to wait a bit since this is called async
				// before the receive handling logic has the chance to update
				// the shared secret
				time.Sleep(time.Second)

				// make sure shared secret got accepted
				shSec, err := alice.sharedSecStorage.GetYoungest(bobIDKey)
				require.Nil(t, err)
				require.NotNil(t, shSec)
				require.True(t, shSec.Accepted)

				done <- struct{}{}
			}

		case msgEv := <-bobReceivedMsgChan:

			msg := msgEv.Message

			// handle received messages
			if msg.Received {

				// make sure the messages is as we expect it to be
				require.Equal(t, "hi bob", string(msg.Message))
				require.True(t, msg.Received)
				require.Equal(t, db.StatusPersisted, msg.Status)
				require.Equal(t, uint(1), msg.Version)

				// make sure shared secret got persisted
				shSec, err := bob.sharedSecStorage.GetYoungest(aliceIDKey)
				require.Nil(t, err)
				require.NotNil(t, shSec)
				require.Equal(t, hex.EncodeToString(aliceIDKey), hex.EncodeToString(shSec.Partner))
				require.True(t, shSec.Accepted)

				// fetch chat
				chat, err := bob.chatStorage.GetChatByPartner(aliceIDKey)
				require.Nil(t, err)

				err = bob.SaveMessage(chat.ID, []byte("hi alice"))
				require.Nil(t, err)
			}

		case <-done:
			return
		}
	}

}

func TestGroupChatBetweenAliceAndBob(t *testing.T) {
	// log.SetDebugLogging()
	alice, aliceTrans, bob, bobTrans, err := createAliceAndBob()

	// listen for bob's messages
	bobReceivedMsgChan := make(chan db.MessagePersistedEvent, 10)
	bob.chatStorage.AddListener(func(e db.MessagePersistedEvent) {
		bobReceivedMsgChan <- e
	})
	bobIDKeyStr, err := bob.km.IdentityPublicKey()
	require.Nil(t, err)
	bobIDKey, err := hex.DecodeString(bobIDKeyStr)
	require.Nil(t, err)

	// listen for alice messages
	aliceReceivedMsgChan := make(chan db.MessagePersistedEvent, 10)
	alice.chatStorage.AddListener(func(e db.MessagePersistedEvent) {
		aliceReceivedMsgChan <- e
	})
	aliceIDKeyStr, err := alice.km.IdentityPublicKey()
	require.Nil(t, err)
	aliceIDKey, err := hex.DecodeString(aliceIDKeyStr)
	require.Nil(t, err)

	// bob go
	go func() {

		bobSignedPreKey := new(bpb.PreKey)
		bobSignedPreKey = nil

		aliceSignedPreKey := new(bpb.PreKey)
		aliceSignedPreKey = nil

		for {
			select {
			case msg := <-bobTrans.sendChan:

				if msg.Request != nil {

					// handle uploaded pre key
					if msg.Request.NewSignedPreKey != nil {
						bobSignedPreKey = msg.Request.NewSignedPreKey
						bobTrans.receiveChan <- &bpb.BackendMessage{
							RequestID: msg.RequestID,
						}
						continue
					}

					if len(msg.Request.PreKeyBundle) != 0 {

						aliceProfile, err := profile.SignProfile("bob", "", "", *bob.km)
						if err != nil {
							panic(err)
						}
						aliceProfileProto, err := aliceProfile.ToProtobuf()
						if err != nil {
							panic(err)
						}

						bobTrans.receiveChan <- &bpb.BackendMessage{
							RequestID: msg.RequestID,
							Response: &bpb.BackendMessage_Response{
								PreKeyBundle: &bpb.BackendMessage_PreKeyBundle{
									SignedPreKey: aliceSignedPreKey,
									Profile:      aliceProfileProto,
								},
							},
						}
						continue
					}

					// handle message send to bob
					// (the nice thing is that we can write it from alice to bob :))
					if len(msg.Request.Messages) != 0 {
						aliceTrans.receiveChan <- msg
						bobTrans.receiveChan <- &bpb.BackendMessage{
							RequestID: msg.RequestID,
						}
						continue
					}

				}
				// alice
			case msg := <-aliceTrans.sendChan:

				if msg.Request != nil {
					// handle the pre key bundle request for bob
					if len(msg.Request.PreKeyBundle) != 0 {

						bobProfile, err := profile.SignProfile("bob", "", "", *bob.km)
						if err != nil {
							panic(err)
						}
						bobProfileProto, err := bobProfile.ToProtobuf()
						if err != nil {
							panic(err)
						}

						<-time.After(time.Second)

						aliceTrans.receiveChan <- &bpb.BackendMessage{
							RequestID: msg.RequestID,
							Response: &bpb.BackendMessage_Response{
								PreKeyBundle: &bpb.BackendMessage_PreKeyBundle{
									SignedPreKey: bobSignedPreKey,
									Profile:      bobProfileProto,
								},
							},
						}
						continue
					}

					// handle upload of our new pre key bundle
					if msg.Request.NewSignedPreKey != nil {
						aliceSignedPreKey = msg.Request.NewSignedPreKey
						aliceTrans.receiveChan <- &bpb.BackendMessage{
							RequestID: msg.RequestID,
						}
						continue
					}

					// handle message send to bob
					// (the nice thing is that we can write it from alice to bob :))
					if len(msg.Request.Messages) != 0 {
						bobTrans.receiveChan <- msg
						aliceTrans.receiveChan <- &bpb.BackendMessage{
							RequestID: msg.RequestID,
						}
						continue
					}
				}

			}

		}
	}()

	require.Nil(t, err)

	groupChatID, err := alice.CreateGroupChat([]ed25519.PublicKey{bobIDKey}, "Group between alice and bob")
	require.Nil(t, err)
	require.Nil(t, alice.SaveMessage(groupChatID, []byte("hi @all")))

	// done signal
	done := make(chan struct{}, 1)

	for {

		select {
		case msgEv := <-aliceReceivedMsgChan:

			msg := msgEv.Message

			if msg.Received {

				// make sure message is as we expect it to be
				require.Equal(t, "Greeting @all", string(msg.Message))
				require.Equal(t, "Group between alice and bob", msgEv.Chat.GroupChatName)
				require.Equal(t, hex.EncodeToString(bobIDKey), hex.EncodeToString(msg.Sender))
				require.Equal(t, uint(1), msg.Version)
				require.Equal(t, db.StatusPersisted, msg.Status)

				// we need to wait a bit since this is called async
				// before the receive handling logic has the chance to update
				// the shared secret
				time.Sleep(time.Second)

				// make sure shared secret got accepted
				shSec, err := alice.sharedSecStorage.GetYoungest(bobIDKey)
				require.Nil(t, err)
				require.NotNil(t, shSec)
				require.True(t, shSec.Accepted)

				done <- struct{}{}
			}

		case msgEv := <-bobReceivedMsgChan:

			msg := msgEv.Message

			// handle received messages
			if msg.Received {

				// make sure group ID has been set
				require.Equal(t, 200, len(msgEv.Chat.GroupChatRemoteID))

				// make sure the messages is as we expect it to be
				require.Equal(t, "hi @all", string(msg.Message))
				require.Equal(t, "Group between alice and bob", msgEv.Chat.GroupChatName)
				require.True(t, msg.Received)
				require.Equal(t, db.StatusPersisted, msg.Status)
				require.Equal(t, uint(1), msg.Version)

				// make sure shared secret got persisted
				shSec, err := bob.sharedSecStorage.GetYoungest(aliceIDKey)
				require.Nil(t, err)
				require.NotNil(t, shSec)
				require.Equal(t, hex.EncodeToString(aliceIDKey), hex.EncodeToString(shSec.Partner))
				require.True(t, shSec.Accepted)

				// send message back
				require.Nil(t, msgEv.Chat.SaveMessage([]byte("Greeting @all")))

			}

		case <-done:
			return
		case <-time.After(time.Second * 6):
			require.FailNow(t, "timed out")
		}
	}

}

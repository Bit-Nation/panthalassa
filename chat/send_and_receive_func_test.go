package chat

import (
	"encoding/hex"
	"fmt"
	"testing"
	"time"

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

					fmt.Println(msg)

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
	require.Nil(t, alice.SavePrivateMessage(bobIDKey, []byte("hi bob")))

	// done signal
	done := make(chan struct{}, 1)

	for {
		select {
		case msgEv := <-aliceReceivedMsgChan:
			fmt.Println(string(msgEv.Message.Message))
			fmt.Println(msgEv.Message.Received)
			if msgEv.Message.Received {
				require.Equal(t, "hi alice", string(msgEv.Message.Message))
				done <- struct{}{}
			}
		case msgEv := <-bobReceivedMsgChan:

			msg := msgEv.Message

			// handle received messages
			if msg.Received {
				// make sure the messages is as we expect it to be
				require.Equal(t, "hi bob", string(msg.Message))
				err := bob.SavePrivateMessage(aliceIDKey, []byte("hi alice"))
				require.Nil(t, err)
			}

		case <-done:
			fmt.Println("done")
		}
	}

}

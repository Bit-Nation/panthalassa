package message

import (
	db "github.com/Bit-Nation/panthalassa/db"
	ed25519 "golang.org/x/crypto/ed25519"
)

type testMessageStorage struct {
	persistMessageToSend   func(to ed25519.PublicKey, msg db.Message) error
	persistReceivedMessage func(partner ed25519.PublicKey, msg db.Message) error
	updateStatus           func(partner ed25519.PublicKey, msgID int64, newStatus db.Status) error
	messages               func(partner ed25519.PublicKey, start int64, amount uint) (map[int64]db.Message, error)
	allChats               func() ([]ed25519.PublicKey, error)
	addListener            func(fn func(e db.MessagePersistedEvent))
	getMessage             func(partner ed25519.PublicKey, messageID int64) (*db.Message, error)
	persistDAppMessage     func(partner ed25519.PublicKey, msg db.DAppMessage) error
}

func (s *testMessageStorage) PersistMessageToSend(partner ed25519.PublicKey, msg db.Message) error {
	return s.persistMessageToSend(partner, msg)
}

func (s *testMessageStorage) UpdateStatus(partner ed25519.PublicKey, msgID int64, newStatus db.Status) error {
	return s.updateStatus(partner, msgID, newStatus)
}

func (s *testMessageStorage) PersistReceivedMessage(partner ed25519.PublicKey, msg db.Message) error {
	return s.persistReceivedMessage(partner, msg)
}

func (s *testMessageStorage) Messages(partner ed25519.PublicKey, start int64, amount uint) (map[int64]db.Message, error) {
	return s.messages(partner, start, amount)
}

func (s *testMessageStorage) AllChats() ([]ed25519.PublicKey, error) {
	return s.allChats()
}

func (s *testMessageStorage) AddListener(fn func(e db.MessagePersistedEvent)) {
	s.addListener(fn)
}

func (s *testMessageStorage) GetMessage(partner ed25519.PublicKey, messageID int64) (*db.Message, error) {
	return s.getMessage(partner, messageID)
}

func (s *testMessageStorage) PersistDAppMessage(partner ed25519.PublicKey, msg db.DAppMessage) error {
	return s.persistDAppMessage(partner, msg)
}

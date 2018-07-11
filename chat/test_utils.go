package chat

import (
	db "github.com/Bit-Nation/panthalassa/db"
	bpb "github.com/Bit-Nation/protobuffers"
	x3dh "github.com/Bit-Nation/x3dh"
	ed25519 "golang.org/x/crypto/ed25519"
)

type testMessageStorage struct {
	persistMessage func(to ed25519.PublicKey, msg bpb.PlainChatMessage) error
	updateStatus   func(partner ed25519.PublicKey, msgID string, newStatus db.Status) error
}

type testSharedSecretStorage struct {
	hasAny      func(key ed25519.PublicKey) (bool, error)
	getYoungest func(key ed25519.PublicKey) (db.SharedSecret, error)
	put         func(key ed25519.PublicKey, proto x3dh.InitializedProtocol) error
}

type testBackend struct {
	fetchPreKeyBundle func(userIDPubKey ed25519.PublicKey) (x3dh.PreKeyBundle, error)
	submitMessage     func(msg bpb.ChatMessage) error
}

type testPreKeyBundle struct {
	identityKey     x3dh.PublicKey
	signedPreKey    x3dh.PublicKey
	preKeySignature []byte
	oneTimePreKey   *x3dh.PublicKey
	validSignature  bool
}

func (b testPreKeyBundle) IdentityKey() x3dh.PublicKey {
	return b.identityKey
}

func (b testPreKeyBundle) SignedPreKey() x3dh.PublicKey {
	return b.signedPreKey
}

func (b testPreKeyBundle) OneTimePreKey() *x3dh.PublicKey {
	return b.oneTimePreKey
}

func (b testPreKeyBundle) ValidSignature() bool {
	return b.validSignature
}

func (s *testSharedSecretStorage) HasAny(key ed25519.PublicKey) (bool, error) {
	return s.hasAny(key)
}

func (s *testSharedSecretStorage) GetYoungest(key ed25519.PublicKey) (db.SharedSecret, error) {
	return s.getYoungest(key)
}

func (s *testSharedSecretStorage) Put(key ed25519.PublicKey, proto x3dh.InitializedProtocol) error {
	return s.put(key, proto)
}

func (s *testMessageStorage) PersistMessage(to ed25519.PublicKey, msg bpb.PlainChatMessage) error {
	return s.persistMessage(to, msg)
}

func (s *testMessageStorage) UpdateStatus(partner ed25519.PublicKey, msgID string, newStatus db.Status) error {
	return s.updateStatus(partner, msgID, newStatus)
}

func (b *testBackend) FetchPreKeyBundle(userIDPubKey ed25519.PublicKey) (x3dh.PreKeyBundle, error) {
	return b.fetchPreKeyBundle(userIDPubKey)
}

func (b *testBackend) SubmitMessage(msg bpb.ChatMessage) error {
	return b.submitMessage(msg)
}

package chat

import (
	backend "github.com/Bit-Nation/panthalassa/backend"
	preKey "github.com/Bit-Nation/panthalassa/chat/prekey"
	db "github.com/Bit-Nation/panthalassa/db"
	km "github.com/Bit-Nation/panthalassa/keyManager"
	ks "github.com/Bit-Nation/panthalassa/keyStore"
	mnemonic "github.com/Bit-Nation/panthalassa/mnemonic"
	bpb "github.com/Bit-Nation/protobuffers"
	x3dh "github.com/Bit-Nation/x3dh"
	ed25519 "golang.org/x/crypto/ed25519"
)

type testMessageStorage struct {
	persistSentMessage     func(to ed25519.PublicKey, msg bpb.PlainChatMessage) error
	persistReceivedMessage func(partner ed25519.PublicKey, msg bpb.PlainChatMessage) error
	updateStatus           func(partner ed25519.PublicKey, msgID string, newStatus db.Status) error
}

type testSharedSecretStorage struct {
	hasAny               func(key ed25519.PublicKey) (bool, error)
	getYoungest          func(key ed25519.PublicKey) (*db.SharedSecret, error)
	put                  func(key ed25519.PublicKey, sharedSecret db.SharedSecret) error
	secretForChatInitMsg func(msg *bpb.ChatMessage) (*db.SharedSecret, error)
	accept               func(sharedSec *db.SharedSecret) error
	get                  func(key ed25519.PublicKey, sharedSecretID []byte) (*db.SharedSecret, error)
}

type testBackend struct {
	fetchPreKeyBundle func(userIDPubKey ed25519.PublicKey) (x3dh.PreKeyBundle, error)
	submitMessages    func(msg []*bpb.ChatMessage) error
	fetchSignedPreKey func(userIdPubKey ed25519.PublicKey) (preKey.PreKey, error)
	addRequestHandler func(backend.RequestHandler)
}

type testSignedPreKeyStore struct {
	getActive func() (*x3dh.KeyPair, error)
	put       func(signedPreKey x3dh.KeyPair) error
	get       func(publicKey x3dh.PublicKey) (*x3dh.PrivateKey, error)
	all       func() []*x3dh.KeyPair
}

type testUserStorage struct {
	getSignedPreKey func(idKey ed25519.PublicKey) (preKey.PreKey, error)
	hasSignedPreKey func(idKey ed25519.PublicKey) (bool, error)
	putSignedPreKey func(idKey ed25519.PublicKey, key preKey.PreKey) error
}

type testPreKeyBundle struct {
	identityKey     x3dh.PublicKey
	signedPreKey    x3dh.PublicKey
	preKeySignature []byte
	oneTimePreKey   *x3dh.PublicKey
	validSignature  func() (bool, error)
}

type testOneTimePreKeyStorage struct {
	cut   func(pubKey []byte) (*x3dh.PrivateKey, error)
	count func() (uint32, error)
	put   func(keyPair []x3dh.KeyPair) error
}

func (t *testOneTimePreKeyStorage) Cut(pubKey []byte) (*x3dh.PrivateKey, error) {
	return t.cut(pubKey)
}

func (t *testOneTimePreKeyStorage) Count() (uint32, error) {
	return t.count()
}

func (t *testOneTimePreKeyStorage) Put(keyPair []x3dh.KeyPair) error {
	return t.put(keyPair)
}

func (s *testSignedPreKeyStore) GetActive() (*x3dh.KeyPair, error) {
	return s.getActive()
}

func (s *testSignedPreKeyStore) Put(signedPreKey x3dh.KeyPair) error {
	return s.put(signedPreKey)
}

func (s *testSignedPreKeyStore) Get(publicKey x3dh.PublicKey) (*x3dh.PrivateKey, error) {
	return s.get(publicKey)
}

func (s *testSignedPreKeyStore) All() []*x3dh.KeyPair {
	return s.all()
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

func (b testPreKeyBundle) ValidSignature() (bool, error) {
	return b.validSignature()
}

func (s *testSharedSecretStorage) HasAny(key ed25519.PublicKey) (bool, error) {
	return s.hasAny(key)
}

func (s *testSharedSecretStorage) GetYoungest(key ed25519.PublicKey) (*db.SharedSecret, error) {
	return s.getYoungest(key)
}

func (s *testSharedSecretStorage) Put(key ed25519.PublicKey, sharedSecret db.SharedSecret) error {
	return s.put(key, sharedSecret)
}

func (s *testSharedSecretStorage) SecretForChatInitMsg(msg *bpb.ChatMessage) (*db.SharedSecret, error) {
	return s.secretForChatInitMsg(msg)
}

func (s *testSharedSecretStorage) Accept(sharedSec *db.SharedSecret) error {
	return s.accept(sharedSec)
}

func (s *testSharedSecretStorage) Get(key ed25519.PublicKey, sharedSecretID []byte) (*db.SharedSecret, error) {
	return s.get(key, sharedSecretID)
}

func (s *testMessageStorage) PersistSentMessage(partner ed25519.PublicKey, msg bpb.PlainChatMessage) error {
	return s.persistSentMessage(partner, msg)
}

func (s *testMessageStorage) UpdateStatus(partner ed25519.PublicKey, msgID string, newStatus db.Status) error {
	return s.updateStatus(partner, msgID, newStatus)
}

func (s *testMessageStorage) PersistReceivedMessage(partner ed25519.PublicKey, msg bpb.PlainChatMessage) error {
	return s.persistReceivedMessage(partner, msg)
}

func (b *testBackend) FetchPreKeyBundle(userIDPubKey ed25519.PublicKey) (x3dh.PreKeyBundle, error) {
	return b.fetchPreKeyBundle(userIDPubKey)
}

func (b *testBackend) SubmitMessages(messages []*bpb.ChatMessage) error {
	return b.submitMessages(messages)
}

func (b *testBackend) FetchSignedPreKey(userIdPubKey ed25519.PublicKey) (preKey.PreKey, error) {
	return b.fetchSignedPreKey(userIdPubKey)
}

func (b testBackend) AddRequestHandler(handler backend.RequestHandler) {
	b.addRequestHandler(handler)
}

func (s *testUserStorage) GetSignedPreKey(idKey ed25519.PublicKey) (preKey.PreKey, error) {
	return s.getSignedPreKey(idKey)
}

func (s *testUserStorage) HasSignedPreKey(idKey ed25519.PublicKey) (bool, error) {
	return s.hasSignedPreKey(idKey)
}

func (s *testUserStorage) PutSignedPreKey(idKey ed25519.PublicKey, key preKey.PreKey) error {
	return s.putSignedPreKey(idKey, key)
}

func createKeyManager() *km.KeyManager {

	mne, err := mnemonic.New()
	if err != nil {
		panic(err)
	}

	keyStore, err := ks.NewFromMnemonic(mne)
	if err != nil {
		panic(err)
	}

	return km.CreateFromKeyStore(keyStore)

}

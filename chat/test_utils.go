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
	persistMessageToSend   func(to ed25519.PublicKey, msg db.Message) error
	persistReceivedMessage func(partner ed25519.PublicKey, msg db.Message) error
	updateStatus           func(partner ed25519.PublicKey, msgID int64, newStatus db.Status) error
	messages               func(partner ed25519.PublicKey, start int64, amount uint) ([]db.Message, error)
	allChats               func() ([]ed25519.PublicKey, error)
	addListener            func(fn func(e db.MessagePersistedEvent))
	getMessage             func(partner ed25519.PublicKey, messageID int64) (*db.Message, error)
	persistDAppMessage     func(partner ed25519.PublicKey, msg db.DAppMessage) error
}

type testSharedSecretStorage struct {
	hasAny      func(key ed25519.PublicKey) (bool, error)
	getYoungest func(key ed25519.PublicKey) (*db.SharedSecret, error)
	put         func(sharedSecret db.SharedSecret) error
	accept      func(sharedSec db.SharedSecret) error
	get         func(key ed25519.PublicKey, sharedSecretID []byte) (*db.SharedSecret, error)
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
	all       func() ([]*x3dh.KeyPair, error)
}

type testUserStorage struct {
	getSignedPreKey func(idKey ed25519.PublicKey) (*preKey.PreKey, error)
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

func (s *testSignedPreKeyStore) Put(signedPreKey x3dh.KeyPair) error {
	return s.put(signedPreKey)
}

func (s *testSignedPreKeyStore) Get(publicKey x3dh.PublicKey) (*x3dh.PrivateKey, error) {
	return s.get(publicKey)
}

func (s *testSignedPreKeyStore) All() ([]*x3dh.KeyPair, error) {
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

func (s *testSharedSecretStorage) Put(sharedSecret db.SharedSecret) error {
	return s.put(sharedSecret)
}

func (s *testSharedSecretStorage) Accept(sharedSec db.SharedSecret) error {
	return s.accept(sharedSec)
}

func (s *testSharedSecretStorage) Get(key ed25519.PublicKey, sharedSecretID []byte) (*db.SharedSecret, error) {
	return s.get(key, sharedSecretID)
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

func (s *testMessageStorage) Messages(partner ed25519.PublicKey, start int64, amount uint) ([]db.Message, error) {
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

func (b *testBackend) FetchPreKeyBundle(userIDPubKey ed25519.PublicKey) (x3dh.PreKeyBundle, error) {
	return b.fetchPreKeyBundle(userIDPubKey)
}

func (b *testBackend) SubmitMessages(messages []*bpb.ChatMessage) error {
	return b.submitMessages(messages)
}

func (b *testBackend) FetchSignedPreKey(userIdPubKey ed25519.PublicKey) (preKey.PreKey, error) {
	return b.fetchSignedPreKey(userIdPubKey)
}

func (b *testBackend) AddRequestHandler(handler backend.RequestHandler) {
	b.addRequestHandler(handler)
}

func (b *testBackend) Close() error {
	return nil
}

func (s *testUserStorage) GetSignedPreKey(idKey ed25519.PublicKey) (*preKey.PreKey, error) {
	return s.getSignedPreKey(idKey)
}

func (s *testUserStorage) PutSignedPreKey(idKey ed25519.PublicKey, key preKey.PreKey) error {
	return s.putSignedPreKey(idKey, key)
}

func (s *testMessageStorage) PersistDAppMessage(partner ed25519.PublicKey, msg db.DAppMessage) error {
	return s.persistDAppMessage(partner, msg)
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

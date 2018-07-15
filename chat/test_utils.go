package chat

import (
	preKey "github.com/Bit-Nation/panthalassa/chat/prekey"
	db "github.com/Bit-Nation/panthalassa/db"
	km "github.com/Bit-Nation/panthalassa/keyManager"
	ks "github.com/Bit-Nation/panthalassa/keyStore"
	mnemonic "github.com/Bit-Nation/panthalassa/mnemonic"
	bpb "github.com/Bit-Nation/protobuffers"
	x3dh "github.com/Bit-Nation/x3dh"
	dr "github.com/tiabc/doubleratchet"
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
	fetchSignedPreKey func(userIdPubKey ed25519.PublicKey) (preKey.PreKey, error)
}

type testSignedPreKeyStore struct {
	hasActive func() (bool, error)
	getActive func() (x3dh.KeyPair, error)
	put       func(signedPreKey x3dh.KeyPair) error
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
	validSignature  bool
}

type drDhPair struct {
	pub  x3dh.PublicKey
	priv x3dh.PrivateKey
}

func (p drDhPair) PrivateKey() dr.Key {
	var k dr.Key
	copy(k[:], p.priv[:])
	return k
}

func (p drDhPair) PublicKey() dr.Key {
	var k dr.Key
	copy(k[:], p.pub[:])
	return k
}

func (s *testSignedPreKeyStore) HasActive() (bool, error) {
	return s.hasActive()
}

func (s *testSignedPreKeyStore) GetActive() (x3dh.KeyPair, error) {
	return s.getActive()
}

func (s *testSignedPreKeyStore) Put(signedPreKey x3dh.KeyPair) error {
	return s.put(signedPreKey)
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

func (b *testBackend) FetchSignedPreKey(userIdPubKey ed25519.PublicKey) (preKey.PreKey, error) {
	return b.fetchSignedPreKey(userIdPubKey)
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

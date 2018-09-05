package backend

import (
	bpb "github.com/Bit-Nation/protobuffers"
	x3dh "github.com/Bit-Nation/x3dh"
)

type testTransport struct {
	send func(msg *bpb.BackendMessage) error
	// this will be the callback that should be called on a message
	// from the transport
	nextMessage func() (*bpb.BackendMessage, error)
}

func (t *testTransport) Send(msg *bpb.BackendMessage) error {
	return t.send(msg)
}

func (t *testTransport) NextMessage() (*bpb.BackendMessage, error) {
	return t.nextMessage()
}

func (t *testTransport) Close() error {
	return nil
}

func (t *testTransport) Start() error {
	return nil
}

type testSignedPreKeyStore struct {
	getActive func() (*x3dh.KeyPair, error)
	put       func(signedPreKey x3dh.KeyPair) error
	get       func(publicKey x3dh.PublicKey) (*x3dh.PrivateKey, error)
	all       func() ([]*x3dh.KeyPair, error)
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

func (s *testSignedPreKeyStore) All() ([]*x3dh.KeyPair, error) {
	return s.all()
}

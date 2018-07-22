package backend

import (
	"errors"
	"time"

	preKey "github.com/Bit-Nation/panthalassa/chat/prekey"
	bpb "github.com/Bit-Nation/protobuffers"
	x3dh "github.com/Bit-Nation/x3dh"
	ed25519 "golang.org/x/crypto/ed25519"
)

// fetch pre key bundle from backend
func (b *Backend) FetchPreKeyBundle(userIDPubKey ed25519.PublicKey) (x3dh.PreKeyBundle, error) {

	// request pre key bundle
	resp, err := b.request(bpb.BackendMessage_Request{
		PreKeyBundle: userIDPubKey,
	}, time.Second*4)
	if err != nil {
		return &PreKeyBundle{}, err
	}

	// unmarshal protobuf
	bundle, err := PreKeyBundleFromProto(userIDPubKey, resp.PreKeyBundle)
	if err != nil {
		return &PreKeyBundle{}, err
	}

	// validate signatures
	valid, err := bundle.ValidSignature()
	if err != nil {
		return &PreKeyBundle{}, err
	}

	// exit if invalid signature
	if !valid {
		return &PreKeyBundle{}, errors.New("invalid pre key bundle signatures")
	}

	// parse pre key bundle
	return bundle, nil

}

// submit messages
func (b *Backend) SubmitMessages(messages []*bpb.ChatMessage) error {
	_, err := b.request(bpb.BackendMessage_Request{Messages: messages}, time.Second*20)
	return err
}

// fetch signed pre key of person
func (b *Backend) FetchSignedPreKey(userIdPubKey ed25519.PublicKey) (preKey.PreKey, error) {
	// request signed pre key of user from backend
	resp, err := b.request(bpb.BackendMessage_Request{SignedPreKey: userIdPubKey}, time.Second*4)
	if err != nil {
		return preKey.PreKey{}, err
	}
	// unmarshal signed pre key
	pk, err := preKey.FromProtoBuf(*resp.SignedPreKey)
	if err != nil {
		return preKey.PreKey{}, err
	}
	// verify signature of signed pre key
	valid, err := pk.VerifySignature(userIdPubKey)
	if err != nil {
		return preKey.PreKey{}, err
	}
	if !valid {
		return preKey.PreKey{}, errors.New("invalid signed pre key signature")
	}
	return pk, nil
}

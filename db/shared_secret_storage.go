package db

import (
	"time"

	bpb "github.com/Bit-Nation/protobuffers"
	x3dh "github.com/Bit-Nation/x3dh"
	ed25519 "golang.org/x/crypto/ed25519"
)

type SharedSecret struct {
	X3dhSS                x3dh.SharedSecret
	Accepted              bool
	CreatedAt             time.Time
	DestroyAt             *time.Time
	EphemeralKey          x3dh.PublicKey
	EphemeralKeySignature []byte
	UsedSignedPreKey      x3dh.PublicKey
	UsedOneTimePreKey     *x3dh.PublicKey
	ID                    []byte
	IDInitParams          []byte
}

type SharedSecretStorage interface {
	HasAny(key ed25519.PublicKey) (bool, error)
	// must return an error if no shared secret found
	GetYoungest(key ed25519.PublicKey) (*SharedSecret, error)
	Put(chatPartner ed25519.PublicKey, ss SharedSecret) error
	// check if a secret for a chat initialization message exists
	SecretForChatInitMsg(msg *bpb.ChatMessage) (*SharedSecret, error)
	// accept will mark the given shared secret as accepted
	// and will set a destroy date for all other shared secrets
	Accept(sharedSec *SharedSecret) error
	// get sender public key and shared secret id
	Get(key ed25519.PublicKey, sharedSecretID []byte) (*SharedSecret, error)
}

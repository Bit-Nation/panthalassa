package backend

import (
	preKey "github.com/Bit-Nation/panthalassa/chat/prekey"
	bpb "github.com/Bit-Nation/protobuffers"
	x3dh "github.com/Bit-Nation/x3dh"
	ed25519 "golang.org/x/crypto/ed25519"
)

// represents the source were we get things
// like the pre key bundles from
type Backend interface {
	FetchPreKeyBundle(userIDPubKey ed25519.PublicKey) (x3dh.PreKeyBundle, error)
	SubmitMessage(msg bpb.ChatMessage) error
	FetchSignedPreKey(userIdPubKey ed25519.PublicKey) (preKey.PreKey, error)
}

package db

import (
	bpb "github.com/Bit-Nation/protobuffers"
	ed25519 "golang.org/x/crypto/ed25519"
)

var (
	messagesBucketName       = []byte("chat")
	messageStorageBucketName = []byte("messages")
	orderBucketName          = []byte("order")
)

// message status
type Status uint

const (
	StatusSent           Status = 100
	StatusFailedToSend   Status = 200
	StatusDelivered      Status = 300
	StatusFailedToHandle Status = 400
)

type ChatMessageStorage interface {
	PersistSentMessage(partner ed25519.PublicKey, msg bpb.PlainChatMessage) error
	PersistReceivedMessage(partner ed25519.PublicKey, msg bpb.PlainChatMessage) error
	UpdateStatus(partner ed25519.PublicKey, msgID string, newStatus Status) error
}

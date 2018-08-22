package db

import (
	"errors"

	preKey "github.com/Bit-Nation/panthalassa/chat/prekey"
	pb "github.com/Bit-Nation/protobuffers"
	bolt "github.com/coreos/bbolt"
	proto "github.com/gogo/protobuf/proto"
	ed25519 "golang.org/x/crypto/ed25519"
)

var (
	userStorageBucketName = []byte("user_metadata_storage")
	signedPreKeyName      = []byte("signed_pre_key")
)

// user storage store meta data about users
type UserStorage interface {
	// don't forget to verify the signature when implementing this
	GetSignedPreKey(idKey ed25519.PublicKey) (*preKey.PreKey, error)
	PutSignedPreKey(idKey ed25519.PublicKey, key preKey.PreKey) error
}

type BoltUserStorage struct {
	db *bolt.DB
}

func NewBoltUserStorage(db *bolt.DB) *BoltUserStorage {
	return &BoltUserStorage{db: db}
}

func (s *BoltUserStorage) GetSignedPreKey(idKey ed25519.PublicKey) (*preKey.PreKey, error) {
	signedPreKey := new(preKey.PreKey)
	signedPreKey = nil
	err := s.db.View(func(tx *bolt.Tx) error {

		// fetch user storage bucket
		userStorageBucket := tx.Bucket(userStorageBucketName)
		if userStorageBucket == nil {
			return nil
		}

		// fetch user bucket
		userBucket := userStorageBucket.Bucket(idKey)
		if userBucket == nil {
			return nil
		}

		// fetch signed pre key
		rawProtoSignedPreKey := userBucket.Get(signedPreKeyName)
		if rawProtoSignedPreKey == nil {
			return nil
		}

		// unmarshal pre key
		protoSignedPreKey := pb.PreKey{}
		if err := proto.Unmarshal(rawProtoSignedPreKey, &protoSignedPreKey); err != nil {
			return err
		}

		// parse proto pre key to pre key
		spk, err := preKey.FromProtoBuf(protoSignedPreKey)
		if err != nil {
			return err
		}
		signedPreKey = &spk

		return nil
	})
	return signedPreKey, err
}

func (s *BoltUserStorage) PutSignedPreKey(idKey ed25519.PublicKey, key preKey.PreKey) error {

	// validate signed pre key
	validSig, err := key.VerifySignature(idKey)
	if err != nil {
		return err
	}
	if !validSig {
		return errors.New("got invalid signature for pre key")
	}

	return s.db.Update(func(tx *bolt.Tx) error {

		// fetch user storage bucket
		userStorageBucket, err := tx.CreateBucketIfNotExists(userStorageBucketName)
		if err != nil {
			return nil
		}

		// fetch user bucket
		userBucket, err := userStorageBucket.CreateBucketIfNotExists(idKey)
		if err != nil {
			return nil
		}

		// proto key
		protoKey, err := key.ToProtobuf()
		if err != nil {
			return err
		}

		// marshaled protobuf key
		rawProtoKey, err := proto.Marshal(&protoKey)
		if err != nil {
			return err
		}

		return userBucket.Put(signedPreKeyName, rawProtoKey)

	})
}

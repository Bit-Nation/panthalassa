package prekey

import (
	"bytes"
	"encoding/hex"
	"errors"
	"time"

	"encoding/binary"
	km "github.com/Bit-Nation/panthalassa/keyManager"
	pb "github.com/Bit-Nation/protobuffers"
	x3dh "github.com/Bit-Nation/x3dh"
	mh "github.com/multiformats/go-multihash"
	ed25519 "golang.org/x/crypto/ed25519"
)

var (
	InvalidIdentityKey = errors.New("invalid identity key - expected length to be 32 bytes")
)

type PreKey struct {
	x3dh.KeyPair
	identityPublicKey [32]byte
	signature         []byte
	time              time.Time
}

func (p *PreKey) hash() (mh.Multihash, error) {

	// make sure it has the right size
	if p.identityPublicKey == [32]byte{} {
		return nil, errors.New("got invalid identity key public key")
	}

	if p.PublicKey == [32]byte{} {
		return nil, errors.New("got invalid pre key public key")
	}

	b := bytes.NewBuffer(p.PublicKey[:])
	b.Write(p.identityPublicKey[:])

	timeStamp := make([]byte, binary.MaxVarintLen64)
	n := binary.PutVarint(timeStamp, p.time.Unix())

	b.Write(timeStamp[:n])
	return mh.Sum(b.Bytes(), mh.SHA3_256, -1)
}

// sign prekey public key
func (p *PreKey) Sign(km km.KeyManager) error {

	idPubKey, err := km.IdentityPublicKey()
	if err != nil {
		return err
	}

	rawIdPubKey, err := hex.DecodeString(idPubKey)
	if err != nil {
		return err
	}

	if len(rawIdPubKey) != 32 {
		return InvalidIdentityKey
	}

	if p.PublicKey == [32]byte{} {
		return errors.New("got invalid pre key public key")
	}

	copy(p.identityPublicKey[:], rawIdPubKey[:32])

	hash, err := p.hash()
	if err != nil {
		return err
	}

	// attach the signature to the pre key
	p.signature, err = km.IdentitySign(hash)
	if err != nil {
		return err
	}

	return nil

}

func (p PreKey) VerifySignature(publicKey ed25519.PublicKey) (bool, error) {
	hash, err := p.hash()
	if err != nil {
		return false, err
	}
	if len(p.identityPublicKey) != 32 {
		return false, InvalidIdentityKey
	}
	return ed25519.Verify(publicKey, hash, p.signature), nil
}

func (p PreKey) ToProtobuf() (pb.PreKey, error) {
	return pb.PreKey{
		Key:                  p.PublicKey[:],
		IdentityKey:          p.identityPublicKey[:],
		IdentityKeySignature: p.signature,
		TimeStamp:            p.time.Unix(),
	}, nil
}

func FromProtoBuf(preKey pb.PreKey) (PreKey, error) {
	if len(preKey.Key) != 32 {
		return PreKey{}, errors.New("pre key public key is invalid - expected be 32 bytes long")
	}
	if len(preKey.IdentityKey) != 32 {
		return PreKey{}, InvalidIdentityKey
	}
	p := PreKey{
		signature: preKey.IdentityKeySignature,
		time:      time.Unix(preKey.TimeStamp, 0),
	}
	copy(p.PublicKey[:], preKey.Key[:32])
	copy(p.identityPublicKey[:], preKey.IdentityKey[:32])
	return p, nil
}

// check if pre key is older than given date
func (p PreKey) OlderThen(past time.Duration) bool {
	return p.time.After(time.Now().Truncate(past))
}

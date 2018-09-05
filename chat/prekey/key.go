package prekey

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"time"

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
	IdentityPublicKey [32]byte
	Signature         []byte
	// unix timestamp in second
	Time int64
}

func (p *PreKey) hash() (mh.Multihash, error) {
	// make sure it has the right size
	if p.IdentityPublicKey == [32]byte{} {
		return nil, errors.New("got invalid identity key public key")
	}

	if p.PublicKey == [32]byte{} {
		return nil, errors.New("got invalid pre key public key")
	}

	b := bytes.NewBuffer(p.PublicKey[:])
	b.Write(p.IdentityPublicKey[:])

	timeStamp := make([]byte, 8)
	n := binary.PutVarint(timeStamp, p.Time)

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

	copy(p.IdentityPublicKey[:], rawIdPubKey[:32])

	hash, err := p.hash()
	if err != nil {
		return err
	}

	// attach the signature to the pre key
	p.Signature, err = km.IdentitySign(hash)
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
	if len(p.IdentityPublicKey) != 32 {
		return false, InvalidIdentityKey
	}
	return ed25519.Verify(publicKey, hash, p.Signature), nil
}

func (p PreKey) ToProtobuf() (pb.PreKey, error) {
	return pb.PreKey{
		Key:                  p.PublicKey[:],
		IdentityKey:          p.IdentityPublicKey[:],
		IdentityKeySignature: p.Signature,
		TimeStamp:            p.Time,
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
		Signature: preKey.IdentityKeySignature,
		Time:      preKey.TimeStamp,
	}
	copy(p.PublicKey[:], preKey.Key[:32])
	copy(p.IdentityPublicKey[:], preKey.IdentityKey[:32])
	return p, nil
}

// check if pre key is older than given date
func (p PreKey) OlderThan(past time.Duration) bool {
	t := time.Unix(p.Time, 0)
	return t.After(time.Now().Truncate(past))
}

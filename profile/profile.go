package profile

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"time"

	km "github.com/Bit-Nation/panthalassa/keyManager"
	pb "github.com/Bit-Nation/protobuffers"
	x3dh "github.com/Bit-Nation/x3dh"
	common "github.com/ethereum/go-ethereum/common"
	crypto "github.com/ethereum/go-ethereum/crypto"
	mh "github.com/multiformats/go-multihash"
	ed25519 "golang.org/x/crypto/ed25519"
)

const (
	profileVersion = 2
	TimeFormat     = time.UnixDate
)

var (
	IdentityPublicKeyTooShort = errors.New("identity public key too short - must have 32 bytes")
	EthereumPublicKeyTooShort = errors.New("ethereum public key too short - must have 33 bytes")
	ChatIDKeyTooShort         = errors.New("chat identity key too short - must have 32 bytes")
)

type Information struct {
	Name           string
	Location       string
	Image          string
	IdentityPubKey []byte
	EthereumPubKey []byte
	ChatIDKey      x3dh.PublicKey
	Timestamp      time.Time
	Version        uint8
}

type Signatures struct {
	IdentityKey []byte
	EthereumKey []byte
}

type Profile struct {
	Information Information
	Signatures  Signatures
}

func ProtobufToProfile(prof *pb.Profile) (*Profile, error) {

	if len(prof.IdentityPubKey) != 32 {
		return nil, IdentityPublicKeyTooShort
	}

	if len(prof.EthereumPubKey) != 33 {
		return nil, EthereumPublicKeyTooShort
	}

	if len(prof.ChatIdentityPubKey) != 32 {
		return nil, ChatIDKeyTooShort
	}

	chatIDPubKey := x3dh.PublicKey{}
	copy(chatIDPubKey[:], prof.ChatIdentityPubKey[:32])

	timeStamp, err := time.Parse(TimeFormat, prof.Timestamp)
	if err != nil {
		return nil, err
	}

	return &Profile{
		Information: Information{
			Name:           prof.Name,
			Location:       prof.Location,
			Image:          prof.Image,
			IdentityPubKey: prof.IdentityPubKey,
			EthereumPubKey: prof.EthereumPubKey,
			ChatIDKey:      chatIDPubKey,
			Timestamp:      timeStamp,
			Version:        uint8(prof.Version),
		},
		Signatures: Signatures{
			IdentityKey: prof.IdentityKeySignature,
			EthereumKey: prof.EthereumKeySignature,
		},
	}, nil

}

// hash all information about a profile
func (p *Profile) Hash() (mh.Multihash, error) {
	b := bytes.NewBuffer([]byte(p.Information.Name))

	toWrite := [][]byte{
		[]byte(p.Information.Location),
		[]byte(p.Information.Image),
		p.Information.IdentityPubKey,
		p.Information.EthereumPubKey,
		p.Information.ChatIDKey[:],
		[]byte(p.Information.Timestamp.Format(TimeFormat)),
	}

	for _, tw := range toWrite {
		if _, err := b.Write(tw); err != nil {
			return nil, err
		}
	}

	if err := b.WriteByte(p.Information.Version); err != nil {
		return nil, err
	}

	return mh.Sum(b.Bytes(), mh.SHA2_256, -1)

}

func (p *Profile) IdentityPublicKey() (ed25519.PublicKey, error) {
	if len(p.Information.IdentityPubKey) != 32 {
		return nil, IdentityPublicKeyTooShort
	}
	return p.Information.IdentityPubKey, nil
}

func (p *Profile) EthereumPublicKey() ([]byte, error) {
	if len(p.Information.EthereumPubKey) != 33 {
		return nil, EthereumPublicKeyTooShort
	}
	return p.Information.EthereumPubKey, nil
}

// Validate signature based on identity pubic key
func (p *Profile) ValidIdentitySignature(msg, signature []byte) (bool, error) {
	idPubKey, err := p.IdentityPublicKey()
	if err != nil {
		return false, err
	}
	return ed25519.Verify(idPubKey, msg, signature), nil
}

// Validate signature based on the ethereum public key
func (p *Profile) ValidEthereumSignature(msg [32]byte, signature [65]byte) (bool, error) {
	ethPubKey, err := p.EthereumPublicKey()
	if err != nil {
		return false, err
	}
	// an the signature contain's [R || S || V] - but we only need R and S
	return crypto.VerifySignature(ethPubKey, msg[:], signature[:64]), nil
}

// check if signatures of profile are correct
func (p Profile) SignaturesValid() (bool, error) {

	h, err := p.Hash()
	if err != nil {
		return false, err
	}

	// check identity key signature
	valid, err := p.ValidIdentitySignature(h, p.Signatures.IdentityKey)
	if err != nil {
		return false, err
	}
	if !valid {
		return false, errors.New("invalid identity signature")
	}

	// check ethereum key signature
	var sig [65]byte
	copy(sig[:], p.Signatures.EthereumKey[:65])
	valid, err = p.ValidEthereumSignature(sha256.Sum256(h), sig)
	if err != nil {
		return false, err
	}
	if !valid {
		return false, errors.New("invalid ethereum signature")
	}

	return true, nil

}

// Convert profile to protobuf
func (p *Profile) ToProtobuf() (*pb.Profile, error) {

	if len(p.Information.IdentityPubKey) != 32 {
		return nil, IdentityPublicKeyTooShort
	}

	if len(p.Information.EthereumPubKey) != 33 {
		return nil, EthereumPublicKeyTooShort
	}

	if len(p.Information.ChatIDKey) != 32 {
		return nil, ChatIDKeyTooShort
	}

	pp := pb.Profile{}
	pp.Name = p.Information.Name
	pp.Location = p.Information.Location
	pp.Image = p.Information.Image
	pp.IdentityPubKey = p.Information.IdentityPubKey
	pp.EthereumPubKey = p.Information.EthereumPubKey
	pp.ChatIdentityPubKey = p.Information.ChatIDKey[:]
	pp.Timestamp = p.Information.Timestamp.Format(TimeFormat)
	pp.Version = uint32(p.Information.Version)

	pp.IdentityKeySignature = p.Signatures.IdentityKey
	pp.EthereumKeySignature = p.Signatures.EthereumKey

	return &pp, nil

}

// get ethereum address based on profile ethereum public key
func (p *Profile) EthereumAddress() (common.Address, error) {

	// validate profile signature
	valid, err := p.SignaturesValid()
	if err != nil {
		return common.Address{}, err
	}
	if !valid {
		return common.Address{}, errors.New("invalid profile signature")
	}

	// raw ethereum key
	ethPubKey, err := p.EthereumPublicKey()
	if err != nil {
		return common.Address{}, err
	}

	// decompressed ethereum key
	decPubKey, err := crypto.DecompressPubkey(ethPubKey)
	if err != nil {
		return common.Address{}, err
	}

	return crypto.PubkeyToAddress(*decPubKey), nil
}

// sign the metadata with identity and ethereum key
func SignProfile(name, location, image string, km km.KeyManager) (*Profile, error) {

	// chat id key pair
	chatIDKeys, err := km.ChatIdKeyPair()
	if err != nil {
		return nil, err
	}

	// identity pub key
	idPubKeyStr, err := km.IdentityPublicKey()
	if err != nil {
		return nil, err
	}
	idPubKeyRaw, err := hex.DecodeString(idPubKeyStr)
	if err != nil {
		return nil, err
	}

	// ethereum pub key
	ethPubKeyStr, err := km.GetEthereumPublicKey()
	if err != nil {
		return nil, err
	}
	ethPubKeyRaw, err := hex.DecodeString(ethPubKeyStr)
	if err != nil {
		return nil, err
	}

	p := Profile{
		Information: Information{
			Name:           name,
			Location:       location,
			Image:          image,
			ChatIDKey:      chatIDKeys.PublicKey,
			Timestamp:      time.Now(),
			Version:        profileVersion,
			IdentityPubKey: idPubKeyRaw,
			EthereumPubKey: ethPubKeyRaw,
		},
	}

	hash, err := p.Hash()
	if err != nil {
		return nil, err
	}

	p.Signatures.IdentityKey, err = km.IdentitySign(hash)
	if err != nil {
		return nil, err
	}

	// sadly we have to hash the hash again to have 32 bytes
	// multihash adds a few bytes to identify what type of hash has been used
	p.Signatures.EthereumKey, err = km.EthereumSign(sha256.Sum256(hash))

	return &p, nil

}

// sign's a profile without the need to start panthalassa
func SignWithKeyManagerStore(name, location, image string, keyManagerStore km.Store, password string) (*Profile, error) {

	keyManager, err := km.OpenWithPassword(keyManagerStore, password)
	if err != nil {
		return &Profile{}, err
	}

	return SignProfile(name, location, image, *keyManager)

}

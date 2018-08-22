package backend

import (
	preKey "github.com/Bit-Nation/panthalassa/chat/prekey"
	profile "github.com/Bit-Nation/panthalassa/profile"
	bpb "github.com/Bit-Nation/protobuffers"
	x3dh "github.com/Bit-Nation/x3dh"
	ed25519 "golang.org/x/crypto/ed25519"
)

func PreKeyBundleFromProto(pubKey ed25519.PublicKey, protoPreKeyBundle *bpb.BackendMessage_PreKeyBundle) (*PreKeyBundle, error) {

	// pre key bundle owner
	preKeyBundleOwner, err := profile.ProtobufToProfile(protoPreKeyBundle.Profile)
	if err != nil {
		return nil, err
	}

	// signed pre key
	signedPreKey, err := preKey.FromProtoBuf(*protoPreKeyBundle.SignedPreKey)
	if err != nil {
		return nil, err
	}

	// pre key bundle
	pkb := &PreKeyBundle{
		expectedBundleOwner: pubKey,
		bundleOwner:         *preKeyBundleOwner,
		signedPreKey:        signedPreKey,
	}

	// one time pre key
	if protoPreKeyBundle.OneTimePreKey != nil {
		oneTimePreKey, err := preKey.FromProtoBuf(*protoPreKeyBundle.OneTimePreKey)
		if err == nil {
			pkb.oneTimePreKey = &oneTimePreKey
		} else {
			logger.Error("failed to unmarshal OneTimePreKey, error: ", err)
		}
	}

	return pkb, nil
}

type PreKeyBundle struct {
	expectedBundleOwner ed25519.PublicKey
	bundleOwner         profile.Profile
	signedPreKey        preKey.PreKey
	oneTimePreKey       *preKey.PreKey
}

func (b *PreKeyBundle) IdentityKey() x3dh.PublicKey {
	return b.bundleOwner.Information.ChatIDKey
}

func (b *PreKeyBundle) SignedPreKey() x3dh.PublicKey {
	return b.signedPreKey.PublicKey
}

func (b *PreKeyBundle) OneTimePreKey() *x3dh.PublicKey {
	if b.oneTimePreKey != nil {
		return &b.oneTimePreKey.PublicKey
	}
	return nil
}

func (b *PreKeyBundle) ValidSignature() (bool, error) {

	// verify profile signatures
	valid, err := b.bundleOwner.SignaturesValid()
	if err != nil || !valid {
		return valid, err
	}

	// verify signed pre key
	valid, err = b.signedPreKey.VerifySignature(b.expectedBundleOwner)
	if err != nil || !valid {
		return valid, err
	}

	// verify one time pre key
	if b.oneTimePreKey != nil {
		valid, err := b.oneTimePreKey.VerifySignature(b.expectedBundleOwner)
		if err != nil || !valid {
			return valid, err
		}
	}

	return true, nil

}

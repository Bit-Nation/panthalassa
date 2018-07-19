package chat

import (
	bpb "github.com/Bit-Nation/protobuffers"
	x3dh "github.com/Bit-Nation/x3dh"
)

func PreKeyBundleFromProto(bpb *bpb.BackendMessage_PreKeyBundle) (*PreKeyBundle, error) {

}

type PreKeyBundle struct {
}

func (b *PreKeyBundle) IdentityKey() x3dh.PublicKey {

}

func (b *PreKeyBundle) SignedPreKey() x3dh.PublicKey {

}

func (b *PreKeyBundle) OneTimePreKey() *x3dh.PublicKey {

}

func (b *PreKeyBundle) ValidSignature() bool {

}

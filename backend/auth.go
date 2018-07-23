package backend

import (
	"crypto/rand"
	"encoding/hex"
	"errors"

	bpb "github.com/Bit-Nation/protobuffers"
)

// authentication request handler
func (b *Backend) auth(req *bpb.BackendMessage_Request) (*bpb.BackendMessage_Response, error) {

	auth := req.Auth

	if auth != nil {

		if len(auth.ToSign) != 4 {
			return nil, errors.New("got invalid amount of toSign bytes")
		}

		myBytes := make([]byte, 4)
		_, err := rand.Read(myBytes)
		if err != nil {
			return nil, err
		}

		toSign := append(auth.ToSign, myBytes...)

		// sign data
		signature, err := b.km.IdentitySign(toSign)
		if err != nil {
			return nil, errors.New("failed to sign data")
		}

		// get identity public key
		idPubStr, err := b.km.IdentityPublicKey()
		if err != nil {
			return nil, err
		}
		rawIDKey, err := hex.DecodeString(idPubStr)
		if err != nil {
			return nil, err
		}

		return &bpb.BackendMessage_Response{
			Auth: &bpb.BackendMessage_Auth{
				Signature:         signature,
				IdentityPublicKey: rawIDKey,
				ToSign:            toSign,
			},
		}, nil
	}

	return nil, nil

}

package backend

import (
	"encoding/hex"
	"errors"

	bpb "github.com/Bit-Nation/protobuffers"
)

// authentication request handler
func (b *ServerBackend) auth() RequestHandler {

	fn := func(req *bpb.BackendMessage_Request) (*bpb.BackendMessage_Response, error) {

		auth := req.Auth

		if auth != nil {
			// sign data
			signature, err := b.km.IdentitySign(auth.ToSign)
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
					ToSign:            auth.ToSign,
				},
			}, nil
		}

		return nil, nil
	}

	return &fn

}

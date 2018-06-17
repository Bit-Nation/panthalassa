package registry

import (
	"github.com/Bit-Nation/panthalassa/dapp"
)

// the client is the "thing" that's using
// the DApps + Registry. It's responsible for
// saving the received DApp
type Client interface {
	HandleReceivedDApp(dApp dapp.JsonRepresentation) error
}

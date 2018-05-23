package panthalassa

import (
	"context"
	"time"

	pid "github.com/libp2p/go-libp2p-peer"
)

func SearchContacts(p *Panthalassa) {

	go func() {

		logger.Info("start the search for contract's")

		for {

			peers := p.mesh.Host.Peerstore().Peers()
			for _, peer := range peers {
				go func(peer pid.ID) {

					isContact, err := p.mesh.Host.Peerstore().Get(peer, "contact")
					if err != nil {
						if err.Error() == "item not found" {
							return
						}
						logger.Error(err)
						return
					}
					if isContact == true {
						pubKey, err := peer.ExtractPublicKey()
						if err != nil {
							logger.Error(pubKey)
							return
						}
						contactKey, err := p.mesh.RendezvousKey.Profile(pubKey)
						if err != nil {
							logger.Error(err)
							return
						}
						profile, err := p.mesh.Dht.GetValue(context.Background(), contactKey.String())
						if err != nil {
							logger.Error(err)
							return
						}
						// @todo verify user profile
						p.mesh.Dht.PutValue(context.Background(), contactKey.String(), profile)
					}

				}(peer)
			}

			time.Sleep(time.Second * 5)

		}

		logger.Info("stop updating user profiles")

	}()

}

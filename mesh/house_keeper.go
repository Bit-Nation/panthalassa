package mesh

import (
	"context"
	"time"

	host "github.com/libp2p/go-libp2p-host"
	"github.com/libp2p/go-libp2p-kad-dht"
	"sync"
)

//This service is responsible for searching for peer's
//@todo peer's that are just "interested" in pangea have a very low priority - we need to make sure that we prefere our friends over them
func SearchPangeaPeers(h host.Host, d *dht.IpfsDHT, rk RendezvousKey, errReporter chan<- error) {

	go func() {

		for {

			logger.Info("Search for pangea peer's")

			if len(h.Network().Peers()) == 25 {

			}

			k, err := rk.Derive(pangea)
			if err != nil {
				errReporter <- err
			}

			peers, err := d.FindProviders(context.Background(), k)
			if err != nil {
				errReporter <- err
			}

			wg := sync.WaitGroup{}

			for _, peerInfo := range peers {
				wg.Add(1)
				go func() {
					if err := h.Connect(context.Background(), peerInfo); err != nil {
						logger.Error(err)
					}
					wg.Done()
				}()
			}

			wg.Wait()

			time.Sleep(time.Second * 5)

		}

	}()

}

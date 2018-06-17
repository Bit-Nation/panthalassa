package registry

import (
	"context"
	"encoding/hex"
	"errors"
	"io/ioutil"
	"sync"

	dapp "github.com/Bit-Nation/panthalassa/dapp"
	module "github.com/Bit-Nation/panthalassa/dapp/module"
	loggerMod "github.com/Bit-Nation/panthalassa/dapp/module/logger"
	uuidv4Mod "github.com/Bit-Nation/panthalassa/dapp/module/uuidv4"
	log "github.com/ipfs/go-log"
	host "github.com/libp2p/go-libp2p-host"
	net "github.com/libp2p/go-libp2p-net"
	pstore "github.com/libp2p/go-libp2p-peerstore"
	ma "github.com/multiformats/go-multiaddr"
	golog "github.com/op/go-logging"
)

var logger = log.Logger("dapp - registry")

// keep track of all running DApps
type Registry struct {
	host           host.Host
	lock           sync.Mutex
	dAppDevStreams map[string]net.Stream
	dAppInstances  map[string]*dapp.DApp
	closeChan      chan *dapp.JsonRepresentation
	client         Client
}

// create new dApp registry
func NewDAppRegistry(h host.Host, client Client) *Registry {

	r := &Registry{
		host:           h,
		lock:           sync.Mutex{},
		dAppDevStreams: map[string]net.Stream{},
		dAppInstances:  map[string]*dapp.DApp{},
		closeChan:      make(chan *dapp.JsonRepresentation),
		client:         client,
	}

	// add worker to remove DApps
	go func() {
		for {
			select {
			case cc := <-r.closeChan:
				r.lock.Lock()
				delete(r.dAppInstances, hex.EncodeToString(cc.SignaturePublicKey))
				// @todo send signal to client that this app was shut down
				r.lock.Unlock()
			}
		}
	}()

	return r

}

// start a DApp
func (r *Registry) StartDApp(dApp *dapp.JsonRepresentation) error {

	l, err := golog.GetLogger("app name")
	if err != nil {
		return err
	}

	vmModules := []module.Module{
		&uuidv4Mod.UUIDV4{},
	}

	// if there is a stream for this DApp
	// we would like to mutate the logger
	// to write to the stream we have for development
	// this will write logs to the stream
	exist, stream := r.getDAppDevStream(dApp.SignaturePublicKey)
	if exist {
		// append log module
		lm, err := loggerMod.New(stream)
		if err != nil {
			return err
		}
		vmModules = append(vmModules, lm)
	} else {
		l.SetBackend(golog.AddModuleLevel(golog.NewLogBackend(ioutil.Discard, "", 0)))
	}

	app, err := dapp.New(l, dApp, vmModules, r.closeChan)
	if err != nil {
		return err
	}

	// add DApp to state
	r.lock.Lock()
	r.dAppInstances[app.ID()] = app
	r.lock.Unlock()

	return nil

}

// use this to connect to a development server
func (r *Registry) ConnectDevelopmentServer(addr ma.Multiaddr) error {

	// address to peer info
	pInfo, err := pstore.InfoFromP2pAddr(addr)
	if err != nil {
		return err
	}

	// connect to peer
	if err := r.host.Connect(context.Background(), *pInfo); err != nil {
		return err
	}

	// create stream to development peer
	str, err := r.host.NewStream(context.Background(), pInfo.ID, "/dapp-development/0.0.0")
	if err != nil {
		return err
	}

	// handle stream
	r.devStreamHandler(str)

	return nil
}

func (r *Registry) ShutDown(dAppJson dapp.JsonRepresentation) error {

	// shut down DApp & remove from state
	r.lock.Lock()
	dApp, exist := r.dAppInstances[hex.EncodeToString(dAppJson.SignaturePublicKey)]
	if !exist {
		return errors.New("DApp is not running")
	}
	dApp.Close()
	r.lock.Unlock()

	return nil

}

package registry

import (
	"encoding/hex"
	"errors"
	"io/ioutil"
	"sync"

	dapp "github.com/Bit-Nation/panthalassa/dapp"
	module "github.com/Bit-Nation/panthalassa/dapp/module"
	loggerMod "github.com/Bit-Nation/panthalassa/dapp/module/logger"
	uuidv4Mod "github.com/Bit-Nation/panthalassa/dapp/module/uuidv4"
	state "github.com/Bit-Nation/panthalassa/state"
	log "github.com/ipfs/go-log"
	host "github.com/libp2p/go-libp2p-host"
	net "github.com/libp2p/go-libp2p-net"
	golog "github.com/op/go-logging"
)

var logger = log.Logger("dapp - registry")

// keep track of all running DApps
type Registry struct {
	host           host.Host
	state          *state.State
	lock           sync.Mutex
	dAppDevStreams map[string]net.Stream
	dAppInstances  map[string]*dapp.DApp
	closeChan      chan *dapp.JsonRepresentation
}

// create new dApp registry
func NewDAppRegistry(h host.Host, state *state.State) *Registry {

	r := &Registry{
		host:           h,
		state:          state,
		lock:           sync.Mutex{},
		dAppDevStreams: map[string]net.Stream{},
		dAppInstances:  map[string]*dapp.DApp{},
		closeChan:      make(chan *dapp.JsonRepresentation),
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

	// set dapp development stream
	// @todo maybe it make sense to register the protocol only
	// @todo if we are in development mode. It will expose less
	// @todo attack vectors.
	h.SetStreamHandler("/dapp-development/0.0.0", r.devStreamHandler)

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
		l.SetBackend(golog.AddModuleLevel(golog.NewLogBackend(stream, "", 0)))
		// append log module
		vmModules = append(vmModules, loggerMod.New(l))
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

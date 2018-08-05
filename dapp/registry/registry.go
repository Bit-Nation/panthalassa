package registry

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"io/ioutil"
	"sync"
	"time"

	api "github.com/Bit-Nation/panthalassa/api"
	dapp "github.com/Bit-Nation/panthalassa/dapp"
	module "github.com/Bit-Nation/panthalassa/dapp/module"
	ethAddrMod "github.com/Bit-Nation/panthalassa/dapp/module/ethAddress"
	ethWSMod "github.com/Bit-Nation/panthalassa/dapp/module/ethWebSocket"
	loggerMod "github.com/Bit-Nation/panthalassa/dapp/module/logger"
	messageModule "github.com/Bit-Nation/panthalassa/dapp/module/message"
	modalMod "github.com/Bit-Nation/panthalassa/dapp/module/modal"
	randBytes "github.com/Bit-Nation/panthalassa/dapp/module/randBytes"
	renderDApp "github.com/Bit-Nation/panthalassa/dapp/module/renderer/dapp"
	renderMsg "github.com/Bit-Nation/panthalassa/dapp/module/renderer/message"
	sendEthTxMod "github.com/Bit-Nation/panthalassa/dapp/module/sendEthTx"
	uuidv4Mod "github.com/Bit-Nation/panthalassa/dapp/module/uuidv4"
	db "github.com/Bit-Nation/panthalassa/db"
	ethws "github.com/Bit-Nation/panthalassa/ethws"
	keyManager "github.com/Bit-Nation/panthalassa/keyManager"
	log "github.com/ipfs/go-log"
	host "github.com/libp2p/go-libp2p-host"
	net "github.com/libp2p/go-libp2p-net"
	pstore "github.com/libp2p/go-libp2p-peerstore"
	ma "github.com/multiformats/go-multiaddr"
	golog "github.com/op/go-logging"
	ed25519 "golang.org/x/crypto/ed25519"
)

var logger = log.Logger("dapp - registry")

// keep track of all running DApps
type Registry struct {
	host           host.Host
	lock           sync.Mutex
	dAppDevStreams map[string]net.Stream
	dAppInstances  map[string]*dapp.DApp
	closeChan      chan *dapp.Data
	conf           Config
	ethWS          *ethws.EthereumWS
	api            *api.API
	km             *keyManager.KeyManager
	dAppDB         dapp.Storage
	msgDB          db.ChatMessageStorage
}

type Config struct {
	EthWSEndpoint string
}

// create new dApp registry
func NewDAppRegistry(h host.Host, conf Config, api *api.API, km *keyManager.KeyManager, dAppDB dapp.Storage, msgDB db.ChatMessageStorage) *Registry {

	r := &Registry{
		host:           h,
		lock:           sync.Mutex{},
		dAppDevStreams: map[string]net.Stream{},
		dAppInstances:  map[string]*dapp.DApp{},
		closeChan:      make(chan *dapp.Data),
		conf:           conf,
		ethWS: ethws.New(ethws.Config{
			Retry: time.Second,
			WSUrl: conf.EthWSEndpoint,
		}),
		api:    api,
		km:     km,
		dAppDB: dAppDB,
		msgDB:  msgDB,
	}

	// add worker to remove DApps
	go func() {
		for {
			select {
			case cc := <-r.closeChan:
				r.lock.Lock()
				delete(r.dAppInstances, hex.EncodeToString(cc.UsedSigningKey))
				// @todo send signal to client that this app was shut down
				r.lock.Unlock()
			}
		}
	}()

	return r

}

// start a DApp
func (r *Registry) StartDApp(dAppSigningKey ed25519.PublicKey, timeOut time.Duration) error {

	// fetch DApp
	dApp, err := r.dAppDB.Get(dAppSigningKey)
	if err != nil {
		return err
	}
	if dApp == nil {
		return fmt.Errorf("failed to fetch DApp for signing key: %x", dAppSigningKey)
	}

	// get logger
	var l *golog.Logger
	if l, err = golog.GetLogger("app name"); err != nil {
		return err
	}

	vmModules := []module.Module{
		uuidv4Mod.New(l),
		ethWSMod.New(l, r.ethWS),
		modalMod.New(l, r.api, dApp.UsedSigningKey),
		sendEthTxMod.New(r.api, l),
		randBytes.New(l),
		ethAddrMod.New(r.km),
		renderMsg.New(l),
		renderDApp.New(l),
		messageModule.New(r.msgDB, dAppSigningKey, l),
	}

	// if there is a stream for this DApp
	// we would like to mutate the logger
	// to write to the stream we have for development
	// this will write logs to the stream
	exist, stream := r.getDAppDevStream(dApp.UsedSigningKey)
	if exist {
		// append log module
		logger, err := loggerMod.New(stream)
		l = logger.Logger
		if err != nil {
			return err
		}
		vmModules = append(vmModules, logger)
	} else {
		if err != nil {
			return err
		}
		l.SetBackend(golog.AddModuleLevel(golog.NewLogBackend(ioutil.Discard, "", 0)))
	}

	app, err := dapp.New(l, dApp, vmModules, r.closeChan, timeOut)
	if err != nil {
		return err
	}

	// add DApp to state
	r.lock.Lock()
	r.dAppInstances[app.ID()] = app
	r.lock.Unlock()

	return nil

}

// open DApp
func (r *Registry) OpenDApp(id, context string) error {
	r.lock.Lock()
	defer r.lock.Unlock()
	if _, exist := r.dAppInstances[id]; !exist {
		return errors.New("it seems like that this app hasn't been started yet")
	}
	return r.dAppInstances[id].RenderDApp(context)
}

func (r *Registry) RenderMessage(id, msg, context string) (string, error) {
	r.lock.Lock()
	defer r.lock.Unlock()
	if _, exist := r.dAppInstances[id]; !exist {
		return "", errors.New("it seems like that this app hasn't been started yet")
	}
	return r.dAppInstances[id].RenderMessage(context)
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

// call a function in a DApp
// @todo when the function call takes too long the lock could stay locked. We should "free" it asap
func (r *Registry) CallFunction(dAppID string, funcId uint, args string) error {

	r.lock.Lock()
	defer r.lock.Unlock()

	dApp, exist := r.dAppInstances[dAppID]
	if !exist {
		return errors.New("can't call function of dApp that hasn't started")
	}

	return dApp.CallFunction(funcId, args)

}

func (r *Registry) ShutDown(signingKey ed25519.PublicKey) error {

	// shut down DApp & remove from state
	r.lock.Lock()
	dApp, exist := r.dAppInstances[hex.EncodeToString(signingKey)]
	if !exist {
		return errors.New("DApp is not running")
	}
	dApp.Close()
	r.lock.Unlock()

	return nil

}

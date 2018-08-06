package dapp

import (
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	module "github.com/Bit-Nation/panthalassa/dapp/module"
	cbModule "github.com/Bit-Nation/panthalassa/dapp/module/callbacks"
	dbModule "github.com/Bit-Nation/panthalassa/dapp/module/db"
	dAppRenderer "github.com/Bit-Nation/panthalassa/dapp/module/renderer/dapp"
	msgRenderer "github.com/Bit-Nation/panthalassa/dapp/module/renderer/message"
	bolt "github.com/coreos/bbolt"
	logger "github.com/op/go-logging"
	otto "github.com/robertkrimen/otto"
)

type DApp struct {
	vm           *otto.Otto
	logger       *logger.Logger
	app          *Data
	closeChan    chan<- *Data
	dAppRenderer *dAppRenderer.Module
	msgRenderer  *msgRenderer.Module
	cbMod        *cbModule.Module
	dbMod        *dbModule.BoltStorage
}

// close DApp
func (d *DApp) Close() {
	d.vm.Interrupt <- func() {
		d.logger.Info(fmt.Sprintf("shutting down: %s (%s)", hex.EncodeToString(d.app.UsedSigningKey), d.app.Name))
		d.closeChan <- d.app
	}
}

func (d *DApp) ID() string {
	return hex.EncodeToString(d.app.UsedSigningKey)
}

func (d *DApp) OpenDApp(context string) error {
	return d.dAppRenderer.OpenDApp(context)
}

func (d *DApp) RenderMessage(payload string) (string, error) {
	return d.msgRenderer.RenderMessage(payload)
}

func (d *DApp) CallFunction(id uint, args string) error {
	return d.cbMod.CallFunction(id, args)
}

// will start a DApp based on the given config file
func New(l *logger.Logger, app *Data, vmModules []module.Module, closer chan<- *Data, timeOut time.Duration, db *bolt.DB) (*DApp, error) {

	// check if app is valid
	valid, err := app.VerifySignature()
	if err != nil {
		return nil, err
	}
	if !valid {
		return nil, InvalidSignature
	}

	// create VM
	vm := otto.New()
	vm.Interrupt = make(chan func(), 1)

	// register all vm modules
	for _, m := range vmModules {
		if err := m.Register(vm); err != nil {
			return nil, err
		}
	}

	// register DApp renderer
	dr := dAppRenderer.New(l)
	if err := dr.Register(vm); err != nil {
		return nil, err
	}

	// register message renderer
	mr := msgRenderer.New(l)
	if err := dr.Register(vm); err != nil {
		return nil, err
	}

	// register callbacks module
	cbm := cbModule.New(l)
	if err := cbm.Register(vm); err != nil {
		return nil, err
	}

	// register database (the one used by the dapp)
	dAppDBStorage, err := dbModule.NewBoltStorage(db, app.UsedSigningKey)
	if err != nil {
		return nil, err
	}
	dbm := dbModule.New(dAppDBStorage)
	if err := dbm.Register(vm); err != nil {
		return nil, err
	}

	dApp := &DApp{
		vm:           vm,
		logger:       l,
		app:          app,
		closeChan:    closer,
		dAppRenderer: dr,
		msgRenderer:  mr,
		cbMod:        cbm,
		dbMod:        dAppDBStorage,
	}

	wait := make(chan error, 1)

	// start the DApp async
	go func() {
		_, err := vm.Run(app.Code)
		if err != nil {
			l.Errorf(err.Error())
		}
		wait <- err
	}()

	// wait for the DApp with given timeout
	select {
	case err := <-wait:
		if err != nil {
			return nil, err
		}
		return dApp, nil
	case <-time.After(timeOut):
		vm.Interrupt <- func() {}
		closer <- app
		return nil, errors.New("timeout - failed to start DApp")
	}

}

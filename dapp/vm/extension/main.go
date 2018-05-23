package extension

import (
	log "github.com/ipfs/go-log"
	otto "github.com/robertkrimen/otto"
)

type OttoFunction = func(otto otto.FunctionCall) otto.Value

var logger = log.Logger("vm_extensions")

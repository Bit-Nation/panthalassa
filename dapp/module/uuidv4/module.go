package uuidv4

import (
	log "github.com/ipfs/go-log"

	validator "github.com/Bit-Nation/panthalassa/dapp/validator"
	logger "github.com/op/go-logging"
	uuid "github.com/satori/go.uuid"
	duktape "gopkg.in/olebedev/go-duktape.v3"
)

var newUuid = uuid.NewV4

var sysLogger = log.Logger("uuidv4")

// register an module on the given
// vm that allows to generate uuid's
// of version 4. it expects a
// callback as it's only argument
type UUIDV4 struct {
	logger *logger.Logger
}

func New(l *logger.Logger) *UUIDV4 {
	return &UUIDV4{logger: l}
}

func (r *UUIDV4) Close() error {
	return nil
}

func (r *UUIDV4) Register(vm *duktape.Context) error {

	_, err := vm.PushGlobalGoFunction("uuidV4", func(context *duktape.Context) int {

		sysLogger.Debug("generate uuidv4")

		// validate function call
		v := validator.New()
		v.Set(0, &validator.TypeFunction)
		if err := v.Validate(context); err != nil {
			//vm.Run(`throw new Error("uuidV4 expects an callback as it's parameter")`)
			return 1
		}

		// create uuid
		id, err := newUuid()
		// call callback with error
		if err != nil {
			context.PushString(err.Error())
			context.PushUndefined()
			context.Call(2)
			return 1
		}

		// call callback with uuid
		context.PushUndefined()
		context.PushString(id.String())
		context.Call(2)

		return 0

	})
	return err
}

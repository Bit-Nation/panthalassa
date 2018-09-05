package randBytes

import (
	"crypto/rand"
	"fmt"
	"strings"

	validator "github.com/Bit-Nation/panthalassa/dapp/validator"
	log "github.com/ipfs/go-log"
	logger "github.com/op/go-logging"
	duktape "gopkg.in/olebedev/go-duktape.v3"
)

var randSource = rand.Reader

var sysLog = log.Logger("rand bytes")

func New(l *logger.Logger) *Module {
	return &Module{
		logger: l,
	}
}

type Module struct {
	logger *logger.Logger
}

func (m *Module) Close() error {
	return nil
}

func (m *Module) Register(vm *duktape.Context) error {
	_, err := vm.PushGlobalGoFunction("randomBytes", func(context *duktape.Context) int {
		sysLog.Debug("generate random bytes")
		var itemsToPopBeforeCallback int
		// validate call
		v := validator.New()
		v.Set(0, &validator.TypeNumber)
		v.Set(1, &validator.TypeFunction)
		handleError := func(errMsg string) int {
			if context.IsFunction(1) {
				context.PopN(itemsToPopBeforeCallback)
				context.PushString(errMsg)
				context.Call(1)
				return 0
			}
			m.logger.Error(errMsg)
			return 1
		}
		if err := v.Validate(context); err != nil {
			m.logger.Error(err.Error())
		}

		// convert to integer
		amount := context.ToInt(0)
		destination := make([]byte, amount)
		_, err := randSource.Read(destination)
		if err != nil {
			handleError(err.Error())
			return 1
		} // if err != nil {

		// call callback
		context.PopN(itemsToPopBeforeCallback)
		context.PushUndefined()
		byteString := fmt.Sprint(destination)
		jsFriendlyString := strings.Replace(strings.TrimSuffix(strings.TrimPrefix(byteString, "["), "]"), " ", ", ", -1)
		context.PushString(jsFriendlyString)
		context.Call(2)
		return 0
	})

	return err
}

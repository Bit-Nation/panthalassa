package logger

import (
	"strings"

	logger "github.com/op/go-logging"
	otto "github.com/robertkrimen/otto"
)

type Logger struct {
	logger *logger.Logger
}

func New(l *logger.Logger) *Logger {
	return &Logger{
		logger: l,
	}
}

func (l *Logger) Name() string {
	return "LOGGER"
}

// Register a module that writes console.log
// to the given logger
func (l *Logger) Register(vm *otto.Otto) {

	vm.Set("console", map[string]interface{}{
		"log": func(call otto.FunctionCall) otto.Value {
			toLog := []string{}
			for _, arg := range call.ArgumentList {
				toLog = append(toLog, arg.String())
			}
			l.logger.Info(strings.Join(toLog, ","))
			return otto.Value{}
		},
	})

}

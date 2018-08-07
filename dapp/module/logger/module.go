package logger

import (
	"bufio"
	"encoding/base64"
	"strings"

	log "github.com/ipfs/go-log"
	net "github.com/libp2p/go-libp2p-net"
	logger "github.com/op/go-logging"
	otto "github.com/robertkrimen/otto"
)

var sysLog = log.Logger("logger module logger")

type streamLogger struct {
	writer *bufio.Writer
}

func (l *streamLogger) Write(data []byte) (int, error) {

	// write log encoded as base64 to make sure there are no spaces in the dataset
	if a, err := l.writer.Write([]byte(base64.StdEncoding.EncodeToString(data))); err != nil {
		return a, err
	}

	// write line break as delimiter
	if err := l.writer.WriteByte(0x0A); err != nil {
		return 0, err
	}

	return len(data), nil
}

type Logger struct {
	Logger *logger.Logger
}

func New(stream net.Stream) (*Logger, error) {

	l, err := logger.GetLogger("")
	if err != nil {
		return nil, err
	}

	w := bufio.NewWriter(stream)
	loggerStream := &streamLogger{writer: w}

	l.SetBackend(logger.AddModuleLevel(logger.NewLogBackend(loggerStream, "", 0)))

	return &Logger{
		Logger: l,
	}, nil
}

// Register a module that writes console.log
// to the given logger
func (l *Logger) Register(vm *otto.Otto) error {

	return vm.Set("console", map[string]interface{}{
		"log": func(call otto.FunctionCall) otto.Value {

			sysLog.Debug("write log statement")

			toLog := []string{}
			for _, arg := range call.ArgumentList {
				sysLog.Debug("log: ", toLog)
				toLog = append(toLog, arg.String())
			}
			l.Logger.Info(strings.Join(toLog, ","))
			return otto.Value{}
		},
	})

}

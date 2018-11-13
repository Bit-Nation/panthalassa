package logger

import (
	"encoding/base64"
	"strings"

	log "github.com/ipfs/go-log"
	net "github.com/libp2p/go-libp2p-net"
	logger "github.com/op/go-logging"
	otto "github.com/robertkrimen/otto"
)

var sysLog = log.Logger("logger module logger")

type streamLogger struct {
	writer net.Stream
}

func (l *streamLogger) Write(data []byte) (int, error) {

	// write log into stream
	i, err := l.writer.Write([]byte(base64.StdEncoding.EncodeToString(data)))
	if err != nil {
		return i, err
	}

	// write delimiter
	_, err = l.writer.Write([]byte{'\n'})
	if err != nil {
		return i, err
	}

	// +1 because the delimiter is just one byte
	return i + +1, nil
}

type Logger struct {
	Logger *logger.Logger
}

func New(stream net.Stream) (*Logger, error) {

	l, err := logger.GetLogger("")
	if err != nil {
		return nil, err
	}

	loggerStream := &streamLogger{writer: stream}

	l.SetBackend(logger.AddModuleLevel(logger.NewLogBackend(loggerStream, "", 0)))

	return &Logger{
		Logger: l,
	}, nil
}

func (l *Logger) Close() error {
	return nil
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

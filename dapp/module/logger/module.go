package logger

import (
	"encoding/base64"
	"strings"

	log "github.com/ipfs/go-log"
	net "github.com/libp2p/go-libp2p-net"
	logger "github.com/op/go-logging"
	duktape "gopkg.in/olebedev/go-duktape.v3"
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
func (l *Logger) Register(vm *duktape.Context) error {
	//@TODO Find a way to overwrite method console.log if neccessary, so that we don't need to call consoleLog
	_, err := vm.PushGlobalGoFunction("consoleLog", func(context *duktape.Context) int {
		sysLog.Debug("write log statement")
		toLog := []string{}
		var i int
		for {
			if context.GetType(i).IsNone() {
				break
			}
			sysLog.Debug("log: ", toLog)
			toLog = append(toLog, context.ToString(i))
			i++
		}
		l.Logger.Info(strings.Join(toLog, ","))
		return 0
	})
	return err
}

package logger

import (
	"bufio"
	"strings"

	pb "github.com/Bit-Nation/panthalassa/dapp/registry/pb"
	net "github.com/libp2p/go-libp2p-net"
	mc "github.com/multiformats/go-multicodec"
	protoMc "github.com/multiformats/go-multicodec/protobuf"
	logger "github.com/op/go-logging"
	otto "github.com/robertkrimen/otto"
)

type streamLogger struct {
	writer  *bufio.Writer
	encoder mc.Encoder
}

func (l *streamLogger) Write(data []byte) (int, error) {

	msg := pb.Message{
		Type: pb.Message_LOG,
		Log:  data,
	}

	if err := l.encoder.Encode(&msg); err != nil {
		return 0, err
	}

	if err := l.writer.Flush(); err != nil {
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
	loggerStream := &streamLogger{
		writer:  w,
		encoder: protoMc.Multicodec(nil).Encoder(w),
	}

	l.SetBackend(logger.AddModuleLevel(logger.NewLogBackend(loggerStream, "", 0)))

	return &Logger{
		Logger: l,
	}, nil
}

func (l *Logger) Name() string {
	return "LOGGER"
}

// Register a module that writes console.log
// to the given logger
func (l *Logger) Register(vm *otto.Otto) error {

	return vm.Set("console", map[string]interface{}{
		"log": func(call otto.FunctionCall) otto.Value {
			toLog := []string{}
			for _, arg := range call.ArgumentList {
				toLog = append(toLog, arg.String())
			}
			l.Logger.Info(strings.Join(toLog, ","))
			return otto.Value{}
		},
	})

}

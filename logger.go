package panthalassa

import (
	"fmt"
	"time"
)

type Logger interface {
	Info(msg string)
	Debug(msg string)
	Error(msg string)
	Warn(msg string)
}

//Logger that print's to console. Create a logger by NewCliLogger
type CliLogger struct{}

func (l *CliLogger) Info(msg string) {
	fmt.Printf("%s - INFO - %s", time.Now().UTC().String(), msg)
}

func (l *CliLogger) Debug(msg string) {
	fmt.Printf("%s - DEBUG - %s", time.Now().UTC().String(), msg)
}

func (l *CliLogger) Error(msg string) {
	fmt.Printf("%s - ERROR - %s", time.Now().UTC().String(), msg)
}

func (l *CliLogger) Warm(msg string) {
	fmt.Printf("%s - Warm - %s", time.Now().UTC().String(), msg)
}

type LoggerCallback = func(msg string)

func NewCliLogger() CliLogger {
	return CliLogger{}
}

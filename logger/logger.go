package logger

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
	fmt.Printf("%s - INFO - %s \n", time.Now().UTC().String(), msg)
}

func (l *CliLogger) Debug(msg string) {
	fmt.Printf("%s - DEBUG - %s \n", time.Now().UTC().String(), msg)
}

func (l *CliLogger) Error(msg string) {
	fmt.Printf("%s - ERROR - %s \n", time.Now().UTC().String(), msg)
}

func (l *CliLogger) Warn(msg string) {
	fmt.Printf("%s - Warm - %s \n", time.Now().UTC().String(), msg)
}

func NewCliLogger() CliLogger {
	return CliLogger{}
}

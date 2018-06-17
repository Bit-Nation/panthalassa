package logger

import (
	"testing"

	logger "github.com/op/go-logging"
	otto "github.com/robertkrimen/otto"
	require "github.com/stretchr/testify/require"
)

type testValue struct {
	js        string
	assertion func(consoleOut string)
}

type testWriter struct {
	assertion func(consoleOut string)
}

func (w testWriter) Write(b []byte) (int, error) {
	w.assertion(string(b))
	return 0, nil
}

func TestLoggerModule(t *testing.T) {

	testValues := []testValue{
		testValue{
			js: `console.log(1, 2, 3)`,
			assertion: func(consoleOut string) {
				require.Equal(t, "1,2,3\n", consoleOut)
			},
		},
		testValue{
			js: `console.log("hi","there")`,
			assertion: func(consoleOut string) {
				require.Equal(t, "hi,there\n", consoleOut)
			},
		},
		testValue{
			js: `console.log({key: 4})`,
			assertion: func(consoleOut string) {
				require.Equal(t, "[object Object]\n", consoleOut)
			},
		},
		testValue{
			js: `
				var cb = function(){};
				console.log(cb)
			`,
			assertion: func(consoleOut string) {
				require.Equal(t, "function(){}\n", consoleOut)
			},
		},
		testValue{
			js: `console.log("hi",1)`,
			assertion: func(consoleOut string) {
				require.Equal(t, "hi,1\n", consoleOut)
			},
		},
	}

	for _, testValue := range testValues {

		// create VM
		vm := otto.New()

		// create logger
		b := logger.NewLogBackend(testWriter{
			assertion: testValue.assertion,
		}, "", 0)
		l, err := logger.GetLogger("-")
		require.Nil(t, err)
		l.SetBackend(logger.AddModuleLevel(b))

		loggerModule, err := New(nil)
		require.Nil(t, err)
		loggerModule.logger = l
		loggerModule.Register(vm)

		_, err = vm.Run(testValue.js)
		require.Nil(t, err)

	}

}

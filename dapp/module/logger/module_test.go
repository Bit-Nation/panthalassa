package logger

import (
	"testing"

	logger "github.com/op/go-logging"
	require "github.com/stretchr/testify/require"
	duktape "gopkg.in/olebedev/go-duktape.v3"
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

	//@TODO Find a way to overwrite method console.log if neccessary, so that we don't need to call consoleLog
	testValues := []testValue{
		testValue{
			js: `consoleLog(1, 2, 3)`,
			assertion: func(consoleOut string) {
				require.Equal(t, "1,2,3\n", consoleOut)
			},
		},
		testValue{
			js: `consoleLog("hi","there")`,
			assertion: func(consoleOut string) {
				require.Equal(t, "hi,there\n", consoleOut)
			},
		},
		testValue{
			js: `consoleLog({key: 4})`,
			assertion: func(consoleOut string) {
				require.Equal(t, "[object Object]\n", consoleOut)
			},
		},
		testValue{
			js: `
				var cb = function(){};
				consoleLog(cb)
			`,
			assertion: func(consoleOut string) {
				require.Equal(t, "function () { [ecmascript code] }\n", consoleOut)
			},
		},
		testValue{
			js: `consoleLog("hi",1)`,
			assertion: func(consoleOut string) {
				require.Equal(t, "hi,1\n", consoleOut)
			},
		},
	}

	for _, testValue := range testValues {

		// create VM
		vm := duktape.New()

		// create logger
		b := logger.NewLogBackend(testWriter{
			assertion: testValue.assertion,
		}, "", 0)
		l, err := logger.GetLogger("-")
		require.Nil(t, err)
		l.SetBackend(logger.AddModuleLevel(b))

		loggerModule, err := New(nil)
		require.Nil(t, err)
		loggerModule.Logger = l
		loggerModule.Register(vm)

		err = vm.PevalString((testValue.js))
		require.Nil(t, err)

	}

}

package glog

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestZapLogger(t *testing.T) {
	buf := new(bytes.Buffer)

	logger := InitZapLogger(&LoggerConfig{
		ServiceName: "my app",
		Level:       DebugLevel,
		LogDir:      "./log",
		InConsole:   true,
	})
	defer logger.Sync()

	hlog.SetLogger(logger)
	hlog.SetOutput(buf)
	hlog.SetLevel(hlog.LevelDebug)

	type logMap map[string]string

	logTestSlice := []logMap{
		{
			"logMessage":       "this is a trace log",
			"formatLogMessage": "this is a trace log: %s",
			"logLevel":         "Trace",
			"zapLogLevel":      "debug",
		},
		{
			"logMessage":       "this is a debug log",
			"formatLogMessage": "this is a debug log: %s",
			"logLevel":         "Debug",
			"zapLogLevel":      "debug",
		},
		{
			"logMessage":       "this is a info log",
			"formatLogMessage": "this is a info log: %s",
			"logLevel":         "Info",
			"zapLogLevel":      "info",
		},
		{
			"logMessage":       "this is a notice log",
			"formatLogMessage": "this is a notice log: %s",
			"logLevel":         "Notice",
			"zapLogLevel":      "warn",
		},
		{
			"logMessage":       "this is a warn log",
			"formatLogMessage": "this is a warn log: %s",
			"logLevel":         "Warn",
			"zapLogLevel":      "warn",
		},
		{
			"logMessage":       "this is a error log",
			"formatLogMessage": "this is a error log: %s",
			"logLevel":         "Error",
			"zapLogLevel":      "error",
		},
		{
			"logMessage":       "this is a fatal log",
			"formatLogMessage": "this is a fatal log: %s",
			"logLevel":         "Fatal",
			"zapLogLevel":      "fatal",
		},
	}

	testHertzLogger := reflect.ValueOf(logger)

	for _, v := range logTestSlice {
		t.Run(v["logLevel"], func(t *testing.T) {
			if v["logLevel"] == "Fatal" {
				defer func() {
					assert.Equal(t, "this is a fatal log", recover())
				}()
			}
			logFunc := testHertzLogger.MethodByName(v["logLevel"])
			logFunc.Call([]reflect.Value{
				reflect.ValueOf(v["logMessage"]),
			})
			assert.Contains(t, buf.String(), v["logMessage"])
			assert.Contains(t, buf.String(), v["zapLogLevel"])

			buf.Reset()

			logfFunc := testHertzLogger.MethodByName(fmt.Sprintf("%sf", v["logLevel"]))
			logfFunc.Call([]reflect.Value{
				reflect.ValueOf(v["formatLogMessage"]),
				reflect.ValueOf(v["logLevel"]),
			})
			assert.Contains(t, buf.String(), fmt.Sprintf(v["formatLogMessage"], v["logLevel"]))
			assert.Contains(t, buf.String(), v["zapLogLevel"])

			buf.Reset()

			ctx := context.Background()
			ctxLogfFunc := testHertzLogger.MethodByName(fmt.Sprintf("Ctx%sf", v["logLevel"]))
			ctxLogfFunc.Call([]reflect.Value{
				reflect.ValueOf(ctx),
				reflect.ValueOf(v["formatLogMessage"]),
				reflect.ValueOf(v["logLevel"]),
			})
			assert.Contains(t, buf.String(), fmt.Sprintf(v["formatLogMessage"], v["logLevel"]))
			assert.Contains(t, buf.String(), v["zapLogLevel"])

			buf.Reset()
		})
	}
}

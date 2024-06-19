package glog

import (
	"go.uber.org/zap/zapcore"
)

type LoggerType uint8

const (
	LoggerTypeZap LoggerType = iota + 1
)

const (
	logOutputTypeStdout      = "stdout"     // 输出到控制台
	logOutputTypeDefaultFile = "file"       // 输出到普通文件
	logOutputTypeWarnFatal   = "warn_fatal" // 输出到警告和致命错误日志文件

	logOutputFileDefaultSuffix   = ".log"
	logOutputFileWarnFatalSuffix = ".wf.log"
)

var logOutputFileSuffixMap = map[string]string{
	logOutputTypeDefaultFile: logOutputFileDefaultSuffix,
	logOutputTypeWarnFatal:   logOutputFileWarnFatalSuffix,
}

const (
	KeyTraceId    = "trace-id"
	KeyTraceFlags = "trace-flags"
	KeySpanId     = "span-id"
	KeyIp         = "ip"
	KeyUri        = "uri"
)

type Level string

const (
	DebugLevel Level = "debug"
	InfoLevel  Level = "info"
	WarnLevel  Level = "warn"
	ErrorLevel Level = "error"
	PanicLevel Level = "panic"
	FatalLevel Level = "fatal"
)

var logLevelMap = map[Level]zapcore.Level{
	DebugLevel: zapcore.DebugLevel,
	InfoLevel:  zapcore.InfoLevel,
	WarnLevel:  zapcore.WarnLevel,
	ErrorLevel: zapcore.ErrorLevel,
	PanicLevel: zapcore.PanicLevel,
	FatalLevel: zapcore.FatalLevel,
}

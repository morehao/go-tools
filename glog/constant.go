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

	logOutputFieldDefaultSuffix  = ".log"
	logOutputFileWarnFatalSuffix = ".wf.log"
)

var logOutputFileSuffixMap = map[string]string{
	logOutputTypeDefaultFile: logOutputFieldDefaultSuffix,
	logOutputTypeWarnFatal:   logOutputFileWarnFatalSuffix,
}

const (
	KeyRequestId  = "requestId"
	KeyTraceId    = "traceId"
	KeyTraceFlags = "traceFlags"
	KeySpanId     = "spanId"

	MsgFlagNotice = "notice"

	KeySkipLog          = "skip"
	KeyService          = "service"
	KeyHost             = "host"
	KeyClientIp         = "clientIp"
	KeyHandle           = "handle"
	KeyProto            = "proto"
	KeyRefer            = "refer"
	KeyUserAgent        = "userAgent"
	KeyHeader           = "header"
	KeyCookie           = "cookie"
	KeyUri              = "uri"
	KeyMethod           = "method"
	KeyHttpStatusCode   = "httpStatusCode"
	KeyRequestQuery     = "requestQuery"
	KeyRequestBody      = "requestBody"
	KeyRequestBodySize  = "requestBodySize"
	KeyResponseCode     = "responseCode"
	KeyResponseBody     = "responseBody"
	KeyResponseBodySize = "responseBodySize"
	KeyRequestStartTime = "requestStartTime"
	KeyRequestEndTime   = "requestEndTime"
	KeyCost             = "cost"
	KeyRequestErr       = "requestErr"
	KeyErrorCode        = "errorCode"
	KeyErrorMsg         = "errorMsg"
	KeyAffectedRows     = "affectedRows"
	KeyAddr             = "addr"
	KeyDatabase         = "database"
	KeySql              = "sql"
	KeyCmd              = "cmd"
	KeyCmdContent       = "cmdContent"
	KeyRalCode          = "ralCode"
	KeyFile             = "file"

	ValueProtoHttp  = "http"
	ValueProtoMysql = "mysql"
	ValueProtoRedis = "redis"
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

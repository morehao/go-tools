package glog

import (
	"context"
	"fmt"
	"io"
	"os"
	"path"
	"strings"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/morehao/go-tools/gutils"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger 是一个封装 zap.Logger 的结构体
type zapLogger struct {
	logger *zap.Logger
	cfg    *LoggerConfig
}

func newZapLogger(cfg *LoggerConfig) (*zap.Logger, error) {
	var zapCores []zapcore.Core
	logLevel, ok := logLevelMap[cfg.Level]
	if !ok {
		logLevel = zapcore.InfoLevel
	}
	var infoLevel = zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= logLevel && lvl <= zapcore.InfoLevel
	})

	var errorLevel = zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= logLevel && lvl >= zapcore.WarnLevel
	})

	var stdLevel = zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= logLevel && lvl >= zapcore.DebugLevel
	})
	if cfg.InConsole {
		c := zapcore.NewCore(
			getZapEncoder(),
			getZapLogWriter(cfg, logOutputTypeStdout),
			stdLevel)
		zapCores = append(zapCores, c)
	}

	zapCores = append(zapCores,
		zapcore.NewCore(
			getZapEncoder(),
			getZapLogWriter(cfg, logOutputTypeDefaultFile),
			infoLevel))

	zapCores = append(zapCores,
		zapcore.NewCore(
			getZapEncoder(),
			getZapLogWriter(cfg, logOutputTypeDefaultFile),
			errorLevel))

	zapCores = append(zapCores,
		zapcore.NewCore(
			getZapEncoder(),
			getZapLogWriter(cfg, logOutputTypeWarnFatal),
			errorLevel))

	core := zapcore.NewTee(zapCores...)

	// 开启开发模式，堆栈跟踪
	caller := zap.WithCaller(true)

	development := zap.Development()

	// 设置初始化字段
	filed := zap.Fields()

	// 构造logger
	logger := zap.New(core, filed, caller, development)

	return logger, nil
}

func (l *zapLogger) Debug(ctx context.Context, args ...interface{}) {
	l.ctxLog(DebugLevel, ctx, args...)
}

func (l *zapLogger) Debugf(ctx context.Context, format string, args ...interface{}) {
	l.ctxLogf(DebugLevel, ctx, format, args...)
}

func (l *zapLogger) Debugw(ctx context.Context, msg string, keysAndValues ...interface{}) {
	l.ctxLogw(DebugLevel, ctx, msg, keysAndValues...)
}

func (l *zapLogger) Info(ctx context.Context, args ...interface{}) {
	l.ctxLog(InfoLevel, ctx, args...)
}

func (l *zapLogger) Infof(ctx context.Context, format string, args ...interface{}) {
	l.ctxLogf(InfoLevel, ctx, format, args...)
}

func (l *zapLogger) Infow(ctx context.Context, msg string, keysAndValues ...interface{}) {
	l.ctxLogw(InfoLevel, ctx, msg, keysAndValues...)
}

func (l *zapLogger) Warn(ctx context.Context, args ...interface{}) {
	l.ctxLog(WarnLevel, ctx, args...)
}

func (l *zapLogger) Warnf(ctx context.Context, format string, args ...interface{}) {
	l.ctxLogf(WarnLevel, ctx, format, args...)
}

func (l *zapLogger) Warnw(ctx context.Context, msg string, keysAndValues ...interface{}) {
	l.ctxLogw(WarnLevel, ctx, msg, keysAndValues...)
}

func (l *zapLogger) Error(ctx context.Context, args ...interface{}) {
	l.ctxLog(ErrorLevel, ctx, args...)
}

func (l *zapLogger) Errorf(ctx context.Context, format string, args ...interface{}) {
	l.ctxLogf(ErrorLevel, ctx, format, args...)
}

func (l *zapLogger) Errorw(ctx context.Context, msg string, keysAndValues ...interface{}) {
	l.ctxLogw(ErrorLevel, ctx, msg, keysAndValues...)
}

func (l *zapLogger) Panic(ctx context.Context, args ...interface{}) {
	l.ctxLog(PanicLevel, ctx, args...)
}

func (l *zapLogger) Panicf(ctx context.Context, format string, args ...interface{}) {
	l.ctxLogf(PanicLevel, ctx, format, args...)
}

func (l *zapLogger) Panicw(ctx context.Context, msg string, keysAndValues ...interface{}) {
	l.ctxLogw(PanicLevel, ctx, msg, keysAndValues...)
}

func (l *zapLogger) Fatal(ctx context.Context, args ...interface{}) {
	l.ctxLog(FatalLevel, ctx, args...)
}

func (l *zapLogger) Fatalf(ctx context.Context, format string, args ...interface{}) {
	l.ctxLogf(PanicLevel, ctx, format, args...)
}

func (l *zapLogger) Fatalw(ctx context.Context, msg string, keysAndValues ...interface{}) {
	l.ctxLogw(FatalLevel, ctx, msg, keysAndValues...)
}

func (l *zapLogger) Sync() {
	_ = l.logger.Sync()
}

func (l *zapLogger) ctxLog(level Level, ctx context.Context, args ...interface{}) {
	if nilCtx(ctx) {
		return
	}
	switch level {
	case DebugLevel:
		l.logger.Sugar().With(l.extraFields(ctx)...).Debug(args...)
	case InfoLevel:
		l.logger.Sugar().With(l.extraFields(ctx)...).Info(args...)
	case WarnLevel:
		l.logger.Sugar().With(l.extraFields(ctx)...).Warn(args...)
	case ErrorLevel:
		l.logger.Sugar().With(l.extraFields(ctx)...).Error(args...)
	case PanicLevel:
		l.logger.Sugar().With(l.extraFields(ctx)...).Panic(args...)
	case FatalLevel:
		l.logger.Sugar().With(l.extraFields(ctx)...).Fatal(args...)
	}
}

func (l *zapLogger) ctxLogf(level Level, ctx context.Context, format string, args ...interface{}) {
	if nilCtx(ctx) {
		return
	}
	switch level {
	case DebugLevel:
		l.logger.Sugar().With(l.extraFields(ctx)...).Debugf(format, args...)
	case InfoLevel:
		l.logger.Sugar().With(l.extraFields(ctx)...).Infof(format, args...)
	case WarnLevel:
		l.logger.Sugar().With(l.extraFields(ctx)...).Warnf(format, args...)
	case ErrorLevel:
		l.logger.Sugar().With(l.extraFields(ctx)...).Errorf(format, args...)
	case PanicLevel:
		l.logger.Sugar().With(l.extraFields(ctx)...).Panicf(format, args...)
	case FatalLevel:
		l.logger.Sugar().With(l.extraFields(ctx)...).Fatalf(format, args...)
	}
}

func (l *zapLogger) ctxLogw(level Level, ctx context.Context, msg string, keysAndValues ...interface{}) {
	if nilCtx(ctx) {
		return
	}
	switch level {
	case DebugLevel:
		l.logger.Sugar().With(l.extraFields(ctx)...).Debugw(msg, keysAndValues...)
	case InfoLevel:
		l.logger.Sugar().With(l.extraFields(ctx)...).Infow(msg, keysAndValues...)
	case WarnLevel:
		l.logger.Sugar().With(l.extraFields(ctx)...).Warnw(msg, keysAndValues...)
	case ErrorLevel:
		l.logger.Sugar().With(l.extraFields(ctx)...).Errorw(msg, keysAndValues...)
	case PanicLevel:
		l.logger.Sugar().With(l.extraFields(ctx)...).Panicw(msg, keysAndValues...)
	case FatalLevel:
		l.logger.Sugar().With(l.extraFields(ctx)...).Fatalw(msg, keysAndValues...)
	}
}

// 提取 context 中的字段
func (l *zapLogger) extraFields(ctx context.Context) []interface{} {
	var fields []interface{}
	for _, key := range l.cfg.ExtraKeys {
		if v := ctx.Value(key); v != nil {
			fields = append(fields, zap.Any(key, v))
		}
	}
	return fields
}

func getZapEncoder() zapcore.Encoder {
	// 设置时间编码格式
	encodeTime := zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05.999999")

	// 配置编码器配置
	encoderCfg := zapcore.EncoderConfig{
		LevelKey:       "level",                       // 日志级别的键名，例如 "INFO", "ERROR"
		TimeKey:        "time",                        // 时间戳的键名，记录日志生成的时间
		StacktraceKey:  "stacktrace",                  // 堆栈跟踪的键名，记录日志产生时的堆栈信息
		CallerKey:      "file",                        // 调用者的键名，记录日志调用的位置 (文件名和行号)
		FunctionKey:    "function",                    // 函数名的键名，记录调用函数的名称
		MessageKey:     "msg",                         // 日志消息的键名，记录实际的日志内容
		LineEnding:     zapcore.DefaultLineEnding,     // 日志行的结束符，默认使用换行符
		EncodeCaller:   zapcore.ShortCallerEncoder,    // 调用者编码器，使用短格式 (文件名:行号)
		EncodeLevel:    zapcore.CapitalLevelEncoder,   // 日志级别编码器，使用大写格式，例如 "INFO", "ERROR"
		EncodeTime:     encodeTime,                    // 时间编码器，使用自定义时间格式 "2006-01-02 15:04:05.999999"
		EncodeDuration: zapcore.StringDurationEncoder, // 持续时间编码器，使用字符串格式记录持续时间
	}

	// 返回一个 JSON 编码器，用于将日志编码为 JSON 格式
	return zapcore.NewJSONEncoder(encoderCfg)
}

// getZapColorEncoder
func getZapColorEncoder() zapcore.Encoder {
	// 设置时间编码格式
	encodeTime := zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05.999999")

	// 配置编码器配置
	encoderCfg := zapcore.EncoderConfig{
		LevelKey:       "level",                          // 日志级别的键名，例如 "INFO", "ERROR"
		TimeKey:        "time",                           // 时间戳的键名，记录日志生成的时间
		CallerKey:      "file",                           // 调用者的键名，记录日志调用的位置 (文件名和行号)
		FunctionKey:    "function",                       // 函数名的键名，记录调用函数的名称
		MessageKey:     "msg",                            // 日志消息的键名，记录实际的日志内容
		StacktraceKey:  "stacktrace",                     // 堆栈跟踪的键名，记录日志产生时的堆栈信息
		LineEnding:     zapcore.DefaultLineEnding,        // 日志行的结束符，默认使用换行符
		EncodeCaller:   zapcore.ShortCallerEncoder,       // 调用者编码器，使用短格式 (文件名:行号)
		EncodeLevel:    zapcore.CapitalColorLevelEncoder, // 日志级别编码器，增加颜色
		EncodeTime:     encodeTime,                       // 时间编码器，使用自定义时间格式 "2006-01-02 15:04:05.999999"
		EncodeDuration: zapcore.StringDurationEncoder,    // 持续时间编码器，使用字符串格式记录持续时间
	}

	// 返回一个 JSON 编码器，用于将日志编码为 JSON 格式
	return zapcore.NewJSONEncoder(encoderCfg)
}

func getZapLogWriter(cfg *LoggerConfig, logOutputType string) (ws zapcore.WriteSyncer) {
	var w io.Writer
	if logOutputType == logOutputTypeStdout {
		w = os.Stdout
	} else {
		var err error
		director := strings.TrimSuffix(cfg.LogDir, "/") + "/" + time.Now().Format("20060102")
		if ok := gutils.FileExists(director); !ok {
			_ = os.MkdirAll(director, os.ModePerm)
		}
		var logFilename string
		switch logOutputType {
		case logOutputTypeDefaultFile:
			logFilename = fmt.Sprintf("%s%s", cfg.ServiceName, logOutputFileDefaultSuffix)
		case logOutputTypeWarnFatal:
			logFilename = fmt.Sprintf("%s%s", cfg.ServiceName, logOutputFileWarnFatalSuffix)
		default:
			logFilename = fmt.Sprintf("%s%s", cfg.ServiceName, logOutputFileDefaultSuffix)
		}
		rotator, err := rotatelogs.New(
			path.Join(strings.TrimSuffix(cfg.LogDir, "/"), "%Y%m%d", logFilename), // 分割后的文件名称
			rotatelogs.WithClock(rotatelogs.Local),                                // 使用本地时间
			rotatelogs.WithRotationTime(time.Hour*24),                             // 日志切割时间间隔
			rotatelogs.WithMaxAge(time.Hour*24*30),                                // 保留旧文件的最大时间
		)
		if err != nil {
			panic(err)
		}
		w = zapcore.AddSync(rotator)
	}

	flushInterval := 5 * time.Second
	if logOutputType == logOutputTypeStdout {
		flushInterval = 1 * time.Second
	}
	ws = &zapcore.BufferedWriteSyncer{
		WS:            zapcore.AddSync(w),
		Size:          256 * 1024,
		FlushInterval: flushInterval,
		Clock:         nil,
	}

	return ws
}

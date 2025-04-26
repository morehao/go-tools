package glog

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
)

// Logger 是一个封装 zap.Logger 的结构体
type zapLogger struct {
	logger *zap.Logger
	cfg    *LoggerConfig
}

type zapLoggerConfig struct {
	callerSkip      int
	fieldHookFunc   FieldHookFunc
	messageHookFunc MessageHookFunc
}

func getZapLogger(cfg *LoggerConfig, optCfg *optConfig) (*zap.Logger, error) {
	var core zapcore.Core

	switch cfg.Type {
	case WriterConsole:
		zapCfg := &zapLoggerConfig{
			callerSkip:      optCfg.callerSkip,
			fieldHookFunc:   optCfg.fieldHookFunc,
			messageHookFunc: optCfg.messageHookFunc,
		}
		core = zapcore.NewCore(
			getZapEncoder(zapCfg),
			zapcore.AddSync(os.Stdout),
			zapcore.Level(logLevelMap[cfg.Level]),
		)
	case WriterFile:
		// 创建日志目录
		if err := os.MkdirAll(cfg.Dir, 0755); err != nil {
			return nil, err
		}

		// 创建主日志文件
		mainFile, err := os.OpenFile(
			filepath.Join(cfg.Dir, fmt.Sprintf("%s.log", cfg.Service)),
			os.O_CREATE|os.O_APPEND|os.O_WRONLY,
			0644,
		)
		if err != nil {
			return nil, err
		}

		// 创建错误日志文件
		errorFile, err := os.OpenFile(
			filepath.Join(cfg.Dir, fmt.Sprintf("%s.error.log", cfg.Service)),
			os.O_CREATE|os.O_APPEND|os.O_WRONLY,
			0644,
		)
		if err != nil {
			return nil, err
		}

		// 创建主日志core
		mainCore := zapcore.NewCore(
			getZapEncoder(&zapLoggerConfig{
				callerSkip:      optCfg.callerSkip,
				fieldHookFunc:   optCfg.fieldHookFunc,
				messageHookFunc: optCfg.messageHookFunc,
			}),
			zapcore.AddSync(mainFile),
			zapcore.Level(logLevelMap[cfg.Level]),
		)

		// 创建错误日志core
		errorCore := zapcore.NewCore(
			getZapEncoder(&zapLoggerConfig{
				callerSkip:      optCfg.callerSkip,
				fieldHookFunc:   optCfg.fieldHookFunc,
				messageHookFunc: optCfg.messageHookFunc,
			}),
			zapcore.AddSync(errorFile),
			zapcore.ErrorLevel,
		)

		// 使用Tee将日志同时写入两个文件
		core = zapcore.NewTee(mainCore, errorCore)
	}

	// 如果有钩子函数，添加钩子
	if optCfg != nil && (optCfg.fieldHookFunc != nil || optCfg.messageHookFunc != nil) {
		zapCfg := &zapLoggerConfig{
			callerSkip:      optCfg.callerSkip,
			fieldHookFunc:   optCfg.fieldHookFunc,
			messageHookFunc: optCfg.messageHookFunc,
		}
		core = zapcore.NewTee(core, zapcore.NewCore(
			getZapEncoder(zapCfg),
			zapcore.AddSync(os.Stdout),
			zapcore.Level(logLevelMap[cfg.Level]),
		))
	}

	// 创建 logger
	logger := zap.New(core)

	// 如果设置了 callerSkip，添加 caller skip
	if optCfg != nil && optCfg.callerSkip > 0 {
		logger = logger.WithOptions(zap.AddCallerSkip(optCfg.callerSkip))
	}

	return logger, nil
}

func (l *zapLogger) Debug(ctx context.Context, args ...any) {
	l.ctxLog(DebugLevel, ctx, args...)
}

func (l *zapLogger) Debugf(ctx context.Context, format string, args ...any) {
	l.ctxLogf(DebugLevel, ctx, format, args...)
}

func (l *zapLogger) Debugw(ctx context.Context, msg string, keysAndValues ...any) {
	l.ctxLogw(DebugLevel, ctx, msg, keysAndValues...)
}

func (l *zapLogger) Info(ctx context.Context, args ...any) {
	l.ctxLog(InfoLevel, ctx, args...)
}

func (l *zapLogger) Infof(ctx context.Context, format string, args ...any) {
	l.ctxLogf(InfoLevel, ctx, format, args...)
}

func (l *zapLogger) Infow(ctx context.Context, msg string, keysAndValues ...any) {
	l.ctxLogw(InfoLevel, ctx, msg, keysAndValues...)
}

func (l *zapLogger) Warn(ctx context.Context, args ...any) {
	l.ctxLog(WarnLevel, ctx, args...)
}

func (l *zapLogger) Warnf(ctx context.Context, format string, args ...any) {
	l.ctxLogf(WarnLevel, ctx, format, args...)
}

func (l *zapLogger) Warnw(ctx context.Context, msg string, keysAndValues ...any) {
	l.ctxLogw(WarnLevel, ctx, msg, keysAndValues...)
}

func (l *zapLogger) Error(ctx context.Context, args ...any) {
	l.ctxLog(ErrorLevel, ctx, args...)
}

func (l *zapLogger) Errorf(ctx context.Context, format string, args ...any) {
	l.ctxLogf(ErrorLevel, ctx, format, args...)
}

func (l *zapLogger) Errorw(ctx context.Context, msg string, keysAndValues ...any) {
	l.ctxLogw(ErrorLevel, ctx, msg, keysAndValues...)
}

func (l *zapLogger) Panic(ctx context.Context, args ...any) {
	l.ctxLog(PanicLevel, ctx, args...)
}

func (l *zapLogger) Panicf(ctx context.Context, format string, args ...any) {
	l.ctxLogf(PanicLevel, ctx, format, args...)
}

func (l *zapLogger) Panicw(ctx context.Context, msg string, keysAndValues ...any) {
	l.ctxLogw(PanicLevel, ctx, msg, keysAndValues...)
}

func (l *zapLogger) Fatal(ctx context.Context, args ...any) {
	l.ctxLog(FatalLevel, ctx, args...)
}

func (l *zapLogger) Fatalf(ctx context.Context, format string, args ...any) {
	l.ctxLogf(PanicLevel, ctx, format, args...)
}

func (l *zapLogger) Fatalw(ctx context.Context, msg string, keysAndValues ...any) {
	l.ctxLogw(FatalLevel, ctx, msg, keysAndValues...)
}

func (l *zapLogger) getLogger(opts ...Option) (Logger, error) {
	cfg := &optConfig{}
	for _, opt := range opts {
		opt.apply(cfg)
	}

	// 创建新的 logger
	logger := l.logger

	// 如果设置了 callerSkip，添加 caller skip
	if cfg.callerSkip > 0 {
		logger = logger.WithOptions(zap.AddCallerSkip(cfg.callerSkip))
	}

	return &zapLogger{
		logger: logger,
		cfg:    l.cfg,
	}, nil
}

func (l *zapLogger) Close() {
	_ = l.logger.Sync()
}

func (l *zapLogger) ctxLog(level Level, ctx context.Context, args ...any) {
	if nilCtx(ctx) || skipLog(ctx) {
		return
	}

	// 执行钩子函数
	executeHooks(ctx, level, fmt.Sprint(args...))

	// 原有的日志记录逻辑
	fields := l.extraFields(ctx)
	switch level {
	case DebugLevel:
		l.logger.Sugar().With(fields...).Debug(args...)
	case InfoLevel:
		l.logger.Sugar().With(fields...).Info(args...)
	case WarnLevel:
		l.logger.Sugar().With(fields...).Warn(args...)
	case ErrorLevel:
		l.logger.Sugar().With(fields...).Error(args...)
	case PanicLevel:
		l.logger.Sugar().With(fields...).Panic(args...)
	case FatalLevel:
		l.logger.Sugar().With(fields...).Fatal(args...)
	}
}

func (l *zapLogger) ctxLogf(level Level, ctx context.Context, format string, args ...any) {
	if nilCtx(ctx) || skipLog(ctx) {
		return
	}

	// 执行钩子函数
	executeHooks(ctx, level, fmt.Sprintf(format, args...))

	// 原有的日志记录逻辑
	fields := l.extraFields(ctx)
	switch level {
	case DebugLevel:
		l.logger.Sugar().With(fields...).Debugf(format, args...)
	case InfoLevel:
		l.logger.Sugar().With(fields...).Infof(format, args...)
	case WarnLevel:
		l.logger.Sugar().With(fields...).Warnf(format, args...)
	case ErrorLevel:
		l.logger.Sugar().With(fields...).Errorf(format, args...)
	case PanicLevel:
		l.logger.Sugar().With(fields...).Panicf(format, args...)
	case FatalLevel:
		l.logger.Sugar().With(fields...).Fatalf(format, args...)
	}
}

func (l *zapLogger) ctxLogw(level Level, ctx context.Context, msg string, keysAndValues ...any) {
	if nilCtx(ctx) || skipLog(ctx) {
		return
	}

	// 将keysAndValues转换为Field切片
	var fields []Field
	for i := 0; i < len(keysAndValues); i += 2 {
		if i+1 < len(keysAndValues) {
			fields = append(fields, Field{
				Key:   fmt.Sprint(keysAndValues[i]),
				Value: keysAndValues[i+1],
			})
		}
	}

	// 执行钩子函数
	executeHooks(ctx, level, msg, fields...)

	// 原有的日志记录逻辑
	extraFields := l.extraFields(ctx)
	switch level {
	case DebugLevel:
		l.logger.Sugar().With(append(extraFields, keysAndValues...)...).Debugw(msg)
	case InfoLevel:
		l.logger.Sugar().With(append(extraFields, keysAndValues...)...).Infow(msg)
	case WarnLevel:
		l.logger.Sugar().With(append(extraFields, keysAndValues...)...).Warnw(msg)
	case ErrorLevel:
		l.logger.Sugar().With(append(extraFields, keysAndValues...)...).Errorw(msg)
	case PanicLevel:
		l.logger.Sugar().With(append(extraFields, keysAndValues...)...).Panicw(msg)
	case FatalLevel:
		l.logger.Sugar().With(append(extraFields, keysAndValues...)...).Fatalw(msg)
	}
}

// 提取 context 中的字段
func (l *zapLogger) extraFields(ctx context.Context) []any {
	var fields []any
	for _, key := range l.cfg.ExtraKeys {
		if v := ctx.Value(key); v != nil {
			fields = append(fields, zap.Any(key, v))
		}
	}
	return fields
}

type gZapEncoder struct {
	fieldHookFunc   FieldHookFunc
	messageHookFunc MessageHookFunc
	zapcore.Encoder
}

func (enc *gZapEncoder) EncodeEntry(ent zapcore.Entry, fields []zapcore.Field) (*buffer.Buffer, error) {
	// 转换 zapcore.Field 到 Field
	var convertedFields []Field
	for _, f := range fields {
		convertedFields = append(convertedFields, Field{
			Key:   f.Key,
			Value: f.Interface,
		})
	}

	// 调用字段钩子函数
	if enc.fieldHookFunc != nil {
		enc.fieldHookFunc(convertedFields)
	}

	// 调用消息钩子函数
	if enc.messageHookFunc != nil {
		ent.Message = enc.messageHookFunc(ent.Message)
	}

	return enc.Encoder.EncodeEntry(ent, fields)
}

func getZapEncoder(cfg *zapLoggerConfig) zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	encoderConfig.EncodeDuration = zapcore.StringDurationEncoder
	encoderConfig.EncodeName = zapcore.FullNameEncoder

	encoder := zapcore.NewJSONEncoder(encoderConfig)

	// 如果配置了字段钩子函数或消息钩子函数，则使用自定义编码器
	if cfg != nil && (cfg.fieldHookFunc != nil || cfg.messageHookFunc != nil) {
		encoder = &gZapEncoder{
			Encoder:         encoder,
			fieldHookFunc:   cfg.fieldHookFunc,
			messageHookFunc: cfg.messageHookFunc,
		}
	}

	return encoder
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
		EncodeTime:     encodeTime,                       // 时间编码器，使用自定义时间格式
		EncodeDuration: zapcore.StringDurationEncoder,    // 持续时间编码器，使用字符串格式
	}

	// 返回一个 JSON 编码器，用于将日志编码为 JSON 格式
	return zapcore.NewJSONEncoder(encoderCfg)
}

// timeWriter 实现按时间切割的writer
type timeWriter struct {
	filePath string
	file     *os.File
	lastTime time.Time
}

func (w *timeWriter) Write(p []byte) (n int, err error) {
	now := time.Now()
	if w.file == nil || !w.isSameTimeUnit(now) {
		if err := w.openFile(now); err != nil {
			return 0, err
		}
	}
	return w.file.Write(p)
}

func (w *timeWriter) isSameTimeUnit(t time.Time) bool {
	if w.file == nil {
		return false
	}
	// 判断是否在同一天
	return t.Year() == w.lastTime.Year() && t.Month() == w.lastTime.Month() && t.Day() == w.lastTime.Day()
}

func (w *timeWriter) openFile(t time.Time) error {
	if w.file != nil {
		w.file.Close()
	}

	// 创建日志目录
	dir := filepath.Dir(w.filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// 生成日志文件名
	filename := fmt.Sprintf("%s.%s.log", w.filePath, t.Format("2006-01-02"))

	file, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	w.file = file
	w.lastTime = t

	return nil
}

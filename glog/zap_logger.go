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
	// 创建基础配置
	zapCfg := &zapLoggerConfig{
		callerSkip:      optCfg.callerSkip,
		fieldHookFunc:   optCfg.fieldHookFunc,
		messageHookFunc: optCfg.messageHookFunc,
	}

	// 创建编码器
	encoder := getZapEncoder(zapCfg)

	// 创建控制台输出
	consoleCore := zapcore.NewCore(
		encoder,
		zapcore.AddSync(os.Stdout),
		logLevelMap[cfg.Level],
	)

	var cores []zapcore.Core

	// 根据配置类型添加其他输出
	switch cfg.Type {
	case WriterConsole:
		cores = append(cores, consoleCore)
	case WriterFile:
		// 创建日志目录
		if err := os.MkdirAll(cfg.Dir, 0755); err != nil {
			return nil, err
		}

		// 创建主日志文件
		mainFile, err := os.OpenFile(
			filepath.Join(cfg.Dir, fmt.Sprintf("%s.log", cfg.service)),
			os.O_CREATE|os.O_APPEND|os.O_WRONLY,
			0644,
		)
		if err != nil {
			return nil, err
		}

		// 创建错误日志文件
		errorFile, err := os.OpenFile(
			filepath.Join(cfg.Dir, fmt.Sprintf("%s.error.log", cfg.service)),
			os.O_CREATE|os.O_APPEND|os.O_WRONLY,
			0644,
		)
		if err != nil {
			return nil, err
		}

		// 创建主日志core
		mainCore := zapcore.NewCore(
			encoder,
			zapcore.AddSync(mainFile),
			logLevelMap[cfg.Level],
		)

		// 创建错误日志core
		errorCore := zapcore.NewCore(
			encoder,
			zapcore.AddSync(errorFile),
			zapcore.ErrorLevel,
		)
		cores = append(cores, consoleCore, mainCore, errorCore)
	}

	// 使用Tee将日志同时写入所有输出
	core := zapcore.NewTee(cores...)

	// 创建 logger，添加 caller 选项
	logger := zap.New(core, zap.Development(), zap.AddCaller(), zap.AddStacktrace(zapcore.PanicLevel))

	// 如果设置了 callerSkip，添加 caller skip
	if optCfg.callerSkip > 0 {
		logger = logger.WithOptions(zap.AddCallerSkip(optCfg.callerSkip))
	}

	return logger, nil
}

func (l *zapLogger) Debug(ctx context.Context, kvs ...any) {
	l.ctxLog(DebugLevel, ctx, kvs...)
}

func (l *zapLogger) Debugf(ctx context.Context, format string, kvs ...any) {
	l.ctxLogf(DebugLevel, ctx, format, kvs...)
}

func (l *zapLogger) Debugw(ctx context.Context, msg string, kvs ...any) {
	l.ctxLogw(DebugLevel, ctx, msg, kvs...)
}

func (l *zapLogger) Info(ctx context.Context, kvs ...any) {
	l.ctxLog(InfoLevel, ctx, kvs...)
}

func (l *zapLogger) Infof(ctx context.Context, format string, kvs ...any) {
	l.ctxLogf(InfoLevel, ctx, format, kvs...)
}

func (l *zapLogger) Infow(ctx context.Context, msg string, kvs ...any) {
	l.ctxLogw(InfoLevel, ctx, msg, kvs...)
}

func (l *zapLogger) Warn(ctx context.Context, kvs ...any) {
	l.ctxLog(WarnLevel, ctx, kvs...)
}

func (l *zapLogger) Warnf(ctx context.Context, format string, kvs ...any) {
	l.ctxLogf(WarnLevel, ctx, format, kvs...)
}

func (l *zapLogger) Warnw(ctx context.Context, msg string, kvs ...any) {
	l.ctxLogw(WarnLevel, ctx, msg, kvs...)
}

func (l *zapLogger) Error(ctx context.Context, kvs ...any) {
	l.ctxLog(ErrorLevel, ctx, kvs...)
}

func (l *zapLogger) Errorf(ctx context.Context, format string, kvs ...any) {
	l.ctxLogf(ErrorLevel, ctx, format, kvs...)
}

func (l *zapLogger) Errorw(ctx context.Context, msg string, kvs ...any) {
	l.ctxLogw(ErrorLevel, ctx, msg, kvs...)
}

func (l *zapLogger) Panic(ctx context.Context, kvs ...any) {
	l.ctxLog(PanicLevel, ctx, kvs...)
}

func (l *zapLogger) Panicf(ctx context.Context, format string, kvs ...any) {
	l.ctxLogf(PanicLevel, ctx, format, kvs...)
}

func (l *zapLogger) Panicw(ctx context.Context, msg string, kvs ...any) {
	l.ctxLogw(PanicLevel, ctx, msg, kvs...)
}

func (l *zapLogger) Fatal(ctx context.Context, kvs ...any) {
	l.ctxLog(FatalLevel, ctx, kvs...)
}

func (l *zapLogger) Fatalf(ctx context.Context, format string, kvs ...any) {
	l.ctxLogf(FatalLevel, ctx, format, kvs...)
}

func (l *zapLogger) Fatalw(ctx context.Context, msg string, kvs ...any) {
	l.ctxLogw(FatalLevel, ctx, msg, kvs...)
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

func (l *zapLogger) ctxLog(level Level, ctx context.Context, kvs ...any) {
	if nilCtx(ctx) || skipLog(ctx) {
		return
	}

	fields := convertKvsToFields(kvs...)

	// 执行钩子函数
	msg := ""
	if len(fields) > 0 {
		msg = fmt.Sprint(fields[0].Value)
	}
	executeHooks(ctx, level, msg, fields...)

	// 获取上下文字段
	zapFields := l.extraFields(ctx)

	// 将 Field 转换为 zap.Field
	for _, f := range fields {
		zapFields = append(zapFields, zap.Any(f.Key, f.Value))
	}

	// 记录日志
	switch level {
	case DebugLevel:
		l.logger.Debug(msg, zapFields...)
	case InfoLevel:
		l.logger.Info(msg, zapFields...)
	case WarnLevel:
		l.logger.Warn(msg, zapFields...)
	case ErrorLevel:
		l.logger.Error(msg, zapFields...)
	case PanicLevel:
		l.logger.Panic(msg, zapFields...)
	case FatalLevel:
		l.logger.Fatal(msg, zapFields...)
	}
}

func (l *zapLogger) ctxLogf(level Level, ctx context.Context, format string, kvs ...any) {
	if nilCtx(ctx) || skipLog(ctx) {
		return
	}

	fields := convertKvsToFields(kvs...)

	// 执行钩子函数
	msg := ""
	if len(fields) > 0 {
		msg = fmt.Sprintf(format, fields[0].Value)
	}
	executeHooks(ctx, level, msg, fields...)

	// 获取上下文字段
	zapFields := l.extraFields(ctx)

	// 将 Field 转换为 zap.Field
	for _, f := range fields {
		zapFields = append(zapFields, zap.Any(f.Key, f.Value))
	}

	// 记录日志
	switch level {
	case DebugLevel:
		l.logger.Debug(msg, zapFields...)
	case InfoLevel:
		l.logger.Info(msg, zapFields...)
	case WarnLevel:
		l.logger.Warn(msg, zapFields...)
	case ErrorLevel:
		l.logger.Error(msg, zapFields...)
	case PanicLevel:
		l.logger.Panic(msg, zapFields...)
	case FatalLevel:
		l.logger.Fatal(msg, zapFields...)
	}
}

func (l *zapLogger) ctxLogw(level Level, ctx context.Context, msg string, kvs ...any) {
	if nilCtx(ctx) || skipLog(ctx) {
		return
	}

	fields := convertKvsToFields(kvs...)

	// 执行钩子函数
	executeHooks(ctx, level, msg, fields...)

	// 获取上下文字段
	zapFields := l.extraFields(ctx)

	// 将 Field 转换为 zap.Field
	for _, f := range fields {
		zapFields = append(zapFields, zap.Any(f.Key, f.Value))
	}

	// 记录日志
	switch level {
	case DebugLevel:
		l.logger.Debug(msg, zapFields...)
	case InfoLevel:
		l.logger.Info(msg, zapFields...)
	case WarnLevel:
		l.logger.Warn(msg, zapFields...)
	case ErrorLevel:
		l.logger.Error(msg, zapFields...)
	case PanicLevel:
		l.logger.Panic(msg, zapFields...)
	case FatalLevel:
		l.logger.Fatal(msg, zapFields...)
	}
}

// 提取 context 中的字段
func (l *zapLogger) extraFields(ctx context.Context) []zap.Field {
	var fields []zap.Field
	// 添加 writer 类型字段
	fields = append(fields, zap.String("writer", string(l.cfg.Type)))

	for _, key := range l.cfg.ExtraKeys {
		if v := ctx.Value(key); v != nil {
			fields = append(fields, zap.Any(key, v))
		}
	}
	return fields
}

type gZapEncoder struct {
	zapcore.Encoder
	fieldHookFunc   FieldHookFunc
	messageHookFunc MessageHookFunc
}

func (enc *gZapEncoder) EncodeEntry(ent zapcore.Entry, fields []zapcore.Field) (*buffer.Buffer, error) {
	// 转换 zapcore.Field 到 Field
	convertedFields := make([]Field, 0, len(fields))
	for _, f := range fields {
		convertedFields = append(convertedFields, Field{
			Key:   f.Key,
			Value: f.Interface,
		})
	}

	// 将修改后的字段转换回 zapcore.Field
	modifiedFields := make([]zapcore.Field, 0, len(convertedFields))
	for _, f := range convertedFields {
		modifiedFields = append(modifiedFields, zapcore.Field{
			Key:       f.Key,
			Type:      zapcore.ReflectType,
			Interface: f.Value,
		})
	}

	// 使用修改后的字段进行编码
	return enc.Encoder.EncodeEntry(ent, modifiedFields)
}

func getZapEncoder(cfg *zapLoggerConfig) zapcore.Encoder {
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	encoder := zapcore.NewJSONEncoder(encoderCfg)
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
		CallerKey:      "caller",                         // 调用者的键名，记录日志调用的位置 (文件名和行号)
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

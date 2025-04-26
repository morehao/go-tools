package glog

import (
	"context"
	"fmt"

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
		getZapStandoutWriter(),
		logLevelMap[cfg.Level],
	)

	var cores []zapcore.Core

	// 根据配置类型添加其他输出
	switch cfg.Writer {
	case WriterConsole:
		cores = append(cores, consoleCore)
	case WriterFile:
		defaultWriter, getDefaultWriterErr := getZapFileWriter(cfg, "_full.log")
		if getDefaultWriterErr != nil {
			return nil, getDefaultWriterErr
		}
		wfWriter, getWfWriterErr := getZapFileWriter(cfg, "_wf.log")
		if getWfWriterErr != nil {
			return nil, getWfWriterErr
		}

		// 创建默认日志core
		defaultCore := zapcore.NewCore(
			encoder,
			defaultWriter,
			logLevelMap[cfg.Level],
		)

		// 创建wf日志core
		wfCore := zapcore.NewCore(
			encoder,
			wfWriter,
			zapcore.WarnLevel,
		)
		cores = append(cores, consoleCore, defaultCore, wfCore)
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
	fields = append(fields, zap.String("writer", string(l.cfg.Writer)))

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

func (enc *gZapEncoder) Clone() zapcore.Encoder {
	encoderClone := enc.Encoder.Clone()
	return &gZapEncoder{
		Encoder:         encoderClone,
		fieldHookFunc:   enc.fieldHookFunc,
		messageHookFunc: enc.messageHookFunc,
	}
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

	// 执行字段钩子函数
	if enc.fieldHookFunc != nil {
		enc.fieldHookFunc(convertedFields)
	}

	// 执行消息钩子函数
	if enc.messageHookFunc != nil {
		ent.Message = enc.messageHookFunc(ent.Message)
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

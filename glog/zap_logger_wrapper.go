package glog

import (
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
)

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

func getZapStandoutWriter() zapcore.WriteSyncer {
	return os.Stdout
}

func getZapFileWriter(cfg *LoggerConfig, fileSuffix string) (zapcore.WriteSyncer, error) {
	// 目录始终按天组织
	dir := strings.TrimSuffix(cfg.Dir, "/") + "/" + time.Now().Format("20060102")
	if ok := fileExists(dir); !ok {
		_ = os.MkdirAll(dir, os.ModePerm)
	}

	// 根据 RotateUnit 确定日志文件名的时间格式
	var timeFormat string
	switch cfg.RotateUnit {
	case RotateUnitHour:
		timeFormat = "15" // 只包含小时
	default:
		timeFormat = "" // 不包含时间
	}

	// 构建日志文件名
	var logFilename string
	if timeFormat != "" {
		logFilename = fmt.Sprintf("%s_%s_%s.log", cfg.service, fileSuffix, time.Now().Format(timeFormat))
	} else {
		logFilename = fmt.Sprintf("%s_%s.log", cfg.service, fileSuffix)
	}

	logFilepath := path.Join(dir, logFilename)

	// 打开日志文件
	file, openErr := os.OpenFile(logFilepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if openErr != nil {
		return nil, openErr
	}

	// 创建带缓冲的写入器
	writer := &zapcore.BufferedWriteSyncer{
		WS:            zapcore.AddSync(file),
		Size:          256 * 1024,
		FlushInterval: time.Second * 5,
		Clock:         nil,
	}

	return writer, nil
}

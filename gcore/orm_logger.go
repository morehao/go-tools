package gcore

import (
	"context"
	"fmt"
	"github.com/morehao/go-tools/glog"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
	"time"
)

type ormLogger struct {
	Service   string
	Address   string
	Database  string
	MaxSqlLen int
	Logger    glog.Logger
}

type ormConfig struct {
	Service   string
	Address   string
	Database  string
	MaxSqlLen int
}

func newOrmLogger(cfg *ormConfig) *ormLogger {
	s := cfg.Service
	if cfg.Service == "" {
		s = cfg.Database
	}

	return &ormLogger{
		Service:   s,
		Address:   cfg.Address,
		Database:  cfg.Database,
		MaxSqlLen: cfg.MaxSqlLen,
		Logger:    glog.GetLogger(glog.WithZapOptions(zap.AddCallerSkip(2))),
	}
}

// LogMode log mode
func (l *ormLogger) LogMode(level logger.LogLevel) logger.Interface {
	return l
}

// Info print info
func (l *ormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	m := fmt.Sprintf(msg, append([]interface{}{utils.FileWithLineNum()}, data...)...)
	l.Logger.Debugw(ctx, m, l.commonFields(ctx)...)
}

// Warn print warn messages
func (l *ormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	m := fmt.Sprintf(msg, append([]interface{}{utils.FileWithLineNum()}, data...)...)
	l.Logger.Warnw(ctx, m, l.commonFields(ctx)...)
}

// Error print error messages
func (l *ormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	m := fmt.Sprintf(msg, append([]interface{}{utils.FileWithLineNum()}, data...)...)
	l.Logger.Errorw(ctx, m, l.commonFields(ctx)...)
}

// Trace print sql message
func (l *ormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	end := time.Now()
	elapsed := end.Sub(begin)
	cost := float64(elapsed.Nanoseconds()/1e4) / 100.0

	msg := "sql execute success"
	ralCode := -0
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		// 没有找到记录不统计在请求错误中
		msg = err.Error()
		ralCode = -1
	}
	sql, rows := fc()

	fileLineNum := utils.FileWithLineNum()
	fields := l.commonFields(ctx)
	fields = append(fields,
		// zap.String("msg", msg),
		zap.Int64("affectedRow", rows),
		zap.String("requestEndTime", glog.FormatRequestTime(end)),
		zap.String("requestStartTime", glog.FormatRequestTime(begin)),
		zap.String("file", fileLineNum),
		zap.Float64("cost", cost),
		zap.Int("ralCode", ralCode),
		zap.String("sql", sql),
	)

	l.Logger.Infow(ctx, msg, fields...)
}

func (l *ormLogger) commonFields(ctx context.Context) []interface{} {
	fields := []interface{}{
		glog.KeyProto, "mysql",
		"service", l.Service,
		"address", l.Address,
		"database", l.Database,
	}
	return fields
}

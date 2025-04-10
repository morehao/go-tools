package dbClient

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/morehao/go-tools/glog"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
)

type MysqlConfig struct {
	Service       string        `yaml:"service"`        // 服务名
	Addr          string        `yaml:"addr"`           // 地址
	Database      string        `yaml:"database"`       // 数据库名
	User          string        `yaml:"user"`           // 用户名
	Password      string        `yaml:"password"`       // 密码
	Charset       string        `yaml:"charset"`        // 字符集
	Timeout       time.Duration `yaml:"timeout"`        // 连接超时
	ReadTimeout   time.Duration `yaml:"read_timeout"`   // 读取超时
	WriteTimeout  time.Duration `yaml:"write_timeout"`  // 写入超时
	SlowThreshold time.Duration `yaml:"slow_threshold"` // 慢SQL阈值
	MaxSqlLen     int           `yaml:"max_sql_len"`    // 日志最大SQL长度
}

func InitMysql(cfg MysqlConfig) (*gorm.DB, error) {
	dns := cfg.buildDns()
	customLogger, newLogErr := newOrmLogger(&ormConfig{
		Service:   cfg.Service,
		Addr:      cfg.Addr,
		Database:  cfg.Database,
		MaxSqlLen: cfg.MaxSqlLen,
	})
	if newLogErr != nil {
		return nil, newLogErr
	}
	db, err := gorm.Open(mysql.Open(dns), &gorm.Config{
		Logger: customLogger,
	})
	if err != nil {
		return nil, err
	}
	return db, nil
}

func (cfg *MysqlConfig) buildDns() string {
	dns := fmt.Sprintf("%s:%s@tcp(%s)/%s?&parseTime=True&loc=Local&timeout=%s&readTimeout=%s&writeTimeout=%s",
		cfg.User, cfg.Password, cfg.Addr, cfg.Database,
		cfg.Timeout, cfg.ReadTimeout, cfg.WriteTimeout)
	var charset = "utf8mb4"
	if cfg.Charset != "" {
		charset = cfg.Charset
	}
	dns += "&charset=" + charset
	return dns
}

type ormLogger struct {
	Service       string
	Addr          string
	Database      string
	MaxSqlLen     int
	SlowThreshold time.Duration
	Logger        glog.Logger
}

type ormConfig struct {
	Service   string
	Addr      string
	Database  string
	MaxSqlLen int
}

func newOrmLogger(cfg *ormConfig) (*ormLogger, error) {
	s := cfg.Service
	if cfg.Service == "" {
		s = cfg.Database
	}
	l, err := glog.GetLogger(glog.WithZapOptions(zap.AddCallerSkip(2)))
	if err != nil {
		return nil, err
	}
	return &ormLogger{
		Service:   s,
		Addr:      cfg.Addr,
		Database:  cfg.Database,
		MaxSqlLen: cfg.MaxSqlLen,
		Logger:    l,
	}, nil
}

// LogMode log mode
func (l *ormLogger) LogMode(level logger.LogLevel) logger.Interface {
	return l
}

// Info print info
func (l *ormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	formatMsg := fmt.Sprintf(msg, append([]interface{}{utils.FileWithLineNum()}, data...)...)
	l.Logger.Infow(ctx, formatMsg, l.commonFields(ctx)...)
}

// Warn print warn messages
func (l *ormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	formatMsg := fmt.Sprintf(msg, append([]interface{}{utils.FileWithLineNum()}, data...)...)
	l.Logger.Warnw(ctx, formatMsg, l.commonFields(ctx)...)
}

// Error print error messages
func (l *ormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	formatMsg := fmt.Sprintf(msg, append([]interface{}{utils.FileWithLineNum()}, data...)...)
	l.Logger.Errorw(ctx, formatMsg, l.commonFields(ctx)...)
}

// Trace print sql message
func (l *ormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	end := time.Now()
	cost := glog.GetRequestCost(begin, end)

	msg := "sql execute success"
	var ralCode int
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		// 过滤未找到数据的错误
		msg = err.Error()
		ralCode = -1
	}
	sql, rows := fc()
	if len(sql) > l.MaxSqlLen && l.MaxSqlLen > 0 {
		sql = sql[:l.MaxSqlLen]
	}

	fileLineNum := utils.FileWithLineNum()
	fields := l.commonFields(ctx)
	fields = append(fields,
		glog.KeyAffectedRows, rows,
		glog.KeyRequestStartTime, glog.FormatRequestTime(begin),
		glog.KeyRequestEndTime, glog.FormatRequestTime(end),
		glog.KeyFile, fileLineNum,
		glog.KeyCost, cost,
		glog.KeyRalCode, ralCode,
		glog.KeySql, sql,
	)

	if l.SlowThreshold > 0 && cost >= float64(l.SlowThreshold/time.Millisecond) {
		msg = "slow sql"
		l.Logger.Warnw(ctx, msg, fields...)
	} else {
		l.Logger.Infow(ctx, msg, fields...)
	}
}

func (l *ormLogger) commonFields(ctx context.Context) []interface{} {
	fields := []interface{}{
		glog.KeyProto, glog.ValueProtoMysql,
		glog.KeyService, l.Service,
		glog.KeyAddr, l.Addr,
		glog.KeyDatabase, l.Database,
	}
	return fields
}

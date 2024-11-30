package dbClient

import (
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
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
	logger, newLogErr := newOrmLogger(&ormConfig{
		Service:   cfg.Service,
		Addr:      cfg.Addr,
		Database:  cfg.Database,
		MaxSqlLen: cfg.MaxSqlLen,
	})
	if newLogErr != nil {
		return nil, newLogErr
	}
	db, err := gorm.Open(mysql.Open(dns), &gorm.Config{
		Logger: logger,
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

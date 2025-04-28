package dbmysql

import (
	"fmt"
	"sync"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	dbMap = map[string]*gorm.DB{}
	lock  sync.RWMutex
)

type MysqlConfig struct {
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
	if cfg.Database == "" {
		return nil, fmt.Errorf("database name is empty")
	}
	dns := cfg.buildDns()
	customLogger, newLogErr := newOrmLogger(&ormConfig{
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
	lock.Lock()
	defer lock.Unlock()
	dbMap[cfg.Database] = db
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

func InitMultiMysql(configs []MysqlConfig) error {
	if len(configs) == 0 {
		return fmt.Errorf("mysql configs is empty")
	}
	for _, cfg := range configs {
		_, err := InitMysql(cfg)
		if err != nil {
			return err
		}
	}
	return nil
}

func GetDB(database string) *gorm.DB {
	lock.RLock()
	defer lock.RUnlock()
	return dbMap[database]
}

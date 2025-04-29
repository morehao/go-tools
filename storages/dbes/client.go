package dbes

import (
	"github.com/elastic/go-elasticsearch/v8"
)

type ESConfig struct {
	Service  string `yaml:"service"`  // 服务名称
	Addr     string `yaml:"addr"`     // 地址
	User     string `yaml:"user"`     // 用户名
	Password string `yaml:"password"` // 密码
}

func InitES(cfg ESConfig) (*elasticsearch.Client, *elasticsearch.TypedClient, error) {
	customLogger, getLoggerErr := newEsLogger(&cfg)
	if getLoggerErr != nil {
		return nil, nil, getLoggerErr
	}
	commonCfg := elasticsearch.Config{
		Addresses: []string{cfg.Addr},
		Username:  cfg.User,
		Password:  cfg.Password,
		Logger:    customLogger,
	}
	simpleClient, newSimpleClientErr := elasticsearch.NewClient(commonCfg)
	if newSimpleClientErr != nil {
		return nil, nil, newSimpleClientErr
	}
	typedClient, newTypedClientErr := elasticsearch.NewTypedClient(commonCfg)
	if newTypedClientErr != nil {
		return nil, nil, newTypedClientErr
	}
	return simpleClient, typedClient, nil
}

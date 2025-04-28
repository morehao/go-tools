package dbes

import (
	"fmt"
	"sync"

	"github.com/elastic/go-elasticsearch/v8"
)

var (
	simpleClientMap = map[string]*elasticsearch.Client{}
	typedClientMap  = map[string]*elasticsearch.TypedClient{}
	lock            sync.RWMutex
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
	lock.Lock()
	defer lock.Unlock()
	simpleClientMap[cfg.Service] = simpleClient
	typedClientMap[cfg.Service] = typedClient
	return simpleClient, typedClient, nil
}

func InitMultiES(configs []ESConfig) error {
	if len(configs) == 0 {
		return fmt.Errorf("es config is empty")
	}
	for _, cfg := range configs {
		_, _, err := InitES(cfg)
		if err != nil {
			return err
		}
	}

	return nil
}

func GetSimpleClient(service string) *elasticsearch.Client {
	lock.RLock()
	defer lock.RUnlock()
	return simpleClientMap[service]
}

func GetTypedClient(service string) *elasticsearch.TypedClient {
	lock.RLock()
	defer lock.RUnlock()
	return typedClientMap[service]
}

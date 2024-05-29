package conf

import (
	"fmt"
	"testing"

	"github.com/morehao/go-tools/utils"
)

func TestLoadConfig(t *testing.T) {
	type MysqlConfig struct {
		Host     string `yaml:"host"`
		Database string `yaml:"database"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
	}
	type Config struct {
		Mysql MysqlConfig `yaml:"mysql"`
	}
	var config Config
	LoadConfig("", "config_example.yaml", &config)
	fmt.Println(utils.ToJson(config))
}

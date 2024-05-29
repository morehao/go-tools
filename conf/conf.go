package conf

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

func LoadConfig(confDir, filename string, dest interface{}) {
	path := filepath.Join(".", confDir, filename)
	if fileContent, readErr := os.ReadFile(path); readErr != nil {
		panic(fmt.Sprintf("load config fail, path: %s error: %s", path, readErr.Error()))
	} else {
		if err := yaml.Unmarshal(fileContent, dest); err != nil {
			panic(fmt.Sprintf("unmarshal config fail, path: %s error: %s", path, err.Error()))
		}
	}
}

package conf

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

func LoadConfig(filepath string, dest interface{}) {
	if fileContent, readErr := os.ReadFile(filepath); readErr != nil {
		panic(fmt.Sprintf("load config fail, filepath: %s error: %s", filepath, readErr.Error()))
	} else {
		if err := yaml.Unmarshal(fileContent, dest); err != nil {
			panic(fmt.Sprintf("unmarshal config fail, filepath: %s error: %s", filepath, err.Error()))
		}
	}
}

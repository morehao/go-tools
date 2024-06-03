package utils

import (
	"os"
	"strings"
)

// GetFileSuffix 获取文件后缀
func GetFileSuffix(name string) string {
	// 获取文件后缀
	return name[strings.LastIndex(name, "."):]
}

// FileExists 检查文件是否存在
func FileExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

func CreateDir(dir string) error {
	if exist := FileExists(dir); !exist {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			return err
		}
	}
	return nil
}

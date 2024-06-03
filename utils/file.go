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
func FileExists(path string) (bool, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false, nil
	}
	return true, nil
}

func CreateDir(dir string) error {
	if exist, err := FileExists(dir); err != nil {
		return err
	} else if !exist {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			return err
		}
	}
	return nil
}

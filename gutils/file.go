package gutils

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
func FileExists(filepath string) bool {
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		return false
	}
	return true
}

func PathExists(path string) (bool, error) {
	fi, err := os.Stat(path)
	if err == nil {
		if fi.IsDir() {
			return true, nil
		}
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func CreateDir(dir string) error {
	if exist := FileExists(dir); !exist {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			return err
		}
	}
	return nil
}

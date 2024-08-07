package gutils

import (
	"regexp"
	"strings"
)

// SnakeToPascal 蛇形转大驼峰
func SnakeToPascal(s string) string {
	if s == "" {
		return ""
	}
	parts := strings.Split(s, "_")

	for i := range parts {
		parts[i] = strings.ToTitle(parts[i][:1]) + parts[i][1:]
	}

	return strings.Join(parts, "")
}

// SnakeToLowerCamel 蛇形转小驼峰
func SnakeToLowerCamel(s string) string {
	if s == "" {
		return ""
	}
	parts := strings.Split(s, "_")

	for i := range parts {
		if i == 0 {
			parts[i] = strings.ToLower(parts[i][:1]) + parts[i][1:]
		} else {
			parts[i] = strings.ToTitle(parts[i][:1]) + parts[i][1:]
		}
	}

	return strings.Join(parts, "")
}

// FirstLetterToUpper 首字母大写
func FirstLetterToUpper(s string) string {
	if len(s) == 0 {
		return ""
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

func FirstLetterToLower(s string) string {
	if len(s) == 0 {
		return ""
	}
	return strings.ToLower(s[:1]) + s[1:]
}

// ReplaceIdToID 将id、Id、iD替换为Id
func ReplaceIdToID(str string) string {
	s := strings.Replace(str, "id", "ID", -1)
	s = strings.Replace(s, "Id", "ID", -1)
	s = strings.Replace(s, "iD", "ID", -1)
	return s
}

func Trim(str string) string {
	if len(str) == 0 {
		return ""
	}
	s := strings.Replace(str, " ", "", -1)
	// 替换所有空白字符（包括空格、制表符、换行符等）
	return regexp.MustCompile(`\s`).ReplaceAllString(s, "")
}

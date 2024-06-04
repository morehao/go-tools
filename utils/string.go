package utils

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

// FirstLetterToUpper 首字母大写
func FirstLetterToUpper(s string) string {
	if len(s) == 0 {
		return ""
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

func Trim(str string) string {
	if len(str) == 0 {
		return ""
	}
	s := strings.Replace(str, " ", "", -1)
	// 替换所有空白字符（包括空格、制表符、换行符等）
	return regexp.MustCompile(`\s`).ReplaceAllString(s, "")
}

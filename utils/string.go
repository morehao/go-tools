package utils

import "strings"

// SnakeToPascal 蛇形转大驼峰
func SnakeToPascal(s string) string {
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

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

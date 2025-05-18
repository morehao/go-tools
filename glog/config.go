/*
 * @Author: morehao morehao@qq.com
 * @Date: 2025-04-26 09:55:22
 * @LastEditors: morehao morehao@qq.com
 * @LastEditTime: 2025-04-26 16:50:59
 * @FilePath: /golib/glog/config.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package glog

// RotateUnit 日志切割的时间单位
type RotateUnit string

const (
	RotateUnitDay  RotateUnit = "day"
	RotateUnitHour RotateUnit = "hour"
)

// LogConfig 模块级别的日志配置
type LogConfig struct {
	// Service 服务名
	Service string
	// Module 模块名称，如 "es", "gorm", "redis" 等
	Module string
	// Level 日志级别
	Level Level `json:"level" yaml:"level"`
	// Writer 日志输出类型
	Writer WriterType `json:"writer" yaml:"writer"`
	// RotateInterval 日志切割周期，单位为天
	RotateInterval RotateIntervalType `json:"rotate_interval" yaml:"rotate_interval"`
	// Dir 日志文件目录
	Dir string `json:"dir" yaml:"dir"`
	// ExtraKeys 需要从上下文中提取的额外字段
	ExtraKeys []string `json:"extra_keys" yaml:"extra_keys"`
	// RotateUnit 日志切割的时间单位
	RotateUnit RotateUnit `json:"rotate_unit" yaml:"rotate_unit"`
}

func GetDefaultLogConfig() *LogConfig {
	return &LogConfig{
		Service:    defaultServiceName,
		Module:     defaultModuleName,
		Level:      DebugLevel,
		Writer:     WriterConsole,
		Dir:        defaultLogDir,
		RotateUnit: RotateUnitDay,
	}
}

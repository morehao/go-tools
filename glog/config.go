/*
 * @Author: morehao morehao@qq.com
 * @Date: 2025-04-26 09:55:22
 * @LastEditors: morehao morehao@qq.com
 * @LastEditTime: 2025-04-26 16:50:59
 * @FilePath: /go-tools/glog/config.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package glog

// LoggerConfig 模块级别的日志配置
type LoggerConfig struct {
	// service 服务名，从 ServiceConfig 继承
	service string
	// module 模块名称，如 "es", "gorm", "redis" 等
	module string
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
}

// ServiceConfig 服务级别的日志配置
type ServiceConfig struct {
	// Service 服务名称，如 "myApp"
	Service string `json:"service" yaml:"service"`
	// Modules 模块配置，key 为模块名称
	Modules map[string]*LoggerConfig `json:"modules" yaml:"modules"`
}

func getDefaultLoggerConfig() *LoggerConfig {
	return &LoggerConfig{
		service: "app",
		module:  "default",
		Level:   InfoLevel,
		Writer:  WriterConsole,
		Dir:     "./log",
	}
}

// FieldHookFunc 字段钩子函数类型
type FieldHookFunc func(fields []Field)

// MessageHookFunc 消息钩子函数类型
type MessageHookFunc func(message string) string

// Option 日志选项
type Option interface {
	apply(cfg *optConfig)
}

type optConfig struct {
	callerSkip      int
	fieldHookFunc   FieldHookFunc
	messageHookFunc MessageHookFunc
}

type option func(cfg *optConfig)

func (fn option) apply(cfg *optConfig) {
	fn(cfg)
}

// WithCallerSkip 设置调用者跳过的层数
func WithCallerSkip(skip int) Option {
	return option(func(cfg *optConfig) {
		cfg.callerSkip = skip
	})
}

// WithFieldHookFunc 设置字段钩子函数
func WithFieldHookFunc(fn FieldHookFunc) Option {
	return option(func(cfg *optConfig) {
		cfg.fieldHookFunc = fn
	})
}

// WithMessageHookFunc 设置消息钩子函数
func WithMessageHookFunc(fn MessageHookFunc) Option {
	return option(func(cfg *optConfig) {
		cfg.messageHookFunc = fn
	})
}

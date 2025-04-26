/*
 * @Author: morehao morehao@qq.com
 * @Date: 2025-04-26 09:55:22
 * @LastEditors: morehao morehao@qq.com
 * @LastEditTime: 2025-04-26 16:50:59
 * @FilePath: /go-tools/glog/config.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package glog

// LoggerConfig 日志配置
type LoggerConfig struct {
	Service   string     `yaml:"service"` // 服务名称
	Module    string     `yaml:"module"`  // 模块名称，如 default、gorm、es、redis 等
	Level     Level      `yaml:"level"`
	Type      WriterType `yaml:"type"` // 输出类型：console 或 file
	Dir       string     `yaml:"dir"`  // 日志文件目录
	ExtraKeys []string   `yaml:"extra_keys"`
}

// ServiceConfig 服务配置
type ServiceConfig struct {
	Service string                   `yaml:"service"` // 服务名称，如 myApp
	Modules map[string]*LoggerConfig `yaml:"modules"` // 模块配置
}

func getDefaultLoggerConfig() *LoggerConfig {
	return &LoggerConfig{
		Service: "app",
		Module:  "default",
		Level:   InfoLevel,
		Type:    WriterConsole,
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

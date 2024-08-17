package glog

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type LoggerConfig struct {
	LoggerType LoggerType `yaml:"logger_type"`
	Service    string     `yaml:"service"`
	Level      Level      `yaml:"level"`
	Dir        string     `yaml:"dir"`
	Stdout     bool       `yaml:"stdout"`
	ExtraKeys  []string   `yaml:"extra_keys"`
}

func getDefaultLoggerConfig() *LoggerConfig {
	return &LoggerConfig{
		LoggerType: LoggerTypeZap,
		Service:    "app",
		Level:      InfoLevel,
		Dir:        "./log",
	}
}

type ZapFieldHookFunc func(fields []zapcore.Field)

type MessageHookFunc func(message string) string

type optConfig struct {
	zapOpts          []zap.Option
	zapFieldHookFunc ZapFieldHookFunc
	messageHookFunc  MessageHookFunc
}

type Option interface {
	apply(cfg *optConfig)
}

type option func(cfg *optConfig)

func (fn option) apply(cfg *optConfig) {
	fn(cfg)
}
func WithZapOptions(opts ...zap.Option) Option {
	return option(func(cfg *optConfig) {
		cfg.zapOpts = append(cfg.zapOpts, opts...)
	})
}

func WithZapFieldHookFunc(fn ZapFieldHookFunc) Option {
	return option(func(cfg *optConfig) {
		cfg.zapFieldHookFunc = fn
	})
}

func WithMessageHookFunc(fn MessageHookFunc) Option {
	return option(func(cfg *optConfig) {
		cfg.messageHookFunc = fn
	})
}

package glog

import "go.uber.org/zap"

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

type optConfig struct {
	zapOpts []zap.Option
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

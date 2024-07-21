package glog

import "go.uber.org/zap"

type Option interface {
	apply(cfg *optConfig)
}

type option func(cfg *optConfig)

func (fn option) apply(cfg *optConfig) {
	fn(cfg)
}

type optConfig struct {
	zapOpts []zap.Option
}

func WithZapOptions(opts ...zap.Option) Option {
	return option(func(cfg *optConfig) {
		cfg.zapOpts = append(cfg.zapOpts, opts...)
	})
}

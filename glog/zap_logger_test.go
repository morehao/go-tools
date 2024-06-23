package glog

import (
	"context"
	"testing"
)

func TestZapLogger(t *testing.T) {
	InitZapLogger(&LoggerConfig{
		ServiceName: "myApp",
		Level:       DebugLevel,
		Dir:         "./log",
		Stdout:      true,
	})
	defer Close()
	ctx := context.Background()
	Info(ctx, "hello world")
	Infof(ctx, "hello %s", "world")
	Infow(ctx, "hello world", "key", "value")
	Error(ctx, "hello world")
	Errorf(ctx, "hello %s", "world")
	Errorw(ctx, "hello world", "key", "value")
}

func TestZapExtraKeys(t *testing.T) {
	InitZapLogger(&LoggerConfig{
		ServiceName: "myApp",
		Level:       DebugLevel,
		Dir:         "./log",
		Stdout:      true,
		ExtraKeys:   []string{"key1", "key2"},
	})
	defer Close()
	ctx := context.Background()
	ctx = context.WithValue(ctx, "key1", "value1")
	ctx = context.WithValue(ctx, "key2", "value2")
	Info(ctx, "hello world")
	Infof(ctx, "hello %s", "world")
	Infow(ctx, "hello world", "key", "value")
	Error(ctx, "hello world")
	Errorf(ctx, "hello %s", "world")
	Errorw(ctx, "hello world", "key", "value")
}

package clog

import "log/slog"

type Options struct {
	Level slog.Level
	Dir   string
}
type Option func(*Options)

func defaultOptions() *Options {
	return &Options{
		Level: slog.LevelInfo,
	}
}

func NewOptions(opt ...Option) *Options {
	d := defaultOptions()
	for _, o := range opt {
		o(d)
	}
	return d
}

func WithLevel(level slog.Level) Option {
	return func(o *Options) {
		o.Level = level
	}
}

func WithDir(dir string) Option {
	return func(o *Options) {
		o.Dir = dir
	}
}

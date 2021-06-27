package middleware

import "log"

// Options represents configs that can be used for StatsD client
type Options struct {
	SampleRate float32
}

// Option represent a func to build up StatsD client options
type Option func(*Options)

// WithSample constructs the SampleRate option
func WithSample(rate float32) Option {
	if rate < 0 || rate > 1 {
		log.Println("warn: rate should be between 0.0 and 1.0")
	}

	return func(opt *Options) {
		opt.SampleRate = rate
	}
}

func NewOptions() *Options {
	return &Options{
		SampleRate: 1.0,
	}
}

func NewOptionsWith(opts ...Option) *Options {
	defaultOpts := NewOptions()

	for _, opt := range opts {
		opt(defaultOpts)
	}

	return defaultOpts
}

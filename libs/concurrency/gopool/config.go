package gopool

import "sync/atomic"

type Config struct {
	ScaleThreshold int64
	Cap            int64
}

func (c *Config) setCap(cap int64) {
	atomic.StoreInt64(&c.Cap, cap)
}

func (c *Config) getCap() int64 {
	return atomic.LoadInt64(&c.Cap)
}

type ConfigOpt func(config *Config)

func WithConfigScaleThreshold(scaleThreshold int64) ConfigOpt {
	return func(config *Config) {
		config.ScaleThreshold = scaleThreshold
	}
}
func WithConfigCap(cap int64) ConfigOpt {
	return func(config *Config) {
		config.Cap = cap
	}
}

func NewConfig(opts ...ConfigOpt) *Config {
	config := &Config{
		ScaleThreshold: 1,
		Cap:            100,
	}
	for _, opt := range opts {
		opt(config)
	}
	return config
}

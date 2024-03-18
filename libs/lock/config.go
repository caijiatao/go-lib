package lock

type Config struct {
	TTL int64
}

func NewConfig() *Config {
	return &Config{TTL: 5}
}

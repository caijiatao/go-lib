package config

import (
	"sync"
)

var (
	config *Config
	once   sync.Once
)

func Conf() *Config {
	return config
}

type Config struct {
	ETCDEndpoints []string
	Env           *Env
	RPCServer     *RPCServer
}

func init() {
	once.Do(func() {
		config = &Config{
			ETCDEndpoints: []string{"localhost:2379"},
			Env: &Env{
				Host: "localhost",
				Port: "13137",
			},
			RPCServer: &RPCServer{
				Network: "tcp",
			},
		}
	})
}

type Env struct {
	Host string
	Port string
}

func (e *Env) GetTarget() string {
	return e.Host + ":" + e.Port
}

type RPCServer struct {
	Network string
}

package config

import "github.com/zeromicro/go-zero/zrpc"

var (
	c = &Config{}
)

func Conf() *Config {
	return c
}

type Config struct {
	zrpc.RpcServerConf

	Chat zrpc.RpcClientConf
}

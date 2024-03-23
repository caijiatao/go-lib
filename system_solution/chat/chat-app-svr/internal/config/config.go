package config

import (
	"fmt"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
	"golib/libs/net_helper"
)

type Config struct {
	rest.RestConf

	User zrpc.RpcClientConf
	//Chat zrpc.RpcClientConf
}

var (
	c = &Config{}
)

func Conf() *Config {
	return c
}

func (c *Config) GetServerID() string {
	return net_helper.GetFigureOutListenOn(fmt.Sprintf("%s:%s", c.RestConf.Host, c.RestConf.Port))
}

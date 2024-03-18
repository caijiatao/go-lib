package etcd_helper

import (
	"github.com/spf13/viper"
	clientv3 "go.etcd.io/etcd/client/v3"
	"time"
)

const (
	etcdConfigKey = "etcd_config"
)

type Config struct {
	// 标识ETCD的配置，如果有多个通过该名字进行获取
	ClientName string
	clientv3.Config
}

func readConfigs() []Config {
	etcdConfigMap := viper.GetStringMap(etcdConfigKey)
	configs := make([]Config, 0)
	for k, _ := range etcdConfigMap {
		etcdConfig := viper.Sub(etcdConfigKey).Sub(k)
		config := Config{
			ClientName: k,
			Config: clientv3.Config{
				Endpoints:   etcdConfig.GetStringSlice("endpoints"),
				DialTimeout: time.Second * time.Duration(etcdConfig.GetInt("dial-timeout")),
			},
		}
		configs = append(configs, config)
	}
	return configs
}

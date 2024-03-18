package redis

import (
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

const (
	redisConfigKey     = "redis_config"
	redisClusterClient = "cluster_client"
)

type Config struct {
	ClusterOptions redis.ClusterOptions
}

func redisConfigs() *Config {
	config := &Config{
		ClusterOptions: redis.ClusterOptions{
			Addrs:       viper.GetStringSlice(fmt.Sprintf("%s.%s.%s", redisConfigKey, redisClusterClient, "addrs")),
			Password:    viper.GetString(fmt.Sprintf("%s.%s.%s", redisConfigKey, redisClusterClient, "password")),
			ClientName:  viper.GetString(fmt.Sprintf("%s.%s.%s", redisConfigKey, redisClusterClient, "client_name")),
			DialTimeout: viper.GetDuration(fmt.Sprintf("%s.%s.%s", redisConfigKey, redisClusterClient, "dial_timeout")),
		},
	}
	return config
}

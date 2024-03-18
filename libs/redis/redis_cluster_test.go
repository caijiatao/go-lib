package redis

import (
	"context"
	"github.com/redis/go-redis/v9"
	"testing"
)

func TestRedis(t *testing.T) {
}

func TestRedisHset(t *testing.T) {
	client = &RedisClient{
		clusterClient: redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:       []string{"192.168.42.40:6370", "192.168.42.40:6371", "192.168.42.40:6372", "192.168.42.40:6373", "192.168.42.40:6374", "192.168.42.40:6375"},
			Password:    "",
			ClientName:  "",
			DialTimeout: 0,
		}),
	}
	err := ClusterClient().HSet(context.Background(), "fq", "name", "0").Err()
	if err != nil {
		t.Log(err)
	}
	err = ClusterClient().HSet(context.Background(), "fq", "age", "0").Err()
	if err != nil {
		t.Log(err)
	}
	res := ClusterClient().HGetAll(context.Background(), "fq").Val()
	for key, val := range res {
		t.Log(key, ":", val)
	}
	t.Log(res)
}

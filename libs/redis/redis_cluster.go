package redis

import (
	"context"
	"errors"
	"github.com/avast/retry-go/v4"
	"github.com/hashicorp/go-uuid"
	"github.com/redis/go-redis/v9"
	"net"
	"sync"
	"time"
)

var client *RedisClient
var clientOnce sync.Once
var onceErr []error

type RedisClient struct {
	clusterClient *redis.ClusterClient
}

func Init() error {
	clientOnce.Do(func() {
		redisConfig := redisConfigs()
		for _, addr := range redisConfig.ClusterOptions.Addrs {
			_, err := net.Dial("tcp", addr)
			onceErr = append(onceErr, err)
		}
		client = &RedisClient{
			clusterClient: redis.NewClusterClient(&redis.ClusterOptions{
				Addrs:       redisConfig.ClusterOptions.Addrs,
				Password:    redisConfig.ClusterOptions.Password,
				ClientName:  redisConfig.ClusterOptions.ClientName,
				DialTimeout: redisConfig.ClusterOptions.DialTimeout,
			}),
		}
	})
	for _, err := range onceErr {
		if err != nil {
			return err
		}
	}
	return nil
}

func ClusterClient() *redis.ClusterClient {
	return client.clusterClient
}

type MutexLock struct {
	ctx        context.Context
	key        string
	expireTime time.Duration
	retryDelay time.Duration
	retryCount uint
}

func NewRedisLock(ctx context.Context, key string, expireTime, retryDelay time.Duration, retryCount uint) *MutexLock {
	return &MutexLock{
		ctx:        ctx,
		key:        key,
		expireTime: expireTime,
		retryDelay: retryDelay,
		retryCount: retryCount,
	}
}

func NewDefaultRedisLock(ctx context.Context) *MutexLock {
	uuidKey, _ := uuid.GenerateUUID()
	return &MutexLock{
		ctx:        ctx,
		key:        uuidKey,
		expireTime: 60 * time.Second,
		retryDelay: 5 * time.Second,
		retryCount: 5,
	}
}

func (self MutexLock) Lock() error {
	err := retry.Do(
		func() error {
			ok, err := ClusterClient().SetNX(self.ctx, self.key, 1, self.expireTime).Result()
			if err != nil {
				return err
			}
			if !ok {
				return errors.New("can't get redis lock")
			}
			return nil
		},
		retry.Delay(self.retryDelay),
		retry.Attempts(self.retryCount),
	)
	if err != nil {
		return err
	}
	return nil
}

func (self MutexLock) Unlock() error {
	return ClusterClient().Del(self.ctx, self.key).Err()
}

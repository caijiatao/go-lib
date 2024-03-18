package etcd_helper

import (
	"context"
	clientv3 "go.etcd.io/etcd/client/v3"
	"golib/libs/logger"
	"time"
)

type etcdContextKeyType string

const (
	etcdClientContextKey etcdContextKeyType = "etcd-context-key"
)

const defaultClientName = "default"

type EtcdClient struct {
	*clientv3.Client
	Config
}

func newEtcdClient(client *clientv3.Client, config Config) *EtcdClient {
	return &EtcdClient{
		Client: client,
		Config: config,
	}
}

func BindContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, etcdClientContextKey, globalManager)
}

func Context(ctx context.Context) *EtcdClient {
	m, ok := ctx.Value(etcdClientContextKey).(*etcdManager)
	if !ok {
		return nil
	}
	return m.get(defaultClientName)
}

func Init() error {
	configs := readConfigs()
	for _, config := range configs {
		err := InitETCDClient(config)
		if err != nil {
			return err
		}
	}
	return nil
}

func InitETCDClient(config Config) (err error) {
	globalManagerInitOnce.Do(func() {
		globalManager = &etcdManager{}
	})
	cli, err := clientv3.New(config.Config)
	if err != nil {
		return err
	}
	client := newEtcdClient(cli, config)

	err = client.Ping(context.Background())
	if err != nil {
		return err
	}

	if len(config.ClientName) == 0 {
		config.ClientName = defaultClientName
	}
	err = globalManager.add(config.ClientName, client)
	if err != nil {
		return err
	}
	return nil
}

func Close() {

}

func Put(ctx context.Context, key, val string, opts ...clientv3.OpOption) (*clientv3.PutResponse, error) {
	resp, err := Context(ctx).Put(ctx, key, val, opts...)
	if err != nil {
		logger.Error(err)
		return resp, err
	}
	return resp, nil
}

func Delete(ctx context.Context, key string, opts ...clientv3.OpOption) (*clientv3.DeleteResponse, error) {
	resp, err := Context(ctx).Delete(ctx, key, opts...)
	if err != nil {
		logger.Error(err)
		return resp, err
	}
	return resp, nil
}

func Get(ctx context.Context, key string, opts ...clientv3.OpOption) (*clientv3.GetResponse, error) {
	resp, err := Context(ctx).Get(ctx, key, opts...)
	if err != nil {
		logger.Error(err)
		return resp, err
	}
	return resp, nil
}

func Watch(ctx context.Context, key string, opts ...clientv3.OpOption) clientv3.WatchChan {
	return Context(ctx).Watch(ctx, key, opts...)
}

func (client *EtcdClient) Ping(ctx context.Context) error {
	ctx, _ = context.WithTimeout(ctx, time.Second)
	_, err := client.Get(ctx, "ping")
	if err != nil {
		return err
	}
	return nil
}

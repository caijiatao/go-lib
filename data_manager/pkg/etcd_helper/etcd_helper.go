package etcd_helper

import (
	"airec_server/pkg/logger"
	"context"
	clientv3 "go.etcd.io/etcd/client/v3"
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

func InitETCDClient(config Config) (err error) {
	globalManagerInitOnce.Do(func() {
		globalManager = &etcdManager{}
	})
	cli, err := clientv3.New(config.Config)
	if err != nil {
		return err
	}
	client := newEtcdClient(cli, config)

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

func Watch(ctx context.Context, key string, opts ...clientv3.OpOption) clientv3.WatchChan {
	return Context(ctx).Watch(ctx, key, opts...)
}
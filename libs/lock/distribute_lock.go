package lock

import (
	"context"
	clientv3 "go.etcd.io/etcd/client/v3"
	"golib/libs/etcd_helper"
)

type DistributedLock interface {
	Lock(ctx context.Context, key string) error
	UnLock(ctx context.Context) error
}

func NewLock() (DistributedLock, error) {
	etcdClient := etcd_helper.GetDefaultClient()
	lease := clientv3.NewLease(etcdClient.Client)

	return &ETCDLock{
		client: etcdClient,
		lease:  lease,
		cancel: nil,
		config: NewConfig(),
	}, nil
}

package lock

import (
	"context"
	clientv3 "go.etcd.io/etcd/client/v3"
	"golib/etcd_helper"
)

type ETCDLock struct {
	client *etcd_helper.EtcdClient
	lease  clientv3.Lease
	key    clientv3.LeaseID
	cancel context.CancelFunc
	config *Config
}

func (l *ETCDLock) Lock(ctx context.Context, key string) (err error) {
	ctx, l.cancel = context.WithCancel(ctx)
	defer func() {
		if err != nil {
			l.cancel()
		}
	}()

	grantResp, err := l.lease.Grant(ctx, l.config.TTL)
	if err != nil {
		return err
	}
	_, err = l.lease.KeepAlive(ctx, grantResp.ID) // 续租
	if err != nil {
		return err
	}

	txn := l.client.Client.Txn(ctx)
	_, err = txn.If(clientv3.Compare(clientv3.CreateRevision(key), "=", 0)).
		Then(clientv3.OpPut(key, "", clientv3.WithLease(grantResp.ID))).
		Commit()
	if err != nil {
		return err
	}

	l.key = grantResp.ID

	return nil
}

func (l *ETCDLock) UnLock(ctx context.Context) error {
	defer l.cancel()
	_, err := l.lease.Revoke(ctx, l.key)
	if err != nil {
		return err
	}
	return nil
}

package etcd_helper

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	clientv3 "go.etcd.io/etcd/client/v3"
	"testing"
)

func TestPing(t *testing.T) {
	InitTestSuite()

	ctx := BindContext(context.Background())
	err := Context(ctx).Ping(ctx)
	assert.Nil(t, err)
}

func TestWatchTimeoutKey(t *testing.T) {
	InitTestSuite()
	ctx := BindContext(context.Background())
	// 设置key的过期时间，例如10秒后过期
	leaseGrantResp, err := Context(ctx).Grant(ctx, 1)
	assert.Nil(t, err)
	assert.NotEqual(t, int64(0), leaseGrantResp.ID)

	_, err = Put(ctx, "/test/watch/timeout/key", "value", clientv3.WithLease(leaseGrantResp.ID))
	assert.Nil(t, err)
	wch := Watch(ctx, "/test/watch/timeout/key")
	for w := range wch {
		for _, event := range w.Events {
			if event.Type == clientv3.EventTypeDelete {
				fmt.Println("delete key")
			}
			fmt.Println("watch key")
		}
	}
}

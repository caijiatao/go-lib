package lock

import (
	"context"
	"github.com/stretchr/testify/assert"
	"golib/etcd_helper"
	"golib/logger"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	etcd_helper.InitTestClient()
	m.Run()
}

func TestLock(t *testing.T) {
	l, err := NewLock()
	assert.Nil(t, err)
	ctx := context.Background()
	lockKey := "test"
	go func() {
		err := l.Lock(ctx, lockKey)
		if err != nil {
			logger.Error("lock1 fail")
			return
		}
		time.Sleep(1 * time.Second)
		err = l.UnLock(ctx)
		assert.Nil(t, err)
	}()

	go func() {
		err := l.Lock(ctx, lockKey)
		if err != nil {
			logger.Error("lock2 fail")
			return
		}
		time.Sleep(1 * time.Second)
		err = l.UnLock(ctx)
		assert.Nil(t, err)
	}()

	time.Sleep(10 * time.Second)
}

package etcd_helper

import (
	clientv3 "go.etcd.io/etcd/client/v3"
	"time"
)

func InitTestClient() {
	_, err := NewClient(Config{
		Config: clientv3.Config{
			Endpoints:   []string{"127.0.0.1:2379"},
			DialTimeout: 10 * time.Second,
		},
	})
	if err != nil {
		panic(err)
	}

	return
}

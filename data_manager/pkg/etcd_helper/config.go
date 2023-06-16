package etcd_helper

import clientv3 "go.etcd.io/etcd/client/v3"

type Config struct {
	// 标识ETCD的配置，如果有多个通过该名字进行获取
	ClientName string
	clientv3.Config
}

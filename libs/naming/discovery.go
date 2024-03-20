package naming

import (
	"context"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
	"golib/libs/etcd_helper"
	"golib/libs/logger"
	"log"
	"sync"
	"time"
)

// ServiceDiscovery 服务发现
type ServiceDiscovery struct {
	serverList map[string]string //服务列表
	lock       sync.Mutex

	changeEvent chan struct{}
}

// NewServiceDiscovery  新建发现服务
func NewServiceDiscovery(endpoints []string, prefix string) *ServiceDiscovery {
	err := etcd_helper.InitETCDClient(etcd_helper.Config{
		ClientName: namingDiscoveryETCDClientKey,
		Config: clientv3.Config{
			Endpoints:   endpoints,
			DialTimeout: 5 * time.Second,
		},
	})
	if err != nil {
		return nil
	}
	svrDisc := &ServiceDiscovery{
		serverList:  make(map[string]string),
		changeEvent: make(chan struct{}, 8),
	}

	err = svrDisc.runWatchService(prefix)
	if err != nil {
		return nil
	}

	return svrDisc
}

// runWatchService 初始化服务列表和监视
func (s *ServiceDiscovery) runWatchService(prefix string) error {
	//根据前缀获取现有的key
	resp, err := etcd_helper.GetClientByName(namingDiscoveryETCDClientKey).Get(context.Background(), prefix, clientv3.WithPrefix())
	if err != nil {
		return err
	}

	for _, ev := range resp.Kvs {
		s.SetServiceList(string(ev.Key), string(ev.Value))
	}

	//监视前缀，修改变更的server
	go s.watcher(prefix)
	return nil
}

// watcher 监听前缀
func (s *ServiceDiscovery) watcher(prefix string) {
	rch := etcd_helper.GetClientByName(namingDiscoveryETCDClientKey).Watch(context.Background(), prefix, clientv3.WithPrefix())
	logger.Infof("watching prefix:%s now...", prefix)
	for wresp := range rch {
		for _, ev := range wresp.Events {
			switch ev.Type {
			case mvccpb.PUT: //修改或者新增
				s.SetServiceList(string(ev.Kv.Key), string(ev.Kv.Value))
			case mvccpb.DELETE: //删除
				s.DelServiceList(string(ev.Kv.Key))
			}
			// 通知变更
			s.changeEvent <- struct{}{}
		}
	}
}

func (s *ServiceDiscovery) Watch() <-chan struct{} {
	return s.changeEvent
}

// SetServiceList 新增服务地址
func (s *ServiceDiscovery) SetServiceList(key, val string) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.serverList[key] = val
	log.Println("put key :", key, "val:", val)
}

// DelServiceList 删除服务地址
func (s *ServiceDiscovery) DelServiceList(key string) {
	s.lock.Lock()
	defer s.lock.Unlock()
	delete(s.serverList, key)
	log.Println("del key:", key)
}

func (s *ServiceDiscovery) GetServiceList() map[string]string {
	return s.serverList
}

// GetServices 获取服务地址
func (s *ServiceDiscovery) GetServices() []string {
	s.lock.Lock()
	defer s.lock.Unlock()
	addrs := make([]string, 0)

	for _, v := range s.serverList {
		addrs = append(addrs, v)
	}
	return addrs
}

func (s *ServiceDiscovery) Close() error {
	close(s.changeEvent)
	return nil
}

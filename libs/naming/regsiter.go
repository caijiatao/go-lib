package naming

import (
	"context"
	clientv3 "go.etcd.io/etcd/client/v3"
	"golib/libs/etcd_helper"
	"golib/libs/logger"
	"log"
	"time"
)

// ServiceRegister 创建租约注册服务
type ServiceRegister struct {
	leaseID clientv3.LeaseID //租约ID
	//租约keepalieve相应chan
	keepAliveChan <-chan *clientv3.LeaseKeepAliveResponse
	key           string //key
	val           string //value
}

// NewServiceRegister 新建注册服务
func NewServiceRegister(endpoints []string, key, val string, lease int64) (*ServiceRegister, error) {
	err := etcd_helper.InitETCDClient(etcd_helper.Config{
		ClientName: namingRegisterETCDClientKey,
		Config: clientv3.Config{
			Endpoints:   endpoints,
			DialTimeout: 5 * time.Second,
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	ser := &ServiceRegister{
		key: key,
		val: val,
	}

	//申请租约设置时间keepalive
	if err := ser.putKeyWithLease(lease); err != nil {
		return nil, err
	}

	return ser, nil
}

// 设置租约
func (s *ServiceRegister) putKeyWithLease(lease int64) error {
	//设置租约时间
	resp, err := etcd_helper.GetClientByName(namingRegisterETCDClientKey).Grant(context.Background(), lease)
	if err != nil {
		return err
	}
	//注册服务并绑定租约
	_, err = etcd_helper.GetClientByName(namingRegisterETCDClientKey).Put(context.Background(), s.key, s.val, clientv3.WithLease(resp.ID))
	if err != nil {
		return err
	}
	//设置续租 定期发送需求请求
	leaseRespChan, err := etcd_helper.GetClientByName(namingRegisterETCDClientKey).KeepAlive(context.Background(), resp.ID)

	if err != nil {
		return err
	}
	s.leaseID = resp.ID
	logger.Infof("%d", s.leaseID)
	s.keepAliveChan = leaseRespChan
	logger.Infof("Put key:%s  val:%s  success!", s.key, s.val)
	return nil
}

// ListenLeaseRespChan 监听 续租情况
func (s *ServiceRegister) ListenLeaseRespChan() {
	for leaseKeepResp := range s.keepAliveChan {
		logger.Infof("续约成功", leaseKeepResp)
	}
	logger.Infof("关闭续租")
}

// Close 注销服务
func (s *ServiceRegister) Close() error {
	//撤销租约
	if _, err := etcd_helper.GetClientByName(namingRegisterETCDClientKey).Revoke(context.Background(), s.leaseID); err != nil {
		return err
	}
	logger.Infof("撤销租约")
	return etcd_helper.GetClientByName(namingRegisterETCDClientKey).Close()
}

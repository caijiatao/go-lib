package election_master

import (
	"context"
	"errors"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
	"golib/goasync"
	"golib/net_helper"
	"time"
)

func Campaign(ctx context.Context, etcdClient *clientv3.Client) (notify chan struct{}) {
	notify = make(chan struct{}, 100)
	go func() {
		defer func() {
			r := recover()
			_ = goasync.PanicErrHandler(r)
		}()
		for {
			select {
			case <-ctx.Done():
				return
			default:
			}
			_ = campaign(ctx, etcdClient, notify)
		}
	}()
	return
}

func campaign(ctx context.Context, etcdClient *clientv3.Client, notify chan struct{}) (err error) {
	session, err := concurrency.NewSession(etcdClient)
	if err != nil {
		return
	}
	defer func() {
		if err != nil {
			fmt.Println(err)
		}
		if session != nil {
			_ = session.Close()
		}
	}()
	election := concurrency.NewElection(session, "/etcd-campaign-demo")
	ip, err := net_helper.GetLocalIP()
	if err != nil {
		return
	}
	// 成为leader节点会运行出来，没选举成功则会阻塞在这里
	err = election.Campaign(ctx, ip)
	if err != nil {
		return
	}

	// 通知程序已经成为leader，可以避免轮训leader状态
	for {
		select {
		case notify <- struct{}{}: // 发送通知说明已经是leader了
		case <-session.Done(): // session断开要重新进行选举
			err = errors.New("session is done")
			return
		case <-ctx.Done():
			// 给context设定了超时时间,避免放弃leader的时间超时
			resignCtx, _ := context.WithTimeout(context.Background(), time.Second)
			_ = election.Resign(resignCtx)
			err = errors.New("ctx is done")
			return
		}
	}
}

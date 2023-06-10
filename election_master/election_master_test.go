package election_master

import (
	"context"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"testing"
	"time"
)

func TestCampaign(t *testing.T) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		// handle error!
	}
	defer cli.Close()

	var (
		ctx = context.Background()
	)
	notify := Campaign(ctx, cli)

	for {
		select {
		case <-notify:
			fmt.Println("master")
		case <-ctx.Done():

		}
	}
}

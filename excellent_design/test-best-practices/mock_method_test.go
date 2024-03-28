package test_best_practices

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_raceGraceTerminateRSList_flushList(t *testing.T) {
	manager := newGracefulTerminationManager()
	go func() {
		for i := 0; i < 100; i++ {
			manager.rsList.add(&item{
				VirtualServer: "virtualServer",
				RealServer:    fmt.Sprint(i),
			})
		}
	}()

	// 等待加入到一定元素再继续往下执行
	for manager.rsList.len() < 20 {
	}

	success := manager.rsList.flushList(func(rsToDelete *item) (bool, error) {
		return true, nil
	})

	assert.True(t, success)
}

package concurrency

import (
	"fmt"
	"sync"
	"testing"
)

func TestUserNameCacheRace(t *testing.T) {
	// 使用 WaitGroup 来等待所有 goroutine 完成
	var wg sync.WaitGroup

	// 模拟并发的读取操作
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			_ = GetUserNameByUserId(fmt.Sprintf("user%d", id%10))
		}(i)
	}

	// 模拟并发的写操作
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			RefreshUserNameCache()
		}()
	}

	// 等待所有操作完成
	wg.Wait()
}

package concurrency

import (
	"fmt"
	"sync"
)

var (
	userId2UserNameCache = make(map[string]string)
	userNameCacheLock    = sync.RWMutex{}
)

func GetUserNameByUserId(userId string) string {
	userNameCacheLock.RLock()
	name := userId2UserNameCache[userId]
	userNameCacheLock.RUnlock()
	return name
}

func RefreshUserNameCache() {
	userId2UserName := make(map[string]string)

	for i := 0; i < 10; i++ {
		userId2UserName[fmt.Sprintf("user%d", i)] = fmt.Sprintf("user%d", i)
	}

	userNameCacheLock.Lock()
	userId2UserNameCache = userId2UserName
	userNameCacheLock.Unlock()
}

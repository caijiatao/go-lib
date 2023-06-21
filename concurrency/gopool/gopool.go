package gopool

import "context"

func init() {
	initManager()
	initTask()
	initWorker()
}

func RegisterPool(pool Pool) error {
	return globalPoolManager.register(pool)
}

func GetPoolByName(name string) Pool {
	return globalPoolManager.getPoolByName(name)
}

func Go(fc func()) {
	globalPoolManager.getPoolByName(defaultPoolName).Go(fc)
}

func GoWithContext(ctx context.Context, fc func()) {
	globalPoolManager.getPoolByName(defaultPoolName).GoWithContext(ctx, fc)
}

func WorkCount() int64 {
	return globalPoolManager.getPoolByName(defaultPoolName).WorkerCount()
}

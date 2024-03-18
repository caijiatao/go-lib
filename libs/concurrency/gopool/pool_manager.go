package gopool

import (
	"errors"
	"sync"
)

var (
	globalPoolManager *poolManager
	defaultPoolName   = "gopool.DefaultPool"
	poolNameDuplicate = errors.New("pool name duplicate")
)

func initManager() {
	globalPoolManager = &poolManager{}
	NewPool(defaultPoolName, NewConfig())
}

type poolManager struct {
	poolMap sync.Map
}

func (p *poolManager) register(pool Pool) error {
	_, ok := globalPoolManager.poolMap.Load(pool.Name())
	if ok {
		return poolNameDuplicate
	}
	globalPoolManager.poolMap.Store(pool.Name(), pool)
	return nil
}

func (p *poolManager) getPoolByName(name string) Pool {
	if value, ok := p.poolMap.Load(name); ok {
		return value.(Pool)
	}
	return nil
}

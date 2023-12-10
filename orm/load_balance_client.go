package orm

import (
	"database/sql"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

var (
	dbEndpoints         = make([]*sql.DB, 0)
	index               atomic.Int64
	initDBEndpointsOnce sync.Once
)

func GetDB() (*sql.DB, error) {
	initDBEndpointsOnce.Do(func() {
		// 初始化所有连接
		initDBEndpoints()
		go updateEndpoints()
	})
	// 简单的轮询，这里可以实现更复杂的负载均衡算法
	i := int(index.Add(1)) % len(dbEndpoints)
	return dbEndpoints[i], nil
}

func updateEndpoints() {
	for {
		select {
		case <-time.After(time.Second):
			// 循环查看连接是否正常，也可以通过监听配置中心热更新数据库连接
		}
	}
}

func initDBEndpoints() {
	endpoints := []string{"endpoint1", "endpoint2", "endpoint3"}
	for _, endpoint := range endpoints {
		dataSourceName := fmt.Sprintf("%s:%s@tcp(%s)/%s", "username", "password", endpoint, "dbname")

		poolDB, err := sql.Open("mysql", dataSourceName)
		if err != nil {
			// 异常处理
		}
		poolDB.SetMaxOpenConns(10)
		poolDB.SetMaxIdleConns(5)
		dbEndpoints = append(dbEndpoints, poolDB)
	}
}

package orm

import (
	"encoding/csv"
	"fmt"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"gorm.io/sharding"
	"os"
	"testing"
)

func TestGetShardingTableIndex(t *testing.T) {
	// 打开CSV文件
	file, err := os.Open("E:\\cas-built-in\\20240328\\behavior.csv")
	if err != nil {
		return
	}
	defer file.Close()

	// 创建CSV读取器
	reader := csv.NewReader(file)

	// 读取所有记录
	records, err := reader.ReadAll()
	if err != nil {
		return
	}

	userIds := make(map[int64]int)
	isRecord := make(map[string]bool)
	for _, record := range records {
		userId := record[5]
		if isRecord[userId] {
			continue
		}
		isRecord[userId] = true
		userIds[GetShardingTableIndex(userId, 128)]++
	}

	fmt.Println(userIds)
	fmt.Println(len(userIds))
}

func TestUserSharding(t *testing.T) {
	// 打开CSV文件
	file, err := os.Open("E:\\cas-built-in\\20240328\\user.csv")
	if err != nil {
		return
	}
	defer file.Close()

	// 创建CSV读取器
	reader := csv.NewReader(file)

	// 读取所有记录
	records, err := reader.ReadAll()
	if err != nil {
		return
	}

	userIds := make(map[int64]int)
	isRecord := make(map[string]bool)
	for _, record := range records {
		userId := record[0]
		if isRecord[userId] {
			continue
		}
		isRecord[userId] = true
		userIds[GetShardingTableIndex(userId, 128)]++
	}

	fmt.Println(userIds)
	fmt.Println(len(userIds))
}

func TestGetShardingTableIndex1(t *testing.T) {
	tableIndex := GetShardingTableIndex("1", 128)
	fmt.Println(tableIndex)
}

func TestRegisterSharding(t *testing.T) {
	var err error
	var db gorm.DB
	err = db.Use(sharding.Register(sharding.Config{
		ShardingKey:    "user_id",
		NumberOfShards: 128,
		// 自定义传入
		ShardingAlgorithm: RecommendResultShardingAlgorithm,
		// 不需要用到推荐结果主键，使用雪花id即可
		PrimaryKeyGenerator: sharding.PKSnowflake,
	}, "recommend_result"))
	assert.Nil(t, err)
	err = db.Exec("select * from recommend_result where user_id = ?", "129").Error
	assert.NotNil(t, err)
}

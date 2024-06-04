package orm

import (
	"encoding/csv"
	"fmt"
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
		userIds[GetShardingTableIndex(userId, 100)]++
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
		userIds[GetShardingTableIndex(userId, 100)]++
	}

	fmt.Println(userIds)
	fmt.Println(len(userIds))
}

package orm

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"math/rand"
	"testing"
)

func TestUpdateCreateTime(t *testing.T) {
	clientName := fmt.Sprintf("test%d", rand.Intn(100000))
	err := NewOrmClient(&Config{
		DBClientName: clientName,
		Config:       &gorm.Config{},
		SourceConfig: &SourceDBConfig{},
		Dial:         mysql.Open("root:123456@tcp(192.168.15.54:3306)/test_zky_1000w"),
	})
	if err != nil {
		t.Log(err)
		return
	}

	client := GetClientByClientName(clientName)
	if client == nil {
		t.Log("client is nil")
		return
	}

	itemIds := make([]string, 0)
	readItemIds := func(data []map[string]interface{}) {
		for _, item := range data {
			itemIds = append(itemIds, item["item_id"].(string))
		}
	}

	_, err = client.QueryByCursor("zky_aiticle_1000w", 100000, []string{"item_id"}, "where new_create_time is NULL", nil, readItemIds)
	if err != nil {
		t.Log(err)
		return
	}

	fmt.Println(len(itemIds))

	batchSize := 10000
	for i := 0; i < len(itemIds); i += batchSize {
		end := i + batchSize
		if end > len(itemIds) {
			end = len(itemIds)
		}
		subset := itemIds[i:end]
		client.Exec("UPDATE test_zky_1000w.zky_aiticle_1000w SET new_create_time = UNIX_TIMESTAMP(create_time) where item_id in (?)", subset)
	}

}

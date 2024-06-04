package orm

import (
	"crypto/sha256"
	"encoding/hex"
	"strconv"
)

func GetShardingTableIndex(userID string, totalTables int) int64 {
	// 使用SHA-256对userID进行哈希计算
	hash := sha256.New()
	hash.Write([]byte(userID))
	hashBytes := hash.Sum(nil)
	// 将哈希值转换为十六进制字符串
	hashHex := hex.EncodeToString(hashBytes)
	// 将十六进制字符串转换为整数
	hashInt, _ := strconv.ParseInt(hashHex[:8], 16, 64) // 只取前8位来减少长度
	// 取模运算，确定表名
	tableIndex := hashInt % int64(totalTables)
	return tableIndex
}

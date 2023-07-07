package postgres_util

import "strings"

func GetColAlias(tableName string, colName string) string {
	return strings.Join([]string{strings.ReplaceAll(tableName, ".", "_"), colName}, "_")
}

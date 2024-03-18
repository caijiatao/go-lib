package postgres_util

import "strings"

func GetColAlias(tableName string, colName string) string {
	return strings.Join([]string{strings.ReplaceAll(tableName, ".", "_"), colName}, "_")
}

func GetSchemaAndTableName(fullTableName string) (schema string, tableName string) {
	fullTableNames := strings.Split(fullTableName, ".")
	if len(fullTableNames) == 2 {
		return fullTableNames[0], fullTableNames[1]
	}
	if len(fullTableNames) == 1 {
		return "", fullTableNames[0]
	}
	return "", ""
}

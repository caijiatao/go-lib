package model

// StandardTable
// @Description: 标准数据模型表
type StandardTable struct {
	TableId    uint64
	Field      string
	TableName  string
	TableDesc  string
	IsNecess   int64 // 是否必须
	CreateTime uint64
}

// StandardTableColumn
// @Description: 标准模型表关联的列meta信息
type StandardTableColumn struct {
	Id         uint64
	TableId    uint64 // 关联的标准模型表
	ColumnName string
	ColumnType string
	ColumnDesc string
	IsNecess   int64
	CreateTime uint64
}

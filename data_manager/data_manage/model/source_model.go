package model

// Source
// @Description: 源的描述信息
type Source struct {
	SourceId       uint64
	SourceType     uint64 // 数据源类型，有文件、DB等
	SourceConfig   string // 数据源配置
	Name           string
	Desc           string
	IsDel          int8
	LastUpdateTime uint64
	CreateTime     uint64
}

// SourceTable
// @Description: 表数据集信息
type SourceTable struct {
	SourceTableId    uint64
	SourceId         uint64
	TableName        string // 数据表
	SyncCompleteTime uint64 // 同步完成时间
	IsDel            int8   // 指表是否被用户删除？
	LastUpdateTime   uint64
	CreateTime       uint64
}

// SourceColumn
// @Description: 源数据列信息
type SourceColumn struct {
	SourceColumnId uint64
	SourceTableId  uint
	ColumnName     string
	ColumnType     string
	ColumnDesc     string
	LastUpdateTime uint64
	CreateTime     uint64
}

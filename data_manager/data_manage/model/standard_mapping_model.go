package model

type StandardMapping struct {
	Id uint64
	// 标准模型的信息
	TableId  uint64
	ColumnId uint64

	// 数据源模型的信息
	SourceId       uint64
	SourceTableId  uint64
	SourceColumnId uint64

	ConvertScript  string
	IsDel          int8
	LastUpdateTime uint64
	CreateTime     uint64
}

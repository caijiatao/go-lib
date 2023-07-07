package orm

// ITableMeta
// @Description: 表的元信息
type ITableMeta interface {
}

// IColumnMeta
// @Description: 列的元信息
type IColumnMeta interface {
}

type gormTableMeta struct {
	name    string
	columns []IColumnMeta
}

type gormColumnMeta struct {
	Field string `gorm:"column:Field"`
	Type  string `gorm:"column:Type"`
}

package orm

import "gorm.io/gorm"

type driverMetaConstructor func(db *gorm.DB) DriverMetaOperator

var (
	driverMetaConstructorMap = map[string]driverMetaConstructor{
		"postgres": newPostgreDriverMetaOperator,
		"mysql":    newMysqlDriverMetaOperator,
	}
)

// DriverMetaOperator
// @Description: 驱动元信息
type DriverMetaOperator interface {
	GetDBMetas() ([]ITableMeta, error)
}

type mysqlDriverMetaOperator struct {
	db *gorm.DB
}

func newMysqlDriverMetaOperator(db *gorm.DB) DriverMetaOperator {
	return &mysqlDriverMetaOperator{
		db: db,
	}
}

func (m *mysqlDriverMetaOperator) GetDBMetas() ([]ITableMeta, error) {
	//TODO implement me
	panic("implement me")
}

type postgreDriverMetaOperator struct {
	db *gorm.DB
}

func (p *postgreDriverMetaOperator) GetDBMetas() ([]ITableMeta, error) {
	//TODO implement me
	panic("implement me")
}

func newPostgreDriverMetaOperator(db *gorm.DB) DriverMetaOperator {
	return &postgreDriverMetaOperator{db: db}
}

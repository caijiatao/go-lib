package adapter

import (
	auditsdk "golib/examples/casbin/audit-sdk"
	"gorm.io/gorm"
)

// MySQLAdapter MySQL 数据库适配器
type MySQLAdapter struct {
	db        *gorm.DB
	tableName string
}

func NewMySQLAdapter(db *gorm.DB, tableName string) *MySQLAdapter {
	return &MySQLAdapter{
		db:        db,
		tableName: tableName,
	}
}

// AutoMigrate 自动创建 MySQL 审计表
func (a *MySQLAdapter) AutoMigrate() error {
	return a.db.Table(a.tableName).AutoMigrate(&auditsdk.CasbinPolicyAudit{})
}

// CreateAuditLog 写入 MySQL 审计日志
func (a *MySQLAdapter) CreateAuditLog(log *auditsdk.CasbinPolicyAudit) error {
	return a.db.Table(a.tableName).Create(log).Error
}

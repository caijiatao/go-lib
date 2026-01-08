package adapter

import (
	auditsdk "golib/examples/casbin/audit-sdk"
	"gorm.io/gorm"
)

// PostgreSQLAdapter PostgreSQL 数据库适配器
type PostgreSQLAdapter struct {
	db        *gorm.DB
	tableName string
}

func NewPostgreSQLAdapter(db *gorm.DB, tableName string) *PostgreSQLAdapter {
	return &PostgreSQLAdapter{
		db:        db,
		tableName: tableName,
	}
}

// AutoMigrate 自动创建 PostgreSQL 审计表
func (a *PostgreSQLAdapter) AutoMigrate() error {
	return a.db.Table(a.tableName).AutoMigrate(&auditsdk.CasbinPolicyAudit{})
}

// CreateAuditLog 写入 PostgreSQL 审计日志
func (a *PostgreSQLAdapter) CreateAuditLog(log *auditsdk.CasbinPolicyAudit) error {
	return a.db.Table(a.tableName).Create(log).Error
}

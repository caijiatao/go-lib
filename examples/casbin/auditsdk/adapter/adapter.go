package adapter

import auditsdk "golib/examples/casbin/audit-sdk"

// DBAdapter 数据库适配器接口
type DBAdapter interface {
	// AutoMigrate 自动创建审计表
	AutoMigrate() error
	// CreateAuditLog 写入审计日志
	CreateAuditLog(log *auditsdk.CasbinPolicyAudit) error
}

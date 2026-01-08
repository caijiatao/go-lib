package auditsdk

import (
	"time"

	"gorm.io/gorm"
)

// AuditConfig SDK 配置项
type AuditConfig struct {
	DB                *gorm.DB      // 数据库连接（由业务方传入，支持 MySQL/PostgreSQL）
	TableName         string        // 审计表名（默认：casbin_policy_audit）
	TenantID          string        // 租户ID（多租户场景用，可选）
	OperatorExtractor func() string // 操作人提取函数（由业务方实现，如从上下文获取当前用户ID）
	RemarkExtractor   func() string // 备注提取函数（可选，如记录操作场景）
	Debug             bool          // 调试模式（是否打印审计日志）
}

// 审计日志结构体（对应数据库表）
type CasbinPolicyAudit struct {
	ID         uint64         `gorm:"primaryKey;autoIncrement" json:"id"`
	TenantID   string         `gorm:"size:64;comment:租户ID" json:"tenant_id"`
	Operator   string         `gorm:"size:64;not null;comment:操作人ID" json:"operator"`
	OpType     string         `gorm:"size:32;not null;comment:操作类型(ADD/REMOVE/UPDATE/SAVE)" json:"op_type"`
	BeforeData string         `gorm:"type:text;comment:变更前数据(JSON)" json:"before_data"`
	AfterData  string         `gorm:"type:text;comment:变更后数据(JSON)" json:"after_data"`
	Remark     string         `gorm:"size:255;comment:备注" json:"remark"`
	CreatedAt  time.Time      `gorm:"comment:操作时间" json:"created_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index;comment:软删除标记" json:"-"`
}

// 默认配置
func DefaultConfig(db *gorm.DB) *AuditConfig {
	return &AuditConfig{
		DB:                db,
		TableName:         "casbin_policy_audit",
		OperatorExtractor: func() string { return "unknown" }, // 默认操作人
		RemarkExtractor:   func() string { return "" },        // 默认无备注
		Debug:             false,
	}
}

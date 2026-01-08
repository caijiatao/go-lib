package auditsdk

import (
	"encoding/json"
	"fmt"
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	"golib/examples/casbin/auditsdk/adapter"
	"time"
)

// AuditEnforcer 带审计功能的 Casbin Enforcer
type AuditEnforcer struct {
	casbin.Enforcer // 嵌入原生 Casbin Enforcer，继承所有方法
	config          *AuditConfig
	dbAdapter       adapter.DBAdapter
	model           model.Model
}

// NewAuditEnforcer 创建带审计功能的 Enforcer
// 参数：
// - modelPath: Casbin 模型文件路径（如 "./model.conf"）
// - casbinAdapter: Casbin 存储适配器（如 MySQL/Redis 适配器）
// - config: 审计配置
func NewAuditEnforcer(modelPath string, casbinAdapter casbin.Adapter, config *AuditConfig) (*AuditEnforcer, error) {
	// 1. 初始化原生 Casbin Enforcer
	enforcer, err := casbin.NewEnforcer(modelPath, casbinAdapter)
	if err != nil {
		return nil, fmt.Errorf("init casbin enforcer failed: %v", err)
	}

	// 2. 初始化数据库适配器（根据 DB 类型自动适配）
	var dbAdapter adapter.DBAdapter
	dialector := config.DB.Dialector.Name()
	switch dialector {
	case "mysql":
		dbAdapter = adapter.NewMySQLAdapter(config.DB, config.TableName)
	case "postgres", "postgresql":
		dbAdapter = adapter.NewPostgreSQLAdapter(config.DB, config.TableName)
	default:
		return nil, fmt.Errorf("unsupported database type: %s", dialector)
	}

	// 3. 自动创建审计表
	if err := dbAdapter.AutoMigrate(); err != nil {
		return nil, fmt.Errorf("auto migrate audit table failed: %v", err)
	}

	// 4. 加载 Casbin 模型
	model, err := enforcer.GetModel()
	if err != nil {
		return nil, fmt.Errorf("get casbin model failed: %v", err)
	}

	// 5. 返回带审计功能的 Enforcer
	auditEnforcer := &AuditEnforcer{
		Enforcer:  *enforcer,
		config:    config,
		dbAdapter: dbAdapter,
		model:     model,
	}

	if config.Debug {
		fmt.Println("casbin audit enforcer init success")
	}
	return auditEnforcer, nil
}

// ------------------------------
// 拦截 AddPolicy 操作，记录审计日志
// ------------------------------
func (a *AuditEnforcer) AddPolicy(params ...interface{}) (bool, error) {
	// 1. 获取变更前数据（当前 Policy 状态）
	beforeData, err := a.getPolicyData()
	if err != nil {
		return false, err
	}

	// 2. 执行原生 AddPolicy
	ok, err := a.Enforcer.AddPolicy(params...)
	if err != nil || !ok {
		return ok, err
	}

	// 3. 获取变更后数据
	afterData, err := a.getPolicyData()
	if err != nil {
		return false, err
	}

	// 4. 记录审计日志（异步，不阻塞主流程）
	go a.logAudit("ADD", beforeData, afterData)

	return ok, nil
}

// ------------------------------
// 拦截 AddPolicies 操作（批量新增）
// ------------------------------
func (a *AuditEnforcer) AddPolicies(rules [][]string) (bool, error) {
	beforeData, err := a.getPolicyData()
	if err != nil {
		return false, err
	}

	ok, err := a.Enforcer.AddPolicies(rules)
	if err != nil || !ok {
		return ok, err
	}

	afterData, err := a.getPolicyData()
	if err != nil {
		return false, err
	}

	go a.logAudit("ADD_BATCH", beforeData, afterData)
	return ok, nil
}

// ------------------------------
// 拦截 RemovePolicy 操作
// ------------------------------
func (a *AuditEnforcer) RemovePolicy(params ...interface{}) (bool, error) {
	beforeData, err := a.getPolicyData()
	if err != nil {
		return false, err
	}

	ok, err := a.Enforcer.RemovePolicy(params...)
	if err != nil || !ok {
		return ok, err
	}

	afterData, err := a.getPolicyData()
	if err != nil {
		return false, err
	}

	go a.logAudit("REMOVE", beforeData, afterData)
	return ok, nil
}

// ------------------------------
// 拦截 RemovePolicies 操作（批量删除）
// ------------------------------
func (a *AuditEnforcer) RemovePolicies(rules [][]string) (bool, error) {
	beforeData, err := a.getPolicyData()
	if err != nil {
		return false, err
	}

	ok, err := a.Enforcer.RemovePolicies(rules)
	if err != nil || !ok {
		return ok, err
	}

	afterData, err := a.getPolicyData()
	if err != nil {
		return false, err
	}

	go a.logAudit("REMOVE_BATCH", beforeData, afterData)
	return ok, nil
}

// ------------------------------
// 拦截 RemoveFilteredPolicy 操作（按条件删除）
// ------------------------------
func (a *AuditEnforcer) RemoveFilteredPolicy(sec, ptype string, fieldIndex int, fieldValues ...string) (bool, error) {
	beforeData, err := a.getPolicyData()
	if err != nil {
		return false, err
	}

	ok, err := a.Enforcer.RemoveFilteredPolicy(sec, ptype, fieldIndex, fieldValues...)
	if err != nil || !ok {
		return ok, err
	}

	afterData, err := a.getPolicyData()
	if err != nil {
		return false, err
	}

	go a.logAudit("REMOVE_FILTERED", beforeData, afterData)
	return ok, nil
}

// ------------------------------
// 拦截 SavePolicy 操作（全量保存）
// ------------------------------
func (a *AuditEnforcer) SavePolicy() error {
	beforeData, err := a.getPolicyData()
	if err != nil {
		return err
	}

	err = a.Enforcer.SavePolicy()
	if err != nil {
		return err
	}

	afterData, err := a.getPolicyData()
	if err != nil {
		return err
	}

	go a.logAudit("SAVE", beforeData, afterData)
	return nil
}

// ------------------------------
// 辅助函数：获取当前 Policy 数据（JSON 格式）
// ------------------------------
func (a *AuditEnforcer) getPolicyData() (string, error) {
	// 获取所有 Policy 规则
	policies, err := a.GetPolicy()
	if err != nil {
		return "", fmt.Errorf("get policy failed: %v", err)
	}

	// 序列化为 JSON 字符串
	data, err := json.MarshalIndent(policies, "", "  ")
	if err != nil {
		return "", fmt.Errorf("marshal policy data failed: %v", err)
	}
	return string(data), nil
}

// ------------------------------
// 辅助函数：写入审计日志
// ------------------------------
func (a *AuditEnforcer) logAudit(opType, beforeData, afterData string) {
	// 构建审计日志
	auditLog := &CasbinPolicyAudit{
		TenantID:   a.config.TenantID,
		Operator:   a.config.OperatorExtractor(), // 业务方自定义提取操作人
		OpType:     opType,
		BeforeData: beforeData,
		AfterData:  afterData,
		Remark:     a.config.RemarkExtractor(), // 业务方自定义提取备注
		CreatedAt:  time.Now(),
	}

	// 写入数据库
	if err := a.dbAdapter.CreateAuditLog(auditLog); err != nil {
		if a.config.Debug {
			fmt.Printf("write audit log failed: %v, log: %+v\n", err, auditLog)
		}
		return
	}

	// 调试模式打印日志
	if a.config.Debug {
		fmt.Printf("audit log recorded: opType=%s, operator=%s, tenantID=%s\n",
			opType, auditLog.Operator, auditLog.TenantID)
	}
}

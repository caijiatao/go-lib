package reconciliation

import (
	"golib/reconciliation/define"
)

type ReconciliationConfig struct {
	SourceDB define.DatabaseConfig

	TargetDBMap map[string]define.DatabaseConfig

	define.FieldsConfig

	Monitor
}

func NewReconciliationConfig(
	sourceDB define.DatabaseConfig,
	targetDBMap map[string]define.DatabaseConfig,
	fieldsConfig define.FieldsConfig,
	monitor Monitor,
) *ReconciliationConfig {
	return &ReconciliationConfig{SourceDB: sourceDB, TargetDBMap: targetDBMap, FieldsConfig: fieldsConfig, Monitor: monitor}
}

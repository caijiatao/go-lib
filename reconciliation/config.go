package reconciliation

import "golib/reconciliation/define"

type ReconciliationConfig struct {
	SourceDB, TargetDB define.DatabaseConfig

	define.FieldsConfig

	Monitor
}

package reconciliation

import "golib/reconciliation/define"

// Monitor
// @Description: 自定义监控上报
type Monitor interface {
	Report(result define.ReconciliationResult)
}

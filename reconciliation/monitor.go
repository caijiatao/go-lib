package reconciliation

import "golib/reconciliation/define"

// Monitor
// @Description: 自定义监控上报
type Monitor interface {
	Report(result define.ReconciliationResult)
	Alert(string)
}

type monitorProxy struct {
	m Monitor
}

func (mp *monitorProxy) Alert(message string) {
	if message == "" {
		return
	}
	if mp.m == nil {
		return
	}
	mp.m.Alert(message)
}

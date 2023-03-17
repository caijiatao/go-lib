package reconciliation

import "github.com/go-mysql-org/go-mysql/canal"

type reconciliationInfo struct {
	c  *canal.Canal
	mp monitorProxy
}

func newReconciliationInfo(c *canal.Canal) *reconciliationInfo {
	return &reconciliationInfo{c: c}
}

package reconciliation

import "testing"

func TestReconciliation(t *testing.T) {
	ron := NewReconciliationImpl()

	ron.RegisterReconciliation(ReconciliationConfig{})

	ron.Run()
}

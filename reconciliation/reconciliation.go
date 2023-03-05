package reconciliation

type ReconciliationImpl struct {
}

func NewReconciliationImpl() *ReconciliationImpl {
	return &ReconciliationImpl{}
}

func (r *ReconciliationImpl) Run() error {
	return nil
}

// RegisterReconciliation
// @Description: 注册对比方法
func (r *ReconciliationImpl) RegisterReconciliation(config ReconciliationConfig) {

}

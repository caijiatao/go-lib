package reconciliation

import (
	"github.com/go-mysql-org/go-mysql/canal"

	"golib/goasync"
	"golib/reconciliation/service"
)

type ReconciliationImpl struct {
	reconInfos []*reconciliationInfo
}

func NewReconciliationImpl() *ReconciliationImpl {
	ron := &ReconciliationImpl{}
	return ron
}

func (r *ReconciliationImpl) Run() error {
	for _, ri := range r.reconInfos {
		copyRi := ri
		go func() {
			defer func() {
				r := recover()
				err := goasync.PanicErrHandler(r)
				copyRi.mp.Alert(err.Error())
			}()

			err := copyRi.c.Run()
			if err != nil {
				copyRi.mp.Alert(err.Error())
			}
		}()
	}
	return nil
}

func (r *ReconciliationImpl) Close() {
	for _, ri := range r.reconInfos {
		ri.c.Close()
	}
}

// RegisterReconciliation
// @Description: 注册对比方法
func (r *ReconciliationImpl) RegisterReconciliation(config ReconciliationConfig) {
	cfg := canal.NewDefaultConfig()
	cfg.Addr = config.SourceDB.Addr
	cfg.User = config.SourceDB.User
	cfg.Password = config.SourceDB.Password
	cfg.Dump.TableDB = config.SourceDB.TableDB
	cfg.Dump.Tables = config.SourceDB.Tables

	c, err := canal.NewCanal(cfg)
	if err != nil {
		panic(any(err))
	}

	// Register a handler to handle RowsEvent

	c.SetEventHandler(service.NewMysqlBinlogHandler())
	r.reconInfos = append(r.reconInfos, newReconciliationInfo(c))

}

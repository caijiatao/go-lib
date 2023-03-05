package binlog_handler

import (
	"fmt"
	"log"

	"github.com/go-mysql-org/go-mysql/canal"
)

const (
	entryTaskSchema = "entry_task"
)

type MysqlBinlogHandler struct {
	canal.DummyEventHandler
	dataHandlerMap map[string]DataHandlerInterface
}

func NewMysqlBinlogHandler() *MysqlBinlogHandler {
	dataHandlerMap := make(map[string]DataHandlerInterface)
	dataHandlerMap[entryTaskSchema] = NewUserDataHandler()
	return &MysqlBinlogHandler{dataHandlerMap: dataHandlerMap}
}

func (h *MysqlBinlogHandler) OnRow(e *canal.RowsEvent) error {
	// "insert [[3 name 1 2021-11-14 18:55:01]] entry_task_db.user &{1636887301 WriteRowsEventV2 123456 50 1395 0}"
	fmt.Printf("%s %v %s %v\n", e.Action, e.Rows, e.Table, e.Header)
	handlerInterface := h.dataHandlerMap[e.Table.Schema]
	err := handlerInterface.SyncBinLogData(e.Action, e.Rows, e.Table.Name)
	if err != nil {
		log.Fatal(err.Error())
		return err
	}
	return nil
}

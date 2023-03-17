package service

import (
	"fmt"
	"log"

	"github.com/go-mysql-org/go-mysql/canal"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"golib/reconciliation/define"
)

type MysqlBinlogHandler struct {
	canal.DummyEventHandler
	dataHandlerMap map[string]DataHandlerInterface
	targetDBMap    map[string]*gorm.DB
}

func NewMysqlBinlogHandler(targetDBMap map[string]define.DatabaseConfig) *MysqlBinlogHandler {
	handler := &MysqlBinlogHandler{}
	for k, v := range targetDBMap {
		dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", v.User, v.Password, v.Addr, v.TableDB)
		var err error
		handler.targetDBMap[k], err = gorm.Open(mysql.Open(dsn))
		if err != nil {
			panic(any(err))
		}
	}
	return handler
}

func (h *MysqlBinlogHandler) parseRows(rows [][]interface{}) []map[string]interface{} {
	return
}

func (h *MysqlBinlogHandler) OnRow(e *canal.RowsEvent) error {
	// "insert [[3 name 1 2021-11-14 18:55:01]] entry_task_db.user &{1636887301 WriteRowsEventV2 123456 50 1395 0}"
	fmt.Printf("%s %v %s %v\n", e.Action, e.Rows, e.Table, e.Header)
	handlerInterface := h.dataHandlerMap[e.Table.Schema]
	changeRowsMap := u.parseUserRows(e.Rows)
	for _, m := range changeRowsMap {

		switch e.Action {
		case "insert":
		case "update":
		case "delete":
		}
	}
	err := handlerInterface.SyncBinLogData(e.Action, e.Rows, e.Table.Name)

	if err != nil {
		log.Fatal(err.Error())
		return err
	}
	return nil
}

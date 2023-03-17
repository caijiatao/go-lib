package service

func (h *MysqlBinlogHandler) String() string {
	return "MysqlBinlogHandler"
}

type DataHandlerInterface interface {
	SyncBinLogData(action string, rows [][]interface{}, table string) error
}

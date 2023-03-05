package main

import (
	"fmt"

	"github.com/go-mysql-org/go-mysql/canal"

	"golib/reconciliation/binlog_handler"
)

func main() {
	cfg := canal.NewDefaultConfig()
	cfg.Addr = "127.0.0.1:33061"
	cfg.User = "root"
	cfg.Password = "secret"

	cfg.Dump.TableDB = "entry_task_db"
	cfg.Dump.Tables = []string{"user"}

	c, err := canal.NewCanal(cfg)
	if err != nil {
		fmt.Println(err)
	}

	// Register a handler to handle RowsEvent
	c.SetEventHandler(&binlog_handler.MysqlBinlogHandler{})

	// Start canal
	err = c.Run()
	if err != nil {
		fmt.Println(err)
	}
}

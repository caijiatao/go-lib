package db

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

func initDB() (*sql.DB, error) {
	connStr := "host=192.168.16.21 user=postgres password=123456 dbname=airec_server port=10432 sslmode=disable TimeZone=Asia/Shanghai"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open a DB connection: %v", err)
	}

	// 连接池设置
	db.SetMaxOpenConns(10)   // 最大打开连接数
	db.SetMaxIdleConns(5)    // 最大空闲连接数
	db.SetConnMaxLifetime(0) // 连接的最大生命周期

	// 尝试ping数据库以确认连接成功
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to connect to DB: %v", err)
	}

	return db, nil
}

func queryExample(db *sql.DB) error {
	rows, err := db.Query("SELECT id FROM airec_server.standard_mapping")
	if err != nil {
		return fmt.Errorf("query failed: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return fmt.Errorf("scan failed: %v", err)
		}
		fmt.Printf("standardMappingId: %v\n", id)
	}

	if err := rows.Err(); err != nil {
		return fmt.Errorf("row iteration failed: %v", err)
	}

	return nil
}

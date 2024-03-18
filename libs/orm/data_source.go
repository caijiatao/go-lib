package orm

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"golib/libs/logger"
	"golib/libs/util"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"strings"
)

type BasicDBConfig struct {
	Port     string `json:"port"`
	Host     string `json:"host"`
	User     string `json:"user"`
	Password string `json:"password"` // 密码加密存储
	DbName   string `json:"dbname"`
}

// SourceDBConfig
// @Description: 源数据库信息
type SourceDBConfig struct {
	BasicDBConfig
	Type uint64
}

func NewSourceDBConfig(sourceType uint64, config string) *SourceDBConfig {
	dbConfig := &SourceDBConfig{}
	err := json.Unmarshal([]byte(config), dbConfig)
	password, err := util.Decrypt(dbConfig.Password)
	dbConfig.Password = password
	dbConfig.Type = sourceType
	if err != nil {
		return nil
	}
	return dbConfig

}

func (d *SourceDBConfig) GetDSN() string {
	switch d.Type {
	case PGSQLSourceType, FileSourceType: // 文件数据源也使用的是postgres
		if d.Password == "" {
			// password为空时，连接数据库有bug
			return fmt.Sprintf("host=%s user=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai", d.Host, d.User, d.DbName, d.Port)
		}
		return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai", d.Host, d.User, d.Password, d.DbName, d.Port)
	case MysqlSourceType:
		return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?timeout=2s&readTimeout=2s&parseTime=true", d.User, d.Password, d.Host, d.Port, d.DbName)
	case HiveSourceType:
		return fmt.Sprintf("%s/%s@%s:%s", d.User, d.Password, d.Host, d.Port)
	}
	return ""
}

func (d *SourceDBConfig) GetDail() gorm.Dialector {
	switch d.Type {
	case PGSQLSourceType, FileSourceType: // 文件数据源也使用的是postgres
		return postgres.Open(d.GetDSN())
	case MysqlSourceType:
		return mysql.Open(d.GetDSN())
	default:
		return nil
	}
}

func GetPGDBConfigByDSN(dsn string) (*SourceDBConfig, error) {
	config := &SourceDBConfig{}
	keyValuePairs := strings.Split(dsn, " ")
	for _, kvPair := range keyValuePairs {
		keyValue := strings.Split(kvPair, "=")
		if len(keyValue) != 2 {
			logger.Errorf("dsn parse failed: %s", dsn)
			return nil, errors.New("dsn parse failed")
		}
		key := keyValue[0]
		value := keyValue[1]
		switch key {
		case "host":
			config.Host = value
		case "user":
			config.User = value
		case "password":
			config.Password = value
		case "dbname":
			config.DbName = value
		case "port":
			config.Port = value
		}
	}
	return config, nil
}

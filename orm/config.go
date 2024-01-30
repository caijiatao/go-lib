package orm

import (
	"github.com/beltran/gohive"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const postgresConfigKey = "postgres_config"
const postgresDemoConfigKey = "postgres_demo_config"
const postgresDemoServerKey = "demo_server"

var DataSourceConfigDSNMap map[string]string

type Config struct {
	DBClientName string
	Dial         gorm.Dialector
	MaxIdle      int
	MaxOpen      int
	*gorm.Config
	//HiveConfig
	SourceConfig *SourceDBConfig
}

type HiveConfig struct {
	Host       string
	Port       int
	auth       string
	Connection *gohive.ConnectConfiguration
}

func readConfigs(configKey string) []*Config {
	DataSourceConfigDSNMap = make(map[string]string, 0)
	configMap := viper.GetStringMap(configKey)
	configs := make([]*Config, 0)
	for key, _ := range configMap {
		pgConfig := viper.Sub(configKey).Sub(key)
		dsn := pgConfig.GetString("dsn")
		DataSourceConfigDSNMap[key] = dsn

		config := &Config{
			Config: &gorm.Config{
				Logger: logger.Default.LogMode(logger.Info),
			},
			Dial:         postgres.Open(dsn),
			MaxIdle:      pgConfig.GetInt("max_idle"),
			MaxOpen:      pgConfig.GetInt("max_open"),
			DBClientName: key,
			SourceConfig: &SourceDBConfig{
				Type: PGSQLSourceType,
			},
		}
		configs = append(configs, config)
	}
	return configs
}

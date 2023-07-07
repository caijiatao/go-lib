package orm

import (
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const postgresConfigKey = "postgres_config"

type Config struct {
	DBClientName string
	Dial         gorm.Dialector
	MaxIdle      int
	MaxOpen      int
	*gorm.Config
}

func readConfigs() []*Config {
	configMap := viper.GetStringMap(postgresConfigKey)
	configs := make([]*Config, 0)
	for key, _ := range configMap {
		pgConfig := viper.Sub(postgresConfigKey).Sub(key)
		dsn := pgConfig.GetString("dsn")

		config := &Config{
			Config: &gorm.Config{
				Logger: logger.Default.LogMode(logger.Info),
			},
			Dial:         postgres.Open(dsn),
			MaxIdle:      pgConfig.GetInt("max_idle"),
			MaxOpen:      pgConfig.GetInt("max_open"),
			DBClientName: key,
		}
		configs = append(configs, config)
	}
	return configs
}

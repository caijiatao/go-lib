package logger

import "github.com/spf13/viper"

var config = &Config{}

type Config struct {
	ProjectName  string
	Level        string `yaml:"level"`
	Path         string `yaml:"path"`
	MaxAge       int64  `yaml:"max_age"`
	RotationTime int64  `yaml:"rotation_time"`
}

// logLevel 日志级别
type logLevel int

// levelMap 字符串和级别映射
func levelMap() map[string]logLevel {
	return map[string]logLevel{
		"debug": logDebug,
		"info":  logInfo,
		"warn":  logWarn,
		"err":   logErr,
		"off":   logOff,
	}
}

const (
	logDebug logLevel = iota
	logInfo
	logWarn
	logErr
	logOff
)

func Init() {
	viper.SetDefault("server_config.server_name", "default_server_name")
	config.ProjectName = viper.GetString("server_config.server_name")
	viper.SetDefault("logger_config.level", "debug")
	config.Level = viper.GetString("logger_config.level")
	viper.SetDefault("logger_config.path", "./var/log")
	config.Path = viper.GetString("logger_config.path")
	viper.SetDefault("logger_config.max_age", "30")
	config.MaxAge = viper.GetInt64("logger_config.max_age")
	viper.SetDefault("logger_config.rotation_time", "24")
	config.RotationTime = viper.GetInt64("logger_config.rotation_time")
}

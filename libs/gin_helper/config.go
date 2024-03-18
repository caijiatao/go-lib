package gin_helper

import "github.com/spf13/viper"

const serverConfigKey = "server_config"

var serverConfig = &ServerConfig{}

type ServerConfig struct {
	HttpPort  string `yaml:"http_port"`
	HttpUse   bool   `yaml:"https_use"`
	HttpsPort string `yaml:"https_port"`
}

func Init() error {
	viper.SetDefault("server_config.http_port", "8080")
	serverConfig.HttpPort = viper.GetString("server_config.http_port")
	viper.SetDefault("server_config.https_use", false)
	serverConfig.HttpsPort = viper.GetString("server_config.https_use")
	viper.SetDefault("server_config.https_port", "8443")
	serverConfig.HttpsPort = viper.GetString("server_config.https_port")
	return nil
}

func GetServerConfig() *ServerConfig {
	return serverConfig
}

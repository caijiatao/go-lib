package sms

import (
	"fmt"
	"github.com/spf13/viper"
	"sync"
)

const (
	smsConfigKey = "sms_config"
)

var smsConfig SmsConfig
var configOnce sync.Once

type SmsConfig struct {
	AccessID     string
	SecretKey    string
	SignName     string
	TemplateCode string
	Endpoint     string
}

func initConfig() {
	configOnce.Do(func() {
		smsConfig = SmsConfig{
			AccessID:     viper.GetString(fmt.Sprintf("%s.%s", smsConfigKey, "access_id")),
			SecretKey:    viper.GetString(fmt.Sprintf("%s.%s", smsConfigKey, "secret_key")),
			SignName:     viper.GetString(fmt.Sprintf("%s.%s", smsConfigKey, "sign_name")),
			TemplateCode: viper.GetString(fmt.Sprintf("%s.%s", smsConfigKey, "template_code")),
			Endpoint:     viper.GetString(fmt.Sprintf("%s.%s", smsConfigKey, "endpoint")),
		}
	})
}

func Config() SmsConfig {
	return smsConfig
}

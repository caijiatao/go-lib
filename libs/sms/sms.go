package sms

import (
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	dysmsapi20170525 "github.com/alibabacloud-go/dysmsapi-20170525/v3/client"
	"github.com/pkg/errors"
	"golib/libs/logger"
	"sync"
)

var client *SmsClient
var clientOnce sync.Once
var clientErr error

type SmsClient struct {
	AliCli *dysmsapi20170525.Client
}

func Init() error {
	initConfig()
	clientOnce.Do(func() {
		client = &SmsClient{}
		config := Config()
		client.AliCli, clientErr = CreateClient(config)

	})
	if clientErr != nil {
		return clientErr
	}
	return nil
}

func Client() *SmsClient {
	return client
}

func CreateClient(config SmsConfig) (_result *dysmsapi20170525.Client, _err error) {
	conf := &openapi.Config{
		// 必填，您的 AccessKey ID
		AccessKeyId: &config.AccessID,
		// 必填，您的 AccessKey Secret
		AccessKeySecret: &config.SecretKey,
		Endpoint:        &config.Endpoint,
	}

	return dysmsapi20170525.NewClient(conf)
}

func (self *SmsClient) SendSms(phoneNumber string, signName string, templateCode string, code string) error {
	req := &dysmsapi20170525.SendSmsRequest{
		PhoneNumbers:  &phoneNumber,
		SignName:      &signName,
		TemplateCode:  &templateCode,
		TemplateParam: &code,
	}
	resp, err := self.AliCli.SendSms(req)
	if err != nil {
		return err
	}
	if *resp.Body.Code != "OK" {
		return errors.New(*resp.Body.Message)
	}
	logger.Info("phone:%s, code:%s, requestID:%s", phoneNumber, code, resp.Body.RequestId)
	return nil
}

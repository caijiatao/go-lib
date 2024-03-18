package mq_test_util

import "golib/libs/mq"

// InitTestKafka
//
//	@Description: 单元测试初始化使用
func InitTestKafka() {
	_, err := mq.NewClient(mq.NewConfig("test", []string{"192.168.13.213:10092", "192.168.13.214:10092", "192.168.13.215:10092"}))
	if err != nil {
		panic(err)
	}
}

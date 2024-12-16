package naming

import (
	"golib/libs/net_helper"
	"log"
	"testing"
	"time"
)

func TestRegister(t *testing.T) {
	var endpoints = []string{"localhost:2379"}
	ipv4, err := net_helper.GetLocalIP()
	if err != nil {
		log.Fatalln(err)
	}
	//ser, err := NewServiceRegister(endpoints, "/api_server/", ipv4, 5)
	_ = NewServiceRegister(endpoints, "/api_server/", ipv4)
	if err != nil {
		log.Fatalln(err)
	}
	//监听续租相应chan
	select {
	case <-time.After(time.Minute):

		//ser.Close()
	}
}

package naming

import (
	"log"
	"testing"
	"time"
)

func TestRegister(t *testing.T) {
	var endpoints = []string{"localhost:2379"}
	ser, err := NewServiceRegister(endpoints, "/api_server/", "localhost:8000", 5)
	if err != nil {
		log.Fatalln(err)
	}
	//监听续租相应chan
	go ser.ListenLeaseRespChan()
	select {
	case <-time.After(20 * time.Second):
		ser.Close()
	}
}

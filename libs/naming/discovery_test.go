package naming

import (
	"log"
	"testing"
	"time"
)

func TestDiscovery(t *testing.T) {
	var endpoints = []string{"localhost:2379"}
	ser := NewServiceDiscovery(endpoints)
	ser.WatchService("/web/")
	ser.WatchService("/gRPC/")
	for {
		select {
		case <-time.Tick(10 * time.Second):
			log.Println(ser.GetServices())
		}
	}
}

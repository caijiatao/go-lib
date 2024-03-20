package naming

import (
	"log"
	"testing"
	"time"
)

func TestDiscovery(t *testing.T) {
	var endpoints = []string{"localhost:2379"}
	ser := NewServiceDiscovery(endpoints)
	ser.runWatchService("/web/")
	ser.runWatchService("/gRPC/")
	for {
		select {
		case <-time.Tick(10 * time.Second):
			log.Println(ser.GetServices())
		}
	}
}

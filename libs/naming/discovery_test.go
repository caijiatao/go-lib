package naming

import (
	"log"
	"testing"
	"time"
)

func TestDiscovery(t *testing.T) {
	var endpoints = []string{"localhost:2379"}
	ser := NewServiceDiscovery(endpoints, "/api_server/")
	for {
		select {
		case <-time.Tick(10 * time.Second):
			log.Println(ser.GetServices())
		}
	}
}

package naming

import (
	"fmt"
	"testing"
	"time"
)

func TestDiscovery(t *testing.T) {
	var endpoints = []string{"localhost:2379"}
	ser := NewServiceDiscovery(endpoints, "/api_server/")
	for {
		select {
		case <-time.Tick(3 * time.Second):
			fmt.Println(ser.GetServices())
			fmt.Println(ser.GetServiceList())
		}
	}
}

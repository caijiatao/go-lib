package net_helper

import (
	"fmt"
	"testing"
)

func TestGetLocalIP(t *testing.T) {
	gotIpv4, err := GetLocalIP()
	fmt.Println(gotIpv4)
	fmt.Println(err)
}

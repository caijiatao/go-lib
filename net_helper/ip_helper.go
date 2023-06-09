package net_helper

import (
	"errors"
	"net"
)

func GetLocalIP() (ipv4 string, err error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}

	for _, addr := range addrs {
		ipNet, isIpNet := addr.(*net.IPNet)
		if !isIpNet {
			continue
		}

		if ipNet.IP.IsLoopback() {
			continue
		}

		// 过滤掉ipv6的地址
		if ipNet.IP.To4() == nil {
			continue
		}
		return ipNet.IP.To4().String(), nil
	}

	return "", errors.New("no local ip")
}

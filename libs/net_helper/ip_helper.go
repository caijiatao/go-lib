package net_helper

import (
	"errors"
	"github.com/zeromicro/go-zero/core/netx"
	"net"
	"os"
	"strings"
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

var (
	allEths  = "0.0.0.0"
	envPodIp = "POD_IP"
)

func GetFigureOutListenOn(listenOn string) string {
	fields := strings.Split(listenOn, ":")
	if len(fields) == 0 {
		return listenOn
	}

	host := fields[0]
	if len(host) > 0 && host != allEths {
		return listenOn
	}

	ip := os.Getenv(envPodIp)
	if len(ip) == 0 {
		ip = netx.InternalIp()
	}
	if len(ip) == 0 {
		return listenOn
	}

	return strings.Join(append([]string{ip}, fields[1:]...), ":")
}

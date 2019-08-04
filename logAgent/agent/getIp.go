package main

import (
	"fmt"
	"net"
)

var (
	// 本地IP
	localIpArray []string
)

func init() {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		panic(fmt.Sprintf("获取本地 Ip失败,%s", err))
	}

	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				localIpArray = append(localIpArray, ipnet.IP.String())
			}
		}
	}
}

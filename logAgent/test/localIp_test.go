package test

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"net"
	"testing"
)

func TestLocalIp(t *testing.T) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		panic(fmt.Sprintf("获取本地 Ip失败%s", err))
	}

	for _, addr := range addrs {
		logs.Info("addr", addr)
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				logs.Info("ip", ipnet.IP.String())
			}
		}
	}
}

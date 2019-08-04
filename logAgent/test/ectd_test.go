package test

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/logs"
	etcd "go.etcd.io/etcd/clientv3"
	"testing"
	"time"
)

type Collect struct {
	Topic   string `json:"topic"`
	LogPath string `json:"logPath"`
}

type LogConfig struct {
	Topic    string `json:"topic"`
	LogPath  string `json:"log_path"`
	Service  string `json:"service"`
	SendRate int    `json:"send_rate"`
}

func TestEtcd(t *testing.T) {
	//初始化一个etcd 的 client

	var etcdAddress = make([]string, 3)
	etcdAddress = append(etcdAddress, "127.0.0.1:2379")
	client, err := etcd.New(etcd.Config{
		Endpoints:   etcdAddress,
		DialTimeout: 5 * time.Second,
	})

	fmt.Print(client)

	if err != nil {
		logs.Info("初始化 ectd错误", err)
		return
	}

	config := &LogConfig{
		Topic:    "nginx_access",
		LogPath:  "/Users/panyundong/Desktop/logs/access.log",
		Service:  "nginx_server",
		SendRate: 100,
	}

	var key string = "/logs/192.168.0.102/nginx/"

	bytes, _ := json.Marshal(config)
	_, _ = client.Put(context.Background(), key, string(bytes))
	response, err := client.Get(context.Background(), key)
	for key, value := range response.Kvs {
		fmt.Printf(string(key), value)
	}

}

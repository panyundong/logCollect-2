package main

import (
	"context"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/coreos/etcd/mvcc/mvccpb"
	etcd "go.etcd.io/etcd/clientv3"
	"sync"
	"time"
)

var (
	etcdClient *etcd.Client
	confChan   = make(chan string, 20)
	waitGroup  sync.WaitGroup
)

func initEtcd(etcdAddress []string, collectKey string) (err error) {
	//初始化一个etcd 的 client
	cli, err := etcd.New(etcd.Config{
		Endpoints:   etcdAddress,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		logs.Error("连接 ectd 失败,%s", err)
		return
	}
	etcdClient = cli

	var etcdKeys []string
	//封装 etcd的key 类似于"/logs/172.0.16.60/nginx/"
	for _, value := range localIpArray {
		key := fmt.Sprintf(collectKey, value)
		etcdKeys = append(etcdKeys, key)
	}

	//从 ectd中拉取key 回来
	for _, key := range etcdKeys {
		ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
		response, err := etcdClient.Get(ctx, key)
		cancelFunc()
		if err != nil {
			logs.Error("get key from etcd error:", err)
			continue
		}

		//把获取到的值往管道里扔
		for _, value := range response.Kvs {
			confChan <- string(value.Value)
			logs.Info("get key from etcd success:", string(value.Value))
		}

	}

	waitGroup.Add(1)

	//去监听 key 的变化
	go etcdWatch(etcdKeys)
	return
}

func etcdWatch(etcdKeys []string) {
	defer waitGroup.Done()
	var watchChans []etcd.WatchChan
	logs.Info("etcdClient = %s", etcdClient)
	for _, key := range etcdKeys {

		//返回一个 watch的 Chan
		watch := etcdClient.Watch(context.Background(), key)
		watchChans = append(watchChans, watch)
	}

	//死循环  一直遍历这些watchChans
	for {
		for _, watch := range watchChans {
			select {
			case wresp := <-watch:
				for _, value := range wresp.Events {

					// key 更新
					if mvccpb.PUT == value.Type {
						confChan <- string(value.Kv.Value)
					}

					//key 删除
					//TODO 删除后停止监听日志文件
					if value.Type == mvccpb.DELETE {
						logs.Warn("删除监听的 key = %s value =%s", string(value.Kv.Key), string(value.Kv.Value))
					}
				}
			default:
				//logs.Warn("do nothing..")
			}

			//短暂休眠一秒吧
			time.Sleep(time.Second)
		}
	}
}

func GetEtcdConfChan() chan string {
	return confChan
}

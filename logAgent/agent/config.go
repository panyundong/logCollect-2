package main

import (
	"github.com/astaxie/beego/config"
	"github.com/astaxie/beego/logs"
	"strings"
)

//agent配置的结构体
type Config struct {
	EtcdAddress  []string //Etcd地址
	EtcdWatchKey string   //Etcd监听的 key

	KafkaAddress []string //kafka 地址

	LogLevel string //agent 的日志级别
	LogPath  string //agent 的日志路径

	ThreadNum int //并发收集的线程数
}

var (
	// 全局agent配置的对象
	agentConfig = &Config{}
)

func initConfig(adapterName string, fileName string) (err error) {
	configer, err := config.NewConfig(adapterName, fileName)
	if err != nil {
		logs.Error("初始化日志失败,%s", err)
		return
	}

	//获取 Etcd的地址
	etcAdders := configer.String("etcd::etcd_address")
	if len(etcAdders) == 0 {
		logs.Error("找不到Ectd的地址")
		return
	}
	agentConfig.EtcdAddress = strings.Split(etcAdders, ",")

	//获取 Etcd的监听的 key
	etcdWatchkey := configer.String("etcd::etcd_watch_key")
	if len(etcdWatchkey) == 0 {
		logs.Error("找不到Etcd的监听的 key")
		return
	}
	agentConfig.EtcdWatchKey = etcdWatchkey

	//获取 kafka的地址
	kafkaAddress := configer.String("kafka::kafka_address")
	if len(kafkaAddress) == 0 {
		logs.Error("找不到kafka的地址")
		return
	}
	agentConfig.KafkaAddress = strings.Split(kafkaAddress, ",")

	//获取并发收集的线程数
	threadNum, err := configer.Int("kafka::thread_num")
	if err != nil {
		logs.Warn("找不到并发收集的线程数,设置为 2 默认值")
		threadNum = 2
	}
	agentConfig.ThreadNum = threadNum

	//获取 agent log 的日志地址
	logPath := configer.String("base::log_path")
	if len(logPath) == 0 {
		logs.Error("找不到agent的日志地址")
		return
	}
	agentConfig.LogPath = logPath

	//获取 agent 的日志级别
	logLevel := configer.String("base::log_level")
	if len(logLevel) == 0 {
		logs.Warn("找不到agent的日志logLevel,设置为 DEBUG 默认值")
		logLevel = "debug"
	}
	agentConfig.LogLevel = logLevel

	return

}

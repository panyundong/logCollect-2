package main

import (
	"errors"
	"github.com/astaxie/beego/config"
	"github.com/astaxie/beego/logs"
	"github.com/panyundong/logCollect-2/logAgent/tailf"
	"strings"
)

type Config struct {
	LogLevel     string
	LogPath      string
	ChanSize     int
	KafkaAddress []string
	EtcdAddress  []string
	CollectKey   string
	Collects     []tailf.Collect
	Ip           string
}

var (
	// 配置信息对象
	agentConfig *Config
)

// 加载配置信息
func LoadConfig(configType string, configPath string) (err error) {
	configer, err := config.NewConfig(configType, configPath)
	if err != nil {
		logs.Error("加载配置文件失败", err)
		return
	}

	agentConfig = &Config{}

	// 获取基础配置
	err = getAgentConfig(configer)
	if err != nil {
		return
	}
	return

}

func getAgentConfig(conf config.Configer) (err error) {
	// 获取日志级别
	logLevel := conf.String("base::log_level")
	if len(logLevel) == 0 {
		logLevel = "debug"
	}
	agentConfig.LogLevel = logLevel

	// 获取日志路径
	logPath := conf.String("base::log_path")
	if len(logPath) == 0 {
		logPath = "./logs/logagent.log"
	}

	agentConfig.LogPath = logPath

	// 日志收集开启chan大小
	chanSize, err := conf.Int("base::queue_size")
	if err != nil {
		chanSize = 200
	}
	agentConfig.ChanSize = chanSize

	//获取 etcd 地址
	etcdAddress := conf.String("etcd::etcd_address")
	if len(etcdAddress) == 0 {
		err = errors.New("找不到 etcd 地址哦!")
		return
	}
	agentConfig.EtcdAddress = strings.Split(etcdAddress, ",")

	kafkaAddress := conf.String("kafka::kafka_address")
	if len(kafkaAddress) == 0 {
		err = errors.New("找不到kafka 的地址哦!")
		return
	}

	agentConfig.KafkaAddress = strings.Split(kafkaAddress, ",")

	//获取日志收集的前缀
	collectKey := conf.String("collect::collectKey")
	if len(collectKey) == 0 {
		err = errors.New("找不到日志收集前缀key")
		return
	}

	agentConfig.CollectKey = collectKey
	return
}

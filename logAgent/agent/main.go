package main

import (
	"github.com/astaxie/beego/logs"
)

func main() {

	//加载配置文件
	err := initConfig("ini", "./logAgent/config/config.ini")
	if err != nil {
		logs.Error("加载配置文件失败,%s", err)
		return
	}
	logs.Info("加载配置文件成功. file =[%s]", agentConfig)

	//初始化日志
	d := initAgentLog(agentConfig.LogPath, agentConfig.LogLevel)
	if d != nil {
		logs.Error("初始化日志失败,", err)
		return
	}
	logs.Debug("初始化日志成功. logPath = [%s]", agentConfig.LogPath)

	//初始化Etcd
	e := initEtcd(agentConfig.EtcdAddress, agentConfig.EtcdWatchKey)
	if e != nil {
		logs.Error("初始化Etcd", err)
		return
	}
	logs.Info("初始化Etcd成功")

	//初始化kafka
	err = InitKafka(agentConfig.KafkaAddress, agentConfig.ThreadNum)
	if err != nil {
		logs.Error("Start logAgent [init kafka] failed, err:", err)
		return
	}
	logs.Debug("初始化kafka 成功")

	runServer()
	logs.Info("logagent 服务退出了")

}

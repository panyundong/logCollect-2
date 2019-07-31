package main

import "github.com/astaxie/beego/logs"

func main() {

	//加载配置文件
	err := LoadConfig("ini", "./logAgent/config/config.ini")
	if err != nil {
		logs.Error("%s", err)
		return
	}
	logs.Debug("加载配置文件成功. file =[%s]", "./logAgent/config/config.ini")

	//初始化日志
	d := initAgentLog()
	if d != nil {
		logs.Error("初始化日志失败", err)
		return
	}
	logs.Debug("初始化日志成功.")

	//初始化Etcd

	//初始化tailf

	//初始化kafka

	//启动logagent服务

}

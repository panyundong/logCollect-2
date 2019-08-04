package main

import (
	"encoding/json"
	"github.com/astaxie/beego/logs"
)

//初始化日志
func initAgentLog(path string, loglevel string) (err error) {
	//日志的配置文件呢,它需要一个 map 参数
	logconfig := make(map[string]interface{})
	logconfig["fileName"] = path
	logconfig["level"] = convertLogLevel(loglevel)
	logconfig["color"] = true //日志级别彩色输出

	bytes, err := json.Marshal(logconfig)
	if err != nil {
		logs.Error("序列化json失败,%s", err)
		return
	}

	err = logs.SetLogger(logs.AdapterConsole, string(bytes)) //打印到控制台
	err = logs.SetLogger(logs.AdapterFile, string(bytes))    //打印到文件
	logs.SetLogFuncCall(true)                                //打印行号和方法名
	logs.SetLogFuncCallDepth(3)                              //调用层次,递归调用时,准确输出行号

	return
}

func convertLogLevel(level string) (logLevel int) {
	switch level {
	case "DEBUG":
		logLevel = logs.LevelDebug
		break
	case "INFO":
		logLevel = logs.LevelInfo
		break
	case "WARN":
		logLevel = logs.LevelWarn
		break
	case "ERROR":
		logLevel = logs.LevelError
		break
	default:
		logLevel = logs.LevelInfo
	}
	return logLevel
}

package main

import (
	"encoding/json"
	"github.com/astaxie/beego/logs"
)

//初始化日志
func initAgentLog() (err error) {

	//日志的配置文件呢需要一个 map 参数
	logconfig := make(map[string]interface{})
	logconfig["fileName"] = agentConfig.LogPath
	logconfig["level"] = convertLogLevel(agentConfig.LogLevel)
	logconfig["color"] = true

	bytes, err := json.Marshal(logconfig)
	if err != nil {
		logs.Error("序列化接送失败,%s", err)
		return
	}

	err = logs.SetLogger(logs.AdapterConsole, string(bytes))
	err = logs.SetLogger(logs.AdapterFile, string(bytes))
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

package main

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/hpcloud/tail"
	"strings"
	"sync"
)

type LogConfig struct {
	Topic    string `json:"topic"`
	LogPath  string `json:"log_path"`
	Service  string `json:"service"`
	SendRate int    `json:"send_rate"`
}

//日志收集对象
type TailObj struct {
	tail     *tail.Tail   //日志收集组件
	offset   int64        //收集的偏移量
	logConf  LogConfig    //日志的配置
	secLimit *SecondLimit //流量速率
	exitChan chan bool    //退出标识
}

//统一日志收集对象管理
type TailMgr struct {
	tailObjMap map[string]*TailObj //多个日志对象,map管理
	lock       sync.Mutex          //锁
}

var tailMgr *TailMgr

func runServer() {
	tailMgr = NewTailMgr()
	tailMgr.Process()
	waitGroup.Wait()

}

//初始化日志管理的对象
func NewTailMgr() *TailMgr {
	return &TailMgr{
		tailObjMap: make(map[string]*TailObj, 16),
	}
}

//收集日志
func (mgr TailMgr) Process() {
	for config := range GetEtcdConfChan() {
		logs.Info("config form etcd , %s", config)
		var logConfArr []LogConfig
		logConfig := &LogConfig{}
		err := json.Unmarshal([]byte(config), logConfig)
		logConfArr = append(logConfArr, *logConfig)
		if err != nil {
			logs.Error("unmarshal failed, err: %v conf :%s", err, config)
			continue
		}

		err = mgr.reloadConfig(logConfArr)
		if err != nil {
			logs.Error("reload config from etcd failed: %v", err)
			continue
		}
		logs.Debug("reload config from etcd success")
	}
}

func (mgr TailMgr) reloadConfig(configs []LogConfig) (err error) {
	for _, value := range configs {
		obj, ok := mgr.tailObjMap[value.LogPath]
		if !ok {
			err = mgr.AddLogFile(value)
			if err != nil {
				logs.Error("add log file failed:%v", err)
				continue
			}
			continue
		}
		obj.logConf = value
		obj.secLimit.limit = int32(value.SendRate)
		mgr.tailObjMap[value.LogPath] = obj
	}

	for key, value := range mgr.tailObjMap {
		var found = false
		for _, newConfig := range configs {
			if key == newConfig.LogPath {
				found = true
				break
			}
		}
		if found == false {
			logs.Warn("log path :%s is remove", key)
			value.exitChan <- true
			delete(mgr.tailObjMap, key)
		}
	}
	return
}

func (mgr TailMgr) AddLogFile(config LogConfig) (err error) {
	mgr.lock.Lock()
	defer mgr.lock.Unlock()

	_, ok := mgr.tailObjMap[config.LogPath]
	if ok {
		err = fmt.Errorf("duplicate filename:%s", config.LogPath)
		return
	}

	//创建日志收集对象
	tail2, err := tail.TailFile(config.LogPath, tail.Config{
		Location:    &tail.SeekInfo{Offset: 0, Whence: 2},
		ReOpen:      true,
		MustExist:   false,
		Poll:        true,
		RateLimiter: nil,
		Follow:      true,
		MaxLineSize: 0,
		Logger:      nil,
	})

	if err != nil {
		logs.Error("创建日志收集对象失败,%s", err)
		return
	}

	//创建日志对象
	obj := &TailObj{
		tail:     tail2,
		offset:   0,
		logConf:  config,
		secLimit: NewSecondLimit(int32(config.SendRate)),
		exitChan: make(chan bool, 1),
	}

	mgr.tailObjMap[config.LogPath] = obj

	waitGroup.Add(1)
	go obj.readLog()
	return
}

func (obj *TailObj) readLog() {

	for line := range obj.tail.Lines {
		if line.Err != nil {
			logs.Error("read line error:%v ", line.Err)
			continue
		}

		trimSpace := strings.TrimSpace(line.Text)
		if len(trimSpace) == 0 || trimSpace[0] == '\n' {
			continue
		}

		kafkaSend.addKafkaMessage(line.Text, obj.logConf.Topic)
		obj.secLimit.Add(1)
		obj.secLimit.Wait()

		select {
		case <-obj.exitChan:
			logs.Warn("tail obj is exited: config:", obj.logConf)
			return
		default:
		}

	}
	waitGroup.Done()
}

package tailf

import (
	"errors"
	"github.com/astaxie/beego/logs"
	"github.com/hpcloud/tail"
	"sync"
)

type Collect struct {
	Topic   string `json:"topic"`
	LogPath string `json:"logPath"`
}

//发送到 kafka的消息
type KafkaMsg struct {
	Msg string `json:"msg"`
	Ip  string `json:"ip"`
}

//发送消息的结构体
type TextMsg struct {
	Msg   KafkaMsg
	Topic string
}

//tailf 任务对象
type TailObj struct {
	tailObj  *tail.Tail
	collect  Collect
	status   int
	exitChan chan int
}

// tailf任务对象管理
type TailsObjMgr struct {
	tailObjs []*TailObj
	msgChan  chan *TextMsg
	lock     sync.Mutex
}

var (
	//初始化tailf任务对象管理
	tailObjMgr *TailsObjMgr
	hostIp     string
)

const (
	StatusNormal = 1
	StatusDelete = 2
)

//初始化tailf
func InitTail(collects []Collect, chanSize int, ip string) (err error) {
	tailObjMgr = &TailsObjMgr{
		msgChan: make(chan *TextMsg, chanSize),
	}
	if len(collects) == 0 {
		err = errors.New("collect task is nill")
		logs.Warn("collect task is nill")
	}

	// 创建tailf task
	for _, v := range collects {
		createTask(v)
	}
	return
}

func createTask(collect Collect) {
	file, e := tail.TailFile(collect.LogPath, tail.Config{
		Location:    &tail.SeekInfo{Offset: 0, Whence: 2},
		ReOpen:      true,
		MustExist:   false,
		Poll:        true,
		Pipe:        false,
		RateLimiter: nil,
		Follow:      true,
		MaxLineSize: 0,
		Logger:      nil,
	})

	if e != nil {
		logs.Warn("tailf create [%v] failed, %v", collect.LogPath, e)
		return
	}

	tailObj := &TailObj{
		tailObj:  file,
		collect:  collect,
		exitChan: make(chan int, 1),
	}

	// 开启goroute去读取监听日志的内容
	go readFromTail(tailObj, collect.Topic)
}

func readFromTail(tailObj *TailObj, topic string) {
	for true {
		select {
		//读取日志
		case linMsg, ok := <-tailObj.tailObj.Lines:
			if !ok {
				logs.Warn("read obj:[%v] topic:[%v] filed continue", tailObj, topic)
				continue
			}
			if linMsg.Text == "" {
				continue
			}

			kafkaMsg := KafkaMsg{
				Msg: linMsg.Text,
				Ip:  hostIp,
			}
			msgObj := &TextMsg{
				Msg:   kafkaMsg,
				Topic: topic,
			}
			tailObjMgr.msgChan <- msgObj

			//停止
		case <-tailObj.exitChan:
			logs.Warn("tail obj will exited, conf:%v", tailObj.collect)
			return
		}
	}
}

//更新 tail 任务
func UpdateTailfTask(collectConfig []Collect) (err error) {
	tailObjMgr.lock.Lock()
	defer tailObjMgr.lock.Unlock()

	for _, value := range collectConfig {
		// 判断tailf运行状态，是否存在
		var isRunning = false
		for _, oldTailObj := range tailObjMgr.tailObjs {
			if value.LogPath == oldTailObj.collect.LogPath {
				isRunning = true
				break
			}
		}
		// 如果tailf任务不存在，创建新的任务
		if isRunning == false {
			createTask(value)
		}
	}

	// 更新tailf任务管理列表内容
	var tailObjs []*TailObj
	for _, oldTailObj := range tailObjMgr.tailObjs {
		oldTailObj.status = StatusDelete
		for _, newColl := range collectConfig {
			if newColl.LogPath == oldTailObj.collect.LogPath {
				oldTailObj.status = StatusNormal
				break
			}
		}
		if oldTailObj.status == StatusDelete {
			oldTailObj.exitChan <- 1
			continue
		}
		tailObjs = append(tailObjs, oldTailObj)
	}
	tailObjMgr.tailObjs = tailObjs
	return
}

// 从chan中获取一行数据
func GetOneLine() (msg *TextMsg) {
	msg = <-tailObjMgr.msgChan
	return
}

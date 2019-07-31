package tailf

import (
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

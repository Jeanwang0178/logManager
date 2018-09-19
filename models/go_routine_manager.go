package models

import (
	"fmt"
	"github.com/beego/bee/logger"
	"github.com/hpcloud/tail"
	"strconv"
	"strings"
	"time"
)

type GoRoutineManager struct {
	grchannelMap *GoroutineChannelMap
}

func NewGoRoutineManager() *GoRoutineManager {
	gm := &GoroutineChannelMap{}
	return &GoRoutineManager{grchannelMap: gm}
}

func (gm GoRoutineManager) StopLoopGoroutine(name string) error {
	stopChannel, ok := gm.grchannelMap.grchannels[name]
	if !ok {
		return fmt.Errorf("not found goroutine name :" + name)
	}
	gm.grchannelMap.grchannels[name].msg <- STOP + strconv.Itoa(int(stopChannel.gid))
	return nil
}

func (gm *GoRoutineManager) RegisterGoroutine(name string) {
	go func(n string) {
		//register channel
		err := gm.grchannelMap.register(n)
		if err != nil {
			beeLogger.Log.Errorf("grchannelMap register: %v", err)
			return
		}
	}(name)

}

func (gm *GoRoutineManager) NewLoopGoroutine(name string, tails *tail.Tail) {

	go func(this *GoRoutineManager, n string, tails tail.Tail) {
		//register channel
		err := this.grchannelMap.register(n)
		if err != nil {
			return
		}
		for {
			select {
			case info := <-this.grchannelMap.grchannels[name].msg:
				taskInfo := strings.Split(info, ":")
				signal, gid := taskInfo[0], taskInfo[1]
				if gid == strconv.Itoa(int(this.grchannelMap.grchannels[name].gid)) {
					if signal == "_STOP" {

						beeLogger.Log.Info(name + "：gid[" + gid + "] quit")
						this.grchannelMap.unregister(name)
						tails.Done()
						return
					} else {
						beeLogger.Log.Info("unknow signal")
					}
				}
			default:
				beeLogger.Log.Info("no signal")
			}
			sendMsg(tails)

		}
	}(gm, name, *tails)

}

func sendMsg(tails tail.Tail) (err error) {
	msg, ok := <-tails.Lines
	if !ok {
		beeLogger.Log.Info("tail file close reopen, filename:%s\n" + tails.Filename)
		time.Sleep(100 * time.Millisecond)
		return
	}
	err = SendToKafka(msg.Text, "topicLog")
	if err != nil {
		beeLogger.Log.Errorf("taild file error : %v ", err)
	}
	return nil
}

func (gm *GoRoutineManager) NewGoroutine(name string, fc interface{}, args ...interface{}) {
	go func(n string, fc interface{}, args ...interface{}) {
		//register channel
		err := gm.grchannelMap.register(n)
		if err != nil {
			beeLogger.Log.Errorf("grchannelMap register: %v", err)
			return
		}
		if len(args) > 1 {
			fc.(func(...interface{}))(args)
		} else if len(args) == 1 {
			fc.(func(interface{}))(args[0])
		} else {
			fc.(func())()
		}
		gm.grchannelMap.unregister(name)
	}(name, fc, args...)

}

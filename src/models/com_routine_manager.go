package models

import (
	"fmt"
	"github.com/hpcloud/tail"
	"logManager/src/common"
	"strconv"
	"strings"
	"time"
)

type GoRoutineManager struct {
	grchannelMap *GoroutineChannelMap
}

var (
	tailCount int32
)

func NewGoRoutineManager() *GoRoutineManager {
	gm := &GoroutineChannelMap{}
	return &GoRoutineManager{grchannelMap: gm}
}

func (gm GoRoutineManager) StopLoopGoroutine(name string) error {
	common.Logger.Info("StopLoopGoroutine ... :" + name)
	stopChannel, ok := gm.grchannelMap.grchannels[name]
	if !ok {
		return fmt.Errorf("not found goroutine name :" + name)
	}
	//stopChannel.tails.Done()
	line := tail.Line{"tailf file done ", time.Now(), nil}
	stopChannel.tails.Lines <- &line
	stopChannel.msg <- common.STOP + strconv.Itoa(int(stopChannel.gid))
	return nil
}

func (gm *GoRoutineManager) NewLoopGoroutine(name string, tails *tail.Tail) {

	go func(this *GoRoutineManager, n string, tails tail.Tail) {
		//register channel
		err := this.grchannelMap.register(n, tails)
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

						common.Logger.Info(name + "：gid[" + gid + "] quit")
						this.grchannelMap.unregister(name)
						//tails.Done()
						return
					} else {
						common.Logger.Info("unknow signal")
					}
				}
			default:
				//common.Logger.Info("no signal")
			}

			//发送KAFKA消息队列
			msg, ok := <-tails.Lines
			if !ok {
				common.Logger.Info("tail file close reopen, filename:%s\n" + tails.Filename)
				time.Sleep(100 * time.Millisecond)
				return
			}
			err = SendToKafka(msg.Text, common.TopicLog)
			if err != nil {
				common.Logger.Error("taild file error : %v ", err)
			}

		}
	}(gm, name, *tails)

}

func (gm *GoRoutineManager) TailfFiles(filePath string) {

	tails, err := tail.TailFile(filePath, tail.Config{
		ReOpen: true,
		Follow: true,
		//Location:&tail.SeekInfo{Offset:0,Whence:2},
		MustExist: false,
		Poll:      true,
	})

	if err != nil {
		common.Logger.Error("taild file error : %v ", err)
	}
	gm.NewLoopGoroutine(common.RoutineName, tails)
	return
}

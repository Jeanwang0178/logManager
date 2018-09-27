package models

import (
	"fmt"
	"github.com/gorilla/websocket"
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

func (gm GoRoutineManager) StopLoopGoroutine(name string, msgKey string) error {
	common.Logger.Info("StopLoopGoroutine ... :" + name)
	stopChannel, ok := gm.grchannelMap.grchannels[name]
	if !ok {
		return fmt.Errorf("not found goroutine name :" + name)
	}
	stopChannel.tails.Done()
	line := tail.Line{"tailf file done ", time.Now(), nil}
	stopChannel.tails.Lines <- &line
	stopChannel.msg <- common.STOP + strconv.Itoa(int(stopChannel.gid))
	return nil
}

func (gm *GoRoutineManager) NewLoopGoroutine(name string, tails *tail.Tail, showType string, msgKey string) {

	go func(this *GoRoutineManager, n string, tails tail.Tail, msgKey string) {
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
						tails.Cleanup()

						//dying := make(chan struct{})
						//tails.Dying()= dying
						tails.Done()
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
			time.Sleep(5 * time.Millisecond)
			if showType == common.ShowKafka {
				err = SendToKafka(msgKey, msg.Text)
				if err != nil {
					common.Logger.Error("taild file error : %v ", err)
				}
			}
			if showType == common.ShowTailf {
				// msg := Message{string(msg.Text)}
				msgt := msg.Text
				broadCast := BroadCastMap[msgKey]
				broadCast.msgChan <- msgt //广播发送至页面
			}

		}
	}(gm, name, *tails, msgKey)

}

func (gm *GoRoutineManager) TailfFiles(filePath string, showType string, socketConn websocket.Conn, msgKey string) {

	tails, err := tail.TailFile(filePath, tail.Config{
		ReOpen: true,
		Follow: true,
		//Location:&tail.SeekInfo{Offset:0,Whence:2},
		MustExist:   false,
		Poll:        true,
		MaxLineSize: 1024,
	})

	if err != nil {
		common.Logger.Error("taild file error : %v ", err)
	}

	gm.NewLoopGoroutine(common.RoutineKafka, tails, showType, msgKey)
	go HandlerMessage(msgKey, socketConn)
	return
}

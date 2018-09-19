package utils

import (
	"github.com/hpcloud/tail"
	"logManager/models"
	"time"
)

const (
	RoutineName = "sendLogFile"
)

var (
	tailCount int32
)

func TailfFiles() {

	gm := models.NewGoRoutineManager()

	fileName := "C:\\data\\logs\\sinochem-oms_20180607.log"
	tails, err := tail.TailFile(fileName, tail.Config{
		ReOpen: true,
		Follow: true,
		//Location:&tail.SeekInfo{Offset:0,Whence:2},
		MustExist: false,
		Poll:      true,
	})

	if err != nil {
		Logger.Error("taild file error : %v ", err)
	}

	gm.NewLoopGoroutine(RoutineName, sendMsg, tails)

	return
}

func sendMsg(sliceParam interface{}) (err error) {
	tails := sliceParam.(tail.Tail)
	msg, ok := <-tails.Lines
	if !ok {
		Logger.Info("tail file close reopen, filename:%s\n", tails.Filename)
		time.Sleep(100 * time.Millisecond)
		return
	}
	err = SendToKafka(msg.Text, TopicLog)
	if err != nil {
		Logger.Error("taild file error : %v ", err)
	}
	return nil
}

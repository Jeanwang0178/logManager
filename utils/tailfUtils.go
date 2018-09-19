package utils

import (
	"github.com/hpcloud/tail"
	"logManager/models"
)

const (
	RoutineName = "sendLogFile"
)

var (
	tailCount int32
)

func TailfFiles(gm *models.GoRoutineManager) {

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

	gm.NewLoopGoroutine(RoutineName, tails)

	return
}

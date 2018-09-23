package utils

import (
	"logManager/src/models"
)

var (
	tailCount int32
)

func TailfFiles(gm *models.GoRoutineManager) {

	/*fileName := "C:\\data\\logs\\sinochem-oms.log"
	tails, err := tail.TailFile(fileName, tail.Config{
		ReOpen: true,
		Follow: true,
		//Location:&tail.SeekInfo{Offset:0,Whence:2},
		MustExist: false,
		Poll:      true,
	})

	if err != nil {
		common.Logger.Error("taild file error : %v ", err)
	}*/
	//broadChan := make(chan models.Message)
	//gm.NewLoopGoroutine(common.RoutineKafka, tails, common.ShowKafka,models.BroadCast{broadChan})

	return
}

package services

import (
	"github.com/gorilla/websocket"
	"logManager/src/common"
	"logManager/src/models"
	"time"
)

/**
 * 1、tailf file文件  2、发送 kafka 3、页面建立webSocket连接 4、监听kafka消息队列，推送页面
 */
func LogFileServiceViewFile(webSocket *websocket.Conn) {

	//filePath1,_:= filepath.Abs("./")
	fileName := "logs/log_manager.log"

	gm := models.NewGoRoutineManager()
	go gm.TailfFiles(fileName)
	models.Clients[webSocket] = true

	for { //处理页面断开
		time.Sleep(time.Second * 1)
		var msg models.Message // Read in a new message as json and map it to a Message object
		err := webSocket.ReadJSON(&msg)
		common.Logger.Info("", err)
		if err != nil {
			common.Logger.Info("页面断开 ws.ReadJson error : %v ", err)
			delete(models.Clients, webSocket)
			defer webSocket.Close()
			err := gm.StopLoopGoroutine(common.RoutineName)
			common.Logger.Info("gm.StopLoopGoroutine failed : %v ", err)
			break
		} else {
			common.Logger.Info("接受从页面反馈回来的信息：", msg.Message)
		}
	}

}

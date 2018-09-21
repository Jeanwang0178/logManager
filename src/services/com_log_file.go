package services

import (
	"github.com/astaxie/beego/httplib"
	"github.com/gorilla/websocket"
	"github.com/satori/go.uuid"
	"logManager/src/common"
	"logManager/src/models"
	"logManager/src/utils"
	"strings"
	"time"
)

/**
 * 1、tailf file文件  2、发送 kafka 3、页面建立webSocket连接 4、监听kafka消息队列，推送页面
 */
func LogFileServiceViewFile(webSocket *websocket.Conn, filePath string) {

	//filePath1,_:= filepath.Abs("./")
	//filePath := "logs/log_manager.log"

	gm := models.NewGoRoutineManager()
	go gm.TailfFiles(filePath, common.ShowKafka)
	models.Clients[webSocket] = true

	for { //处理页面断开
		time.Sleep(time.Second * 1)
		var msg models.Message // Read in a new message as json and map it to a Message object
		err := webSocket.ReadJSON(&msg)
		common.Logger.Info("", err)
		if err != nil {
			common.Logger.Info("页面断开 ws.ReadJson ... : %v ", err)
			delete(models.Clients, webSocket)
			defer webSocket.Close()
			err := gm.StopLoopGoroutine(common.RoutineKafka)
			common.Logger.Info("gm.StopLoopGoroutine ... : %v ", err)
			break
		} else {
			common.Logger.Info("接受从页面反馈回来的信息：", msg.Message)
		}
	}

}

//远程调用接口
func LogFileServiceViewFile_Remote(webSocket *websocket.Conn, filePath string) {

	//filePath1,_:= filepath.Abs("./")
	//filePath := "logs/log_manager.log"

	//gm := models.NewGoRoutineManager()
	chanName := strings.Replace(uuid.Must(uuid.NewV4()).String(), "-", "", -1)
	remoteServiceStart(chanName, filePath)
	models.Clients[webSocket] = true

	for { //处理页面断开
		time.Sleep(time.Second * 1)
		var msg models.Message // Read in a new message as json and map it to a Message object
		err := webSocket.ReadJSON(&msg)
		common.Logger.Info("", err)
		if err != nil {
			common.Logger.Info("页面断开 ws.ReadJson ... : %v ", err)
			delete(models.Clients, webSocket)
			defer webSocket.Close()
			/*err := gm.StopLoopGoroutine(common.RoutineKafka)
			common.Logger.Info("gm.StopLoopGoroutine ... : %v ", err)*/
			remoteServiceStop(chanName)
			break
		} else {
			common.Logger.Info("接受从页面反馈回来的信息：", msg.Message)
		}
	}

}

func remoteServiceStop(chanName string) {
	vModel := models.ConfigRemote{}
	vModel.RemoteAddr = "http://192.168.3.151:9901/open/logFile/stopTail"
	vModel.Method = "POST"
	vModel.Header = "{}"
	vModel.Param = "{}"
	vModel.Body = "{\"chanName\":\"" + chanName + "\"}"

	request, err := utils.SendPost(vModel)
	if err != nil {
		common.Logger.Error(err.Error())
		return
	}
	var req = request.(*httplib.BeegoHTTPRequest)

	resp, err := req.Response()
	common.Logger.Info(resp.Status)

	strBody, err := req.String()

	if err != nil {
		common.Logger.Error("request faile %v ", err)
	}

	common.Logger.Info(strBody)
}

func remoteServiceStart(chanName string, filePath string) (err error) {
	vModel := models.ConfigRemote{}
	vModel.RemoteAddr = "http://192.168.3.151:9901/open/logFile/startTail"
	vModel.Method = "POST"
	vModel.Header = "{}"
	vModel.Param = "{}"
	vModel.Body = "{\"filePath\":\"" + filePath + "\",\"chanName\":\"" + chanName + "\"}"

	request, err := utils.SendPost(vModel)
	if err != nil {
		common.Logger.Error(err.Error())
		return err
	}
	var req = request.(*httplib.BeegoHTTPRequest)

	resp, err := req.Response()
	common.Logger.Info(resp.Status)

	strBody, err := req.String()

	if err != nil {
		common.Logger.Error("request faile %v ", err)
	}

	common.Logger.Info(strBody)
	return nil

}

// tailf 日志文件
func LogFileServiceTailfFile(webSocket *websocket.Conn, filePath string) {
	gm := models.NewGoRoutineManager()
	go gm.TailfFiles(filePath, common.ShowTailf)
	models.Clients[webSocket] = true

	for { //处理页面断开
		time.Sleep(time.Second * 1)
		var msg models.Message // Read in a new message as json and map it to a Message object
		err := webSocket.ReadJSON(&msg)
		common.Logger.Info("", err)
		if err != nil {
			common.Logger.Info("页面断开 ws.ReadJson ... : %v ", err)
			delete(models.Clients, webSocket)
			defer webSocket.Close()
			err := gm.StopLoopGoroutine(common.RoutineKafka)
			common.Logger.Info("gm.StopLoopGoroutine ... : %v ", err)
			break
		} else {
			common.Logger.Info("接受从页面反馈回来的信息：", msg.Message)
		}
	}
}

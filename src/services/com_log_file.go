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
func LogFileServiceViewFile(socketConn *websocket.Conn, filePath string) {

	gm := models.NewGoRoutineManager()
	msgKey := strings.Replace(uuid.Must(uuid.NewV4()).String(), "-", "", -1)

	models.Clients[socketConn] = models.NewBroadCast()
	models.BroadCastMap[msgKey] = models.Clients[socketConn] //key->broadCast

	go gm.TailfFiles(filePath, common.ShowKafka, *socketConn, msgKey)

	for { //处理页面断开
		time.Sleep(time.Second * 1)
		var msg models.Message // Read in a new message as json and map it to a Message object

		err := socketConn.ReadJSON(&msg)
		common.Logger.Info("", err)
		if err != nil {
			common.Logger.Info("页面断开 ws.ReadJson ... : %v ", err)
			delete(models.Clients, socketConn)
			defer socketConn.Close()
			err := gm.StopLoopGoroutine(common.RoutineKafka, msgKey)
			common.Logger.Info("gm.StopLoopGoroutine ... : %v ", err)
			break
		} else {
			common.Logger.Info("接受从页面反馈回来的信息：", msg.Message)
		}
	}

}

//远程调用接口
func LogFileServiceViewFile_remote(socketConn *websocket.Conn, remoteAddr string, filePath string) {

	chanName := strings.Replace(uuid.Must(uuid.NewV4()).String(), "-", "", -1)
	msgKey := strings.Replace(uuid.Must(uuid.NewV4()).String(), "-", "", -1)

	models.Clients[socketConn] = models.NewBroadCast()
	models.BroadCastMap[msgKey] = models.Clients[socketConn] //key->broadCast
	go models.HandlerMessage(msgKey, *socketConn)

	vModel := models.ConfigRemote{}
	vModel.RemoteAddr = remoteAddr + "/open/logFile/startTail"
	vModel.Method = "POST"
	vModel.Header = "{}"
	vModel.Param = "{}"
	vModel.Body = "{\"filePath\":\"" + filePath + "\",\"chanName\":\"" + chanName + "\",\"msgKey\":\"" + msgKey + "\"}"

	remoteServiceStart(vModel)

	for { //处理页面断开
		time.Sleep(time.Second * 1)
		var msg models.Message // Read in a new message as json and map it to a Message object
		err := socketConn.ReadJSON(&msg)
		common.Logger.Info("", err)
		if err != nil {
			common.Logger.Info("页面断开 ws.ReadJson ... : %v ", err)
			delete(models.Clients, socketConn)
			defer socketConn.Close()
			remoteServiceStop(chanName, msgKey)

			break
		} else {
			common.Logger.Info("接受从页面反馈回来的信息：", msg.Message)
		}
	}

}

//远程调用KAFKA 查看文件
func remoteServiceStart(vModel models.ConfigRemote) (err error) {

	request, err := utils.SendPost(vModel)
	if err != nil {
		common.Logger.Error(err.Error())
		return err
	}
	var req = request.(*httplib.BeegoHTTPRequest)

	resp, err := req.Response()
	if resp != nil {
		defer resp.Body.Close()
	}
	common.Logger.Info(resp.Status)

	strBody, err := req.String()

	if err != nil {
		common.Logger.Error("request faile %v ", err)
	}

	common.Logger.Info(strBody)
	return nil

}

//停止远程进程
func remoteServiceStop(chanName string, msgKey string) {
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
	if resp != nil {
		defer resp.Body.Close()
	}
	common.Logger.Info(resp.Status)

	strBody, err := req.String()

	if err != nil {
		common.Logger.Error("request faile %v ", err)
	}
	delete(models.BroadCastMap, msgKey) //删除key
	common.Logger.Debug(strBody)
}

// tailf 日志文件
func LogFileServiceTailfFile(socketConn *websocket.Conn, filePath string) {
	gm := models.NewGoRoutineManager()

	models.Clients[socketConn] = models.NewBroadCast()
	msgKey, _ := uuid.NewV4()
	models.BroadCastMap[msgKey.String()] = models.Clients[socketConn] //key->broadCast

	go gm.TailfFiles(filePath, common.ShowTailf, *socketConn, msgKey.String())
	//models.Clients[webSocket] = true

	for { //处理页面断开
		time.Sleep(time.Second * 1)
		var msg models.Message // Read in a new message as json and map it to a Message object
		err := socketConn.ReadJSON(&msg)
		common.Logger.Info("", err)
		if err != nil {
			common.Logger.Info("页面断开 ws.ReadJson ... : %v ", err)
			delete(models.Clients, socketConn)
			defer socketConn.Close()
			err := gm.StopLoopGoroutine(common.RoutineKafka, msgKey.String())
			common.Logger.Info("gm.StopLoopGoroutine ... : %v ", err)
			break
		} else {
			common.Logger.Info("接受从页面反馈回来的信息：", msg.Message)
		}
	}
}

//停止远程进程
func RemoteServiceListFile(remoteAddr string, foldPath string) (strBody string, err error) {
	vModel := models.ConfigRemote{}
	vModel.RemoteAddr = remoteAddr + "/open/logFile/listFile"
	vModel.Method = "POST"
	vModel.Header = "{}"
	vModel.Param = "{}"
	vModel.Body = "{\"foldPath\":\"" + foldPath + "\"}"

	request, err := utils.SendPost(vModel)
	if err != nil {
		common.Logger.Error(err.Error())
		return
	}
	var req = request.(*httplib.BeegoHTTPRequest)

	resp, err := req.Response()
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return "", err
	}
	common.Logger.Info(resp.Status)

	strBody, err = req.String()

	if err != nil {
		common.Logger.Error("request faile %v ", err)
		return "", err
	}
	return strBody, nil
}

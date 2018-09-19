package controllers

import (
	"github.com/gorilla/websocket"
	"logManager/models"
	"logManager/utils"
	"time"
)

type LogFileController struct {
	BaseController
}

var upgrader = websocket.Upgrader{}

// @router /view [get]
func (ctl *LogFileController) View() {

	response := make(map[string]interface{})

	response["code"] = utils.SuccessCode
	response["msg"] = utils.SuccessMsg

	ctl.display()

}

// @router /viewLog [get]
func (ctl *LogFileController) ViewLog() {

	gm := models.NewGoRoutineManager()

	go utils.TailfFiles()

	ws, err := upgrader.Upgrade(ctl.Ctx.ResponseWriter, ctl.Ctx.Request, nil)
	if err != nil {
		utils.Logger.Error("get connection failed：%v ", err)
	}
	utils.Clients[ws] = true

	for {
		//发送广播至页面
		time.Sleep(time.Second * 1)

		var msg models.Message // Read in a new message as json and map it to a Message object
		err := ws.ReadJSON(&msg)
		utils.Logger.Info("", err)
		if err != nil {
			utils.Logger.Info("页面断开 ws.ReadJson error : %v ", err)
			delete(utils.Clients, ws)
			defer ws.Close()
			gm.StopLoopGoroutine(utils.RoutineName)
			break
		} else {
			utils.Logger.Info("接受从页面反馈回来的信息：", msg.Message)
		}

	}
	ctl.TplName = "logfile/view.html"

}

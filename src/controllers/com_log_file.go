package controllers

import (
	"github.com/gorilla/websocket"
	"logManager/src/common"
	"logManager/src/services"
	"logManager/src/utils"
)

type LogFileController struct {
	BaseController
}

var upgrader = websocket.Upgrader{}

// @router /view [get]
func (ctl *LogFileController) View() {

	response := make(map[string]interface{})

	fileNames, err := utils.ListFile("./logs")

	if err != nil {
		common.Logger.Error("utils listFile error : %v ", err)
	}

	response["code"] = utils.SuccessCode
	response["msg"] = utils.SuccessMsg
	response["data"] = fileNames

	ctl.Data["result"] = response
	ctl.display()

}

// @router /viewLog [get]
func (ctl *LogFileController) ViewLog() {
	filePath := ctl.GetString("filePath")
	webSocket, err := upgrader.Upgrade(ctl.Ctx.ResponseWriter, ctl.Ctx.Request, nil)

	if err != nil {
		common.Logger.Error("get connection failed：%v ", err)
	} else {
		services.LogFileServiceViewFile(webSocket, filePath)
	}

	ctl.TplName = "logfile/view.html"

}

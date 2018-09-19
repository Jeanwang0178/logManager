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

	response["code"] = utils.SuccessCode
	response["msg"] = utils.SuccessMsg

	ctl.display()

}

// @router /viewLog [get]
func (ctl *LogFileController) ViewLog() {

	webSocket, err := upgrader.Upgrade(ctl.Ctx.ResponseWriter, ctl.Ctx.Request, nil)

	if err != nil {
		common.Logger.Error("get connection failed：%v ", err)
	} else {
		services.LogFileServiceViewFile(webSocket)
	}

	ctl.TplName = "logfile/view.html"

}

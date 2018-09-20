package controllers

import (
	"github.com/gorilla/websocket"
	"logManager/src/common"
	"logManager/src/services"
	"logManager/src/utils"
	"strings"
)

type LogFileController struct {
	BaseController
}

var upgrader = websocket.Upgrader{}

// @router /view [get]
func (ctl *LogFileController) View() {

	response := make(map[string]interface{})
	var fileNames = []string{}
	var err error
	foldPath := strings.TrimSpace(ctl.GetString("foldPath"))
	if foldPath == "nil" || foldPath == "" {
		foldPath = "C:/data/logs"
	}

	fileNames, err = utils.ListFile(foldPath)
	if err != nil {
		common.Logger.Error("utils listFile error : %v ", err)
	}

	response["code"] = utils.SuccessCode
	response["msg"] = utils.SuccessMsg
	response["data"] = fileNames
	response["param"] = foldPath
	ctl.Data["result"] = response
	ctl.display()

}

// @router /viewLog [get]
func (ctl *LogFileController) ViewLog() {
	filePath := strings.TrimSpace(ctl.GetString("filePath"))
	webSocket, err := upgrader.Upgrade(ctl.Ctx.ResponseWriter, ctl.Ctx.Request, nil)

	if err != nil {
		common.Logger.Error("get connection failed：%v ", err)
	} else {
		services.LogFileServiceViewFile(webSocket, filePath)
	}

	ctl.TplName = "logfile/view.html"

}

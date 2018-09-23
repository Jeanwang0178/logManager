package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/gorilla/websocket"
	"logManager/src/common"
	"logManager/src/models"
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

	var query = make(map[string]string)
	var fileNames = []string{}
	var err error
	foldPath := strings.TrimSpace(ctl.GetString("foldPath"))
	queryType := strings.TrimSpace(ctl.GetString("queryType"))
	if foldPath == "nil" || foldPath == "" {
		foldPath = "C:/data/logs"
	}

	if queryType == "nil" || queryType == "" {
		queryType = "kafka"
	}

	remoteTail := beego.AppConfig.String("tailf.kafka.type")

	fileNames, err = utils.ListFile(foldPath)
	if err != nil {
		common.Logger.Error("utils listFile error : %v ", err)
	}

	query["foldPath"] = foldPath
	query["queryType"] = queryType

	response["code"] = utils.SuccessCode
	response["msg"] = utils.SuccessMsg
	response["data"] = fileNames
	response["remoteTail"] = remoteTail
	response["param"] = query
	ctl.Data["result"] = response
	ctl.display()

}

// @router /viewLog [get]
func (ctl *LogFileController) ViewLog() {
	filePath := strings.TrimSpace(ctl.GetString("filePath"))
	remoteAddr := strings.TrimSpace(ctl.GetString("remoteAddr"))
	webSocket, err := upgrader.Upgrade(ctl.Ctx.ResponseWriter, ctl.Ctx.Request, nil)

	if err != nil {
		common.Logger.Error("get connection failed：%v ", err)
	} else {
		remoteTail := beego.AppConfig.String("tailf.kafka.type")
		if "remote" != remoteTail {
			services.LogFileServiceViewFile(webSocket, filePath)
		} else {
			services.LogFileServiceViewFile_remote(webSocket, remoteAddr, filePath)
		}
	}

	ctl.TplName = "logfile/view.html"

}

// @router /tailfLog [get]
func (ctl *LogFileController) TailfLog() {
	filePath := strings.TrimSpace(ctl.GetString("filePath"))
	webSocket, err := upgrader.Upgrade(ctl.Ctx.ResponseWriter, ctl.Ctx.Request, nil)
	if err != nil {
		common.Logger.Error("get connection failed：%v ", err)
	} else {
		services.LogFileServiceTailfFile(webSocket, filePath)
	}

	ctl.TplName = "logfile/view.html"
}

// @router /listRemoteFile [post]
func (ctl *LogFileController) ListRemoteFile() {

	response := make(map[string]interface{})

	var err error
	remoteAddr := strings.TrimSpace(ctl.GetString("remoteAddr"))
	foldPath := strings.TrimSpace(ctl.GetString("foldPath"))
	if foldPath == "nil" || foldPath == "" {
		foldPath = "C:/data/logs"
	}

	if remoteAddr == "nil" || remoteAddr == "" {
		response["code"] = utils.FailedCode
		response["msg"] = "缺少远程调用地址"
		return
	}

	strBody, err := services.RemoteServiceListFile(remoteAddr, foldPath)
	resData := models.ResponseData{}
	json.Unmarshal([]byte(strBody), &resData)

	if err != nil {
		common.Logger.Error("utils listFile error : %v ", err)
		response["code"] = utils.FailedCode
		response["msg"] = err.Error()
	} else {
		common.Logger.Info(strBody)
		response["code"] = utils.SuccessCode
		response["msg"] = utils.SuccessMsg
		response["data"] = resData.Data
	}

	ctl.Data["json"] = response
	ctl.ServeJSON()

}

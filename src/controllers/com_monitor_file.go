package controllers

import (
	"encoding/json"
	"logManager/src/common"
	"logManager/src/models"
	"logManager/src/services"
	"logManager/src/utils"
	"strings"
)

type MonitorController struct {
	BaseController
}

// @router /view [get]
func (ctl *MonitorController) View() {

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
		queryType = "local"
	}

	fileNames, err = utils.ListFile(foldPath)
	if err != nil {
		common.Logger.Error("utils listFile error : %v ", err)
	}

	query["foldPath"] = foldPath
	query["queryType"] = queryType
	query["lineNum"] = "0"
	query["preLineNum"] = "0"

	response["code"] = utils.SuccessCode
	response["msg"] = utils.SuccessMsg
	response["data"] = fileNames
	response["param"] = query
	ctl.Data["result"] = response
	ctl.display()

}

// @router /queryContent [post]
func (ctl *MonitorController) QueryContent() {

	response := make(map[string]interface{})

	resData := models.RequestFileParam{}
	remoteAddr := strings.TrimSpace(ctl.GetString("remoteAddr"))
	filePath := strings.TrimSpace(ctl.GetString("filePath"))
	content := strings.TrimSpace(ctl.GetString("content"))
	nextLineNum, err := ctl.GetInt64("nextLineNum")
	preLineNum, err := ctl.GetInt64("preLineNum")
	queryType := strings.TrimSpace(ctl.GetString("queryType"))
	operType := strings.TrimSpace(ctl.GetString("operType"))
	position, err := ctl.GetInt("position")

	if err != nil {
		response["code"] = utils.FailedCode
		response["msg"] = err.Error()
	} else {

		resData.RemoteAddr = remoteAddr
		resData.FilePath = filePath
		resData.Content = content
		resData.NextLineNum = nextLineNum
		resData.PreLineNum = preLineNum
		resData.QueryType = queryType
		resData.OperType = operType
		resData.Position = position
		data, preOffset, nextOffset, position, err := services.MonitorFileServiceQueryContent(resData)
		if err != nil && !strings.Contains(err.Error(), "1001") {
			common.Logger.Error("utils listFile error : %v ", err)
			response["code"] = utils.FailedCode
			response["msg"] = err.Error()
		} else {
			response["code"] = utils.SuccessCode
			response["msg"] = utils.SuccessMsg
			response["data"] = data
			if err != nil {
				response["err"] = string(err.Error())
			}
		}
		response["nextOffset"] = nextOffset
		response["preOffset"] = preOffset
		response["position"] = position
	}
	ctl.Data["json"] = response
	ctl.ServeJSON()
}

// @router /listRemoteFile [post]
func (ctl *MonitorController) ListRemoteFile() {

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

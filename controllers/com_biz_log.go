package controllers

import (
	"ERP/utils"
	"errors"
	"github.com/astaxie/beego"
	"logManager/services"
	logger "logManager/utils"
	"strings"
	"webcron/app/libs"
)

type BizLogController struct {
	BaseController
}

// @router /list [post,get]
func (ctl *BizLogController) LogList() {

	logger.Logger.Debug("log manager list ")

	response := make(map[string]interface{})

	page, _ := ctl.GetInt("page")
	ctl.pageSize, _ = ctl.GetInt("pageSize")

	moduleName := ctl.Input().Get("moduleName")
	className := ctl.Input().Get("className")
	methodName := ctl.Input().Get("methodName")
	status := ctl.Input().Get("status")

	logger.Logger.Info(methodName, className, moduleName, status)

	if ctl.pageSize == 0 {
		ctl.pageSize = 10
	}

	var sortby = []string{"create_time"}
	var order = []string{"desc"}
	var query = make(map[string]string)
	var limit int64 = 10
	var offset = (int64)((page - 1) * ctl.pageSize)

	if moduleName != "" {
		query["moduleName"] = moduleName
	}
	if className != "" {
		query["className"] = className
	}
	if methodName != "" {
		query["methodName"] = methodName
	}
	if status != "" && status != "-1" {
		query["status"] = status
	}
	// query: k:v,k:v
	if v := ctl.GetString("query"); v != "" {
		for _, cond := range strings.Split(v, ",") {
			kv := strings.SplitN(cond, ":", 2)
			if len(kv) != 2 {
				err := errors.New("Error: invalid query key/value pair")
				response["code"] = utils.FailedCode
				response["msg"] = utils.FailedMsg
				response["err"] = err
				ctl.Data["json"] = response
				ctl.ServeJSON()
				return
			}
			k, v := kv[0], kv[1]
			query[k] = v
		}
	}

	logList, count, err := services.BizLogServiceGetList(query, []string{}, sortby, order, offset, limit)

	if err != nil {
		response["code"] = utils.FailedCode
		response["msg"] = utils.FailedMsg
		response["err"] = err
	} else {
		response["code"] = utils.SuccessCode
		response["msg"] = utils.SuccessMsg
		response["data"] = logList
	}

	query["status"] = status

	ctl.Data["param"] = query
	ctl.Data["result"] = response

	ctl.Data["pageBar"] = libs.NewPager(page, int(count), ctl.pageSize, beego.URLFor("BizLogController.LogList"), true).ToString()

	ctl.display()

}

// @router /findById [get]
func (ctl *BizLogController) LogView() {

	id := ctl.GetString("id")
	logger.Logger.Debug("log manager list ", id)
	response := make(map[string]interface{})

	bizLog, err := services.BizLogServiceGetById(id)

	if err != nil {
		response["code"] = utils.FailedCode
		response["msg"] = utils.FailedMsg
		response["err"] = err
	} else {
		response["code"] = utils.SuccessCode
		response["msg"] = utils.SuccessMsg
		response["data"] = bizLog
	}

	ctl.Data["result"] = response
	ctl.display()

}

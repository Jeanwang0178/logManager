package controllers

import (
	"errors"
	"fmt"
	"logManager/services"
	"logManager/utils"
)

type ManagerController struct {
	BaseController
}

// @router /query [get]
func (ctl *ManagerController) QueryList() {

	utils.Logger.Debug("log manager list ")

	response := make(map[string]interface{})

	ctl.pageSize, _ = ctl.GetInt("pageSize")

	aliasName := ctl.GetString("aliasName")
	tableName := ctl.Input().Get("tableName")
	utils.Logger.Info("query param", aliasName, tableName)

	var query = make(map[string]string)

	if aliasName != "" {
		query["aliasName"] = aliasName
	}
	if tableName != "" {
		query["tableName"] = tableName
	}

	response["code"] = utils.SuccessCode
	response["msg"] = utils.SuccessMsg

	ctl.Data["param"] = query
	ctl.Data["result"] = response

	ctl.display()

}

// @router /querylist [post]
func (ctl *ManagerController) QueryDataList() {

	utils.Logger.Debug("log manager list ")

	response := make(map[string]interface{})

	page, _ := ctl.GetInt("page")
	ctl.pageSize, _ = ctl.GetInt("pageSize")

	aliasName := ctl.GetString("aliasName")
	tableName := ctl.GetString("tableName")
	utils.Logger.Info("query param", aliasName, tableName)

	var sortby = []string{"field_sort"}
	var order = []string{"desc"}
	var query = make(map[string]string)
	var titleMap = make(map[string]string)
	var limit int64 = 10
	var offset = (int64)((page - 1) * ctl.pageSize)

	if aliasName != "" {
		query["aliasName"] = aliasName
	} else {
		response["code"] = utils.FailedCode
		response["msg"] = utils.FailedMsg
		response["err"] = errors.New("请选择数据库")
	}
	if tableName != "" {
		query["tableName"] = tableName
	} else {
		response["code"] = utils.FailedCode
		response["msg"] = utils.FailedMsg
		response["err"] = errors.New("请输入表名称")
	}

	response["param"] = query

	mappingList, titleMap, sortFields, count, err := services.ManagerServiceGetLogList(query, sortby, order, offset, limit)
	response["titles"] = titleMap
	fmt.Println(count)
	if err != nil {
		response["code"] = utils.FailedCode
		response["msg"] = err.Error()
	} else {
		response["code"] = utils.SuccessCode
		response["msg"] = utils.SuccessMsg
		response["data"] = mappingList
		response["sortFields"] = sortFields
	}

	ctl.Data["json"] = response
	ctl.ServeJSON()

}

// @router /findById [get]
func (ctl *ManagerController) View() {

	id := ctl.GetString("id")
	utils.Logger.Debug("log manager list ", id)
	response := make(map[string]interface{})

	bizLog, err := services.BizLogServiceGetById(id)

	if err != nil {
		response["code"] = utils.FailedCode
		response["msg"] = err.Error()
	} else {
		response["code"] = utils.SuccessCode
		response["msg"] = utils.SuccessMsg
		response["data"] = bizLog
	}

	ctl.Data["result"] = response
	ctl.display()

}

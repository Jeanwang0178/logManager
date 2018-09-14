package controllers

import (
	"errors"
	"github.com/astaxie/beego"
	"logManager/services"
	"logManager/utils"
	"webcron/app/libs"
)

type ManagerController struct {
	BaseController
}

// @router /list [get]
func (ctl *ManagerController) List() {

	utils.Logger.Debug("log manager list ")

	response := make(map[string]interface{})
	var query = make(map[string]string)

	ctl.pageSize, _ = ctl.GetInt("pageSize")

	aliasName := ctl.GetString("aliasName")
	tableName := ctl.Input().Get("tableName")
	utils.Logger.Info("query param", aliasName, tableName)

	query["aliasName"] = aliasName
	query["tableName"] = tableName

	aliasNames := make([]interface{}, 0)
	err := utils.GetCache(utils.AliasName, &aliasNames)
	if err != nil {
		response["code"] = utils.FailedCode
		response["msg"] = err.Error()
	} else {
		response["code"] = utils.SuccessCode
		response["msg"] = utils.SuccessMsg
	}

	ctl.Data["aliasNames"] = aliasNames

	ctl.Data["param"] = query
	ctl.Data["result"] = response

	ctl.display()

}

// @router /dataList [get,post]
func (ctl *ManagerController) DataList() {

	utils.Logger.Debug("log manager list ")

	response := make(map[string]interface{})

	page, _ := ctl.GetInt("page")
	ctl.pageSize, _ = ctl.GetInt("pageSize")
	if page == 0 {
		page = 1
	}
	if ctl.pageSize == 0 {
		ctl.pageSize = 20
	}

	aliasName := ctl.GetString("aliasName")
	tableName := ctl.GetString("tableName")
	utils.Logger.Info("query param", aliasName, tableName)

	var query = make(map[string]string)
	var titleMap = make(map[string]string)
	var limit int64 = 20
	var offset = (int64)((page - 1) * ctl.pageSize)
	query["aliasName"] = aliasName
	query["tableName"] = tableName
	if aliasName == "" {
		response["code"] = utils.FailedCode
		response["msg"] = errors.New("请选择数据库")
	}
	if tableName == "" {
		response["code"] = utils.FailedCode
		response["msg"] = errors.New("请输入表名称")
	}

	response["param"] = query

	aliasNames := make([]interface{}, 0)
	err := utils.GetCache(utils.AliasName, &aliasNames)
	if err != nil {
		utils.Logger.Error("utils.GetCache failed , key || %s", utils.AliasName)
	}

	mappingList, titleMap, fieldsSort, count, err := services.ManagerServiceGetDataList(query, offset, limit)
	response["titleMap"] = titleMap

	pageBar := libs.NewPager(page, int(count), ctl.pageSize, beego.URLFor("ManagerController.DataList", "aliasName", aliasName, "tableName", tableName), true).ToString()

	if err != nil {
		response["code"] = utils.FailedCode
		response["msg"] = err.Error()
	} else {
		response["code"] = utils.SuccessCode
		response["msg"] = utils.SuccessMsg
		response["data"] = mappingList
		response["pageBar"] = beego.Str2html(pageBar)
		response["fieldsSort"] = fieldsSort
	}

	ctl.Data["aliasNames"] = aliasNames

	ctl.Data["json"] = response
	ctl.ServeJSON()

}

// @router /view [get]
func (ctl *ManagerController) View() {

	id := ctl.GetString("id")
	aliasName := ctl.GetString("aliasName")
	tableName := ctl.GetString("tableName")

	utils.Logger.Debug("log manager list ", id)
	response := make(map[string]interface{})

	var query = make(map[string]string)

	query["id"] = id
	query["aliasName"] = aliasName
	query["tableName"] = tableName

	if id == "" {
		response["code"] = utils.FailedCode
		response["msg"] = errors.New("缺少参数ID")
	}
	if aliasName == "" {
		response["code"] = utils.FailedCode
		response["msg"] = errors.New("请选择数据库")
	}
	if tableName == "" {
		response["code"] = utils.FailedCode
		response["msg"] = errors.New("请输入表名称")
	}
	dataMap, titleMap, fieldsSort, err := services.ManagerServiceGetDataById(query)

	if err != nil {
		response["code"] = utils.FailedCode
		response["msg"] = err.Error()
	} else {
		response["code"] = utils.SuccessCode
		response["msg"] = utils.SuccessMsg
		response["dataMap"] = dataMap
		response["titleMap"] = titleMap
		response["fieldsSort"] = fieldsSort
	}

	ctl.Data["result"] = response
	ctl.display()

}

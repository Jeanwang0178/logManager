package controllers

import (
	"errors"
	"github.com/astaxie/beego"
	"logManager/services"
	utils "logManager/utils"
	"strings"
	"webcron/app/libs"
)

//查看common日志
type BizLogController struct {
	BaseController
}

// @Title Post,Get
// @Description 获取日志列表
// @Param page query int true "页码"
// @Param pageSize query  int true "分页大小"
// @Param moduleName query  string flase "模块名称"
// @Param className query  string flase "类名称"
// @Param methodName query  string flase "方法名称"
// @Param status query   string flase "状态"
// @Success 200 {object} models.BizLog "0k"
// @Failure 403 : other err
// @router /list [post,get]
func (ctl *BizLogController) List() {

	utils.Logger.Debug("log manager list ")

	response := make(map[string]interface{})

	page, _ := ctl.GetInt("page")
	ctl.pageSize, _ = ctl.GetInt("pageSize")

	moduleName := ctl.Input().Get("moduleName")
	className := ctl.Input().Get("className")
	methodName := ctl.Input().Get("methodName")
	status := ctl.Input().Get("status")

	if ctl.pageSize == 0 {
		ctl.pageSize = 20
	}

	var sortby = []string{"create_time"}
	var order = []string{"desc"}
	var query = make(map[string]string)
	var limit int64 = 20
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
				response["msg"] = err.Error()
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
		response["msg"] = err.Error()
	} else {
		response["code"] = utils.SuccessCode
		response["msg"] = utils.SuccessMsg
		response["data"] = logList
	}

	query["status"] = status

	ctl.Data["param"] = query
	ctl.Data["result"] = response

	ctl.Data["pageBar"] = libs.NewPager(page, int(count), ctl.pageSize, beego.URLFor("BizLogController.LogList", "status", status, "moduleName", moduleName, "className", className, "methodName", methodName), true).ToString()

	ctl.display()

}

// @router /view [get]
func (ctl *BizLogController) View() {

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

// @router /edit [get]
func (ctl *BizLogController) Edit() {

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

// @Title Post
// @Description 获取日志列表
// @Param	body body 	models.BizLog	true		"日志详情"
// @Success 201 {int} models.BizLog
// @Success 200 {object} models.BizLog "0k"
// @Failure 403 : other err
// @router /save [post]
func (ctl *BizLogController) Save() {

	id := ctl.GetString("Id")

	response := make(map[string]interface{})

	bizLog, err := services.BizLogServiceGetById(id)
	if err != nil {
		ctl.showMsg(err.Error())
	}

	if ctl.isPost() {
		bizLog.ModuleName = ctl.GetString("ModuleName")
		bizLog.UserId = ctl.GetString("UserId")
		ctime := ctl.GetString("CreateTime")
		utils.Logger.Info(ctime)
		createT, _ := beego.DateParse(ctl.GetString("CreateTime"), "Y-m-d H:i:s")
		bizLog.CreateTime = createT
		bizLog.Ip = ctl.GetString("Ip")
		bizLog.Status, _ = ctl.GetInt("Status")
		bizLog.ModuleName = ctl.GetString("ModuleName")
		bizLog.ClassName = ctl.Input().Get("ClassName")
		bizLog.MethodName = ctl.Input().Get("MethodName")
		bizLog.Params = ctl.Input().Get("Params")
		bizLog.Commemts = ctl.Input().Get("Commemts")
		err = services.BizLogServiceUpdate(bizLog)
	}

	if err != nil {
		response["code"] = utils.FailedCode
		response["msg"] = err.Error()
	} else {
		response["code"] = utils.SuccessCode
		response["msg"] = utils.SuccessMsg
		response["data"] = bizLog
	}
	ctl.Data["json"] = response
	ctl.ServeJSON()
	return
}

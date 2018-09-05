package controllers

import (
	"errors"
	"logManager/services"
	"logManager/utils"
	"strings"
)

type ManagerController struct {
	BaseController
}

// @router /query [get]
func (ctl *ManagerController) QueryList() {

	utils.Logger.Debug("log manager list ")

	response := make(map[string]interface{})

	page, _ := ctl.GetInt("page")
	ctl.pageSize, _ = ctl.GetInt("pageSize")

	aliasName := ctl.GetString("aliasName")
	tableName := ctl.Input().Get("tableName")
	utils.Logger.Info("query param", aliasName, tableName)

	var sortby = []string{"create_time"}
	var order = []string{"desc"}
	var query = make(map[string]string)
	var limit int64 = 10
	var offset = (int64)((page - 1) * ctl.pageSize)

	if aliasName != "" {
		query["aliasName"] = aliasName
	}
	if tableName != "" {
		query["tableName"] = tableName
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

	mappingList, err := services.MappingServiceGetList(query, []string{}, sortby, order, offset, limit)

	if err != nil {
		response["code"] = utils.FailedCode
		response["msg"] = utils.FailedMsg
		response["err"] = err
	} else {
		response["code"] = utils.SuccessCode
		response["msg"] = utils.SuccessMsg
		response["data"] = mappingList
	}

	ctl.Data["param"] = query
	ctl.Data["result"] = response

	ctl.display()

}

// @router /findById [get]
func (ctl *ManagerController) View() {

	id := ctl.GetString("id")
	utils.Logger.Debug("log manager list ", id)
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

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

// @router /list [get]
func (ctl *BizLogController) LogList() {

	logger.Logger.Debug("log manager list ")

	logId, _ := ctl.GetInt("id")
	page, _ := ctl.GetInt("page")
	response := make(map[string]interface{})

	var fields []string
	var sortby []string
	var order []string
	var query = make(map[string]string)
	var limit int64 = 10
	var offset int64

	// fields: col1,col2,entity.col3
	if v := ctl.GetString("fields"); v != "" {
		fields = strings.Split(v, ",")
	}
	// limit: 10 (default is 10)
	if v, err := ctl.GetInt64("limit"); err == nil {
		limit = v
	}
	// offset: 0 (default is 0)
	if v, err := ctl.GetInt64("offset"); err == nil {
		offset = v
	}
	// sortby: col1,col2
	if v := ctl.GetString("sortby"); v != "" {
		sortby = strings.Split(v, ",")
	}
	// order: desc,asc
	if v := ctl.GetString("order"); v != "" {
		order = strings.Split(v, ",")
	}
	// query: k:v,k:v
	if v := ctl.GetString("query"); v != "" {
		for _, cond := range strings.Split(v, ",") {
			kv := strings.SplitN(cond, ":", 2)
			if len(kv) != 2 {
				err := errors.New("Error: invalid query key/value pair")
				ctl.ServeJSON()
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

	logList, err := services.BizLogServiceGetList(query, fields, sortby, order, offset, limit)
	count := 100
	if err != nil {
		response["code"] = utils.FailedCode
		response["msg"] = utils.FailedMsg
		response["err"] = err
	} else {
		response["code"] = utils.SuccessCode
		response["msg"] = utils.SuccessMsg
		response["data"] = logList
	}

	ctl.Data["json"] = response

	ctl.Data["pageBar"] = libs.NewPager(page, int(count), ctl.pageSize, beego.URLFor("BizLogController.LogList", "id", logId), true).ToString()

	ctl.display()

	/*ctl.Data["json"] = response
	ctl.ServeJSON()*/

}

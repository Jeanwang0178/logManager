package controllers

import (
	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/httplib"
	"logManager/src/common"
	"logManager/src/models"
	"logManager/src/services"
	"logManager/src/utils"
	"regexp"
	"webcron/app/libs"
)

// RemoteController operations for ConfigRemote
type RemoteController struct {
	BaseController
}

// @router /edit [get]
func (ctl *RemoteController) Edit() {
	response := make(map[string]interface{})

	idStr := strings.TrimSpace(ctl.GetString("Id"))
	var vmodel = &models.ConfigRemote{}
	var err error
	if idStr != "" {
		vmodel, err = services.ConfigRemoteServiceGetById(idStr)
		if err != nil {
			response["code"] = utils.FailedCode
			response["msg"] = err.Error()
		} else {
			response["code"] = utils.SuccessCode
			response["msg"] = utils.SuccessMsg

		}
	}
	response["data"] = vmodel
	ctl.Data["result"] = response
	ctl.display()
}

// @router /save [post]
func (ctl *RemoteController) Save() {
	response := make(map[string]interface{})

	sId := strings.TrimSpace(ctl.GetString("Id"))
	sRemoteAddr := strings.TrimSpace(ctl.GetString("RemoteAddr"))
	sMethod := strings.TrimSpace(ctl.GetString("Method"))
	sHeader := ctl.GetString("Header")
	sParam := ctl.GetString("Param")
	sBody := ctl.GetString("Body")
	operType := ctl.GetString("operType")

	reg := regexp.MustCompile(`\n|\r`)

	sHeader = strings.TrimSpace(reg.ReplaceAllString(sHeader, ""))
	sParam = strings.TrimSpace(reg.ReplaceAllString(sParam, ""))
	sBody = strings.TrimSpace(reg.ReplaceAllString(sBody, ""))

	vModel := models.ConfigRemote{}
	vModel.Id = sId
	vModel.RemoteAddr = sRemoteAddr
	vModel.Method = sMethod
	vModel.Header = sHeader
	vModel.Param = sParam
	vModel.Body = sBody

	var num = int64(0)
	var err error
	if vModel.Id == "" {
		num, err = services.ConfigRemoteServiceAdd(&vModel)
	} else {
		err = services.ConfigRemoteServiceUpdate(&vModel)
	}

	common.Logger.Info("services.ConfigRemoteServiceAdd  num %d ", num)

	if err != nil {
		response["code"] = utils.FailedCode
		response["msg"] = err.Error()
	} else {

		response["code"] = utils.SuccessCode
		response["msg"] = "接口调用成功"

		request, err := SendRequest(vModel)
		if err != nil {
			common.Logger.Error(err.Error())
		}
		var req = request.(*httplib.BeegoHTTPRequest)

		resp, err := req.Response()
		headMap := resp.Header
		var isBinary = false
		contentType := headMap.Get("Content-Type")
		if strings.Index(contentType, "octet-stream") >= 0 {
			isBinary = true
		}

		for k, v := range headMap {
			ctl.Ctx.Output.Header(k, v[0])
		}

		if isBinary && operType == "down" {
			strBody, err := req.Bytes()

			if err != nil {
				response["code"] = utils.FailedCode
				response["msg"] = err.Error()
			} else {
				ctl.Ctx.Output.Body(strBody)
			}
			response["data"] = vModel
		} else {
			strBody, err := req.String()

			if err != nil {
				response["code"] = utils.FailedCode
				response["msg"] = err.Error()
			} else {
				response["data"] = strBody

			}
		}
		ctl.Data["json"] = response
	}

	ctl.ServeJSON()
}

func SendRequest(remote models.ConfigRemote) (resp interface{}, err error) {
	if remote.Method == "GET" {
		resp, err = utils.SendGet(remote)
	} else if remote.Method == "POST" {
		resp, err = utils.SendPost(remote)
	} else if remote.Method == "PUT" {
		resp, err = utils.SendPut(remote)
	} else if remote.Method == "DELETE" {
		resp, err = utils.SendDelete(remote)
	} else if remote.Method == "HEADER" {
		resp, err = utils.SendHeader(remote)
	}
	return resp, err
}

// @Title Get
// @Description 获取接口调用历史
// @Param page query int true "页码"
// @Param pageSize query  int true "分页大小"
// @Param moduleName query  string flase "模块名称"
// @Success 200 {object} models.BizLog "0k"
// @Failure 403 : other err
// @router /list [get]
func (ctl *RemoteController) List() {

	response := make(map[string]interface{})

	page, _ := ctl.GetInt("page")
	ctl.pageSize, _ = ctl.GetInt("pageSize")

	RemoteAddr := ctl.Input().Get("RemoteAddr")
	Body := ctl.Input().Get("Body")

	if ctl.pageSize == 0 {
		ctl.pageSize = 20
	}

	var query = make(map[string]string)
	var limit int64 = 20
	var offset = (int64)((page - 1) * ctl.pageSize)

	query["RemoteAddr__icontains"] = RemoteAddr
	query["Body__icontains"] = Body

	logList, count, err := services.ConfigRemoteServiceGetList(query, offset, limit)

	if err != nil {
		response["code"] = utils.FailedCode
		response["msg"] = err.Error()
	} else {
		response["code"] = utils.SuccessCode
		response["msg"] = utils.SuccessMsg
		response["data"] = logList
	}

	ctl.Data["param"] = query
	ctl.Data["result"] = response

	ctl.Data["pageBar"] = libs.NewPager(page, int(count), ctl.pageSize, beego.URLFor("RemoteController.List", "RemoteAddr", RemoteAddr, "Body", Body), true).ToString()

	ctl.display()

}

// Delete ...
// @Title Delete
// @Description delete the ConfigRemote
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /delete [post]
func (ctl *RemoteController) Delete() {
	response := make(map[string]interface{})

	Ids := strings.TrimSpace(ctl.GetString("Ids"))

	idArr := strings.Split(Ids, ",")

	for index, id := range idArr {
		err := services.ConfigRemoteServiceDelete(id)
		if err != nil {
			common.Logger.Info("services.ConfigRemoteServiceDelete  field id %d, %s ", index, id)
		}
	}

	response["code"] = utils.SuccessCode
	response["msg"] = utils.SuccessMsg
	ctl.Data["json"] = response
	ctl.ServeJSON()
}

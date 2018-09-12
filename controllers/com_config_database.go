package controllers

import (
	"errors"
	"logManager/models"
	"logManager/services"
	"logManager/utils"
	"strings"
)

type DatabaseController struct {
	BaseController
}

// @router /edit [get]
func (ctl *DatabaseController) Edit() {
	response := make(map[string]interface{})

	id := ctl.GetString("id")
	utils.Logger.Debug("database list ", id)

	if id != "" {
		vModel, err := services.ConfigDatabaseServiceGetById(id)
		if err != nil {
			response["code"] = utils.FailedCode
			response["msg"] = err.Error()
		} else {
			response["code"] = utils.SuccessCode
			response["msg"] = utils.SuccessMsg
		}
		response["data"] = vModel
	} else {
		response["code"] = utils.SuccessCode
		response["msg"] = utils.SuccessMsg
		response["data"] = models.ConfigDatabase{}
	}

	ctl.Data["result"] = response
	ctl.display()
}

// @router /save [post]
func (ctl *DatabaseController) Save() {

	response := make(map[string]interface{})

	vModel := models.ConfigDatabase{}
	if err := ctl.ParseForm(&vModel); err != nil {
		response["code"] = utils.FailedCode
		response["msg"] = err.Error()
	} else {
		if vModel.Id == "" {
			_, err = services.ConfigDatabaseServiceAdd(&vModel)
		} else {
			err = services.ConfigDatabaseServiceUpdate(&vModel)
		}

		if err != nil {
			response["code"] = utils.FailedCode
			response["msg"] = err.Error()
		} else {
			ctl.Ctx.Output.SetStatus(201)
			response["code"] = utils.SuccessCode
			response["msg"] = utils.SuccessMsg
		}

	}

	ctl.Data["json"] = response
	ctl.ServeJSON()

}

// @router /view [get]
func (ctl *DatabaseController) View() {

	id := ctl.GetString("id")

	utils.Logger.Debug("database list ", id)
	response := make(map[string]interface{})

	if id == "" {
		response["code"] = utils.FailedCode
		response["msg"] = errors.New("缺少参数ID")
	}
	vModel, err := services.ConfigDatabaseServiceGetById(id)

	if err != nil {
		response["code"] = utils.FailedCode
		response["msg"] = err.Error()
	} else {
		response["code"] = utils.SuccessCode
		response["msg"] = utils.SuccessMsg
		response["data"] = vModel
	}

	ctl.Data["result"] = response
	ctl.display()

}

// @router /list [get,post]
func (ctl *DatabaseController) List() {

	utils.Logger.Debug("database list ")

	response := make(map[string]interface{})
	var query = make(map[string]string)
	vlist, err := services.ConfigDatabaseServiceGetList(query)

	if err != nil {
		response["code"] = utils.FailedCode
		response["msg"] = err.Error()
	} else {
		response["code"] = utils.SuccessCode
		response["msg"] = utils.SuccessMsg
		response["data"] = vlist
	}

	ctl.Data["result"] = response
	ctl.display()

}

// Delete ...
// @Title Delete
// @Description delete the ConfigRemote
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /delete [post]
func (ctl *DatabaseController) Delete() {
	response := make(map[string]interface{})

	Ids := strings.TrimSpace(ctl.GetString("Ids"))

	idArr := strings.Split(Ids, ",")

	for index, id := range idArr {
		if id != "" {
			err := services.ConfigDatabaseServiceDelete(id)
			if err != nil {
				utils.Logger.Info("services.ConfigRemoteServiceDelete  field id %d, %s ", index, id)
			}
		}
	}

	response["code"] = utils.SuccessCode
	response["msg"] = utils.SuccessMsg
	ctl.Data["json"] = response
	ctl.ServeJSON()
}

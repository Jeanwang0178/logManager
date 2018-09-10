package controllers

import (
	"encoding/json"
	"errors"
	"logManager/models"
	"logManager/services"
	"logManager/utils"
)

type DatabaseController struct {
	BaseController
}

// @router /add [post]
func (ctl *DatabaseController) Add() {

	response := make(map[string]interface{})

	vModel := models.ConfigDatabase{}
	var vbyte = ctl.Ctx.Input.RequestBody
	err := json.Unmarshal(vbyte, &vModel)
	if err == nil {
		_, err = services.ConfigDatabaseServiceAdd(&vModel)
		if err != nil {
			response["code"] = utils.FailedCode
			response["msg"] = err.Error()
		} else {
			ctl.Ctx.Output.SetStatus(201)
			response["code"] = utils.SuccessCode
			response["msg"] = utils.SuccessMsg
		}

	} else {
		response["code"] = utils.FailedCode
		response["msg"] = err.Error()
	}

	ctl.Data["json"] = response
	ctl.ServeJSON()

}

// @router /edit [get,post]
func (ctl *DatabaseController) Edit() {

	response := make(map[string]interface{})

	vModel := models.ConfigDatabase{}
	var vbyte = ctl.Ctx.Input.RequestBody
	err := json.Unmarshal(vbyte, &vModel)
	if err == nil {
		err = services.ConfigDatabaseServiceUpdate(&vModel)
		if err != nil {
			response["code"] = utils.FailedCode
			response["msg"] = err.Error()
		} else {
			ctl.Ctx.Output.SetStatus(201)
			response["code"] = utils.SuccessCode
			response["msg"] = utils.SuccessMsg
		}

	} else {
		response["code"] = utils.FailedCode
		response["msg"] = err.Error()
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
